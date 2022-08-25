package server

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/qingwave/weave/docs"
	"github.com/qingwave/weave/pkg/authentication"
	"github.com/qingwave/weave/pkg/authentication/oauth"
	"github.com/qingwave/weave/pkg/authorization"
	"github.com/qingwave/weave/pkg/common"
	"github.com/qingwave/weave/pkg/config"
	"github.com/qingwave/weave/pkg/controller"
	"github.com/qingwave/weave/pkg/controller/kubecontroller"
	"github.com/qingwave/weave/pkg/database"
	"github.com/qingwave/weave/pkg/library/docker"
	"github.com/qingwave/weave/pkg/library/kubernetes"
	"github.com/qingwave/weave/pkg/middleware"
	"github.com/qingwave/weave/pkg/repository"
	"github.com/qingwave/weave/pkg/service"
	"github.com/qingwave/weave/pkg/utils/request"
	"github.com/qingwave/weave/pkg/utils/set"
	"github.com/qingwave/weave/pkg/version"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func New(conf *config.Config, logger *logrus.Logger) (*Server, error) {
	rateLimitMiddleware, err := middleware.RateLimitMiddleware(conf.Server.LimitConfigs)
	if err != nil {
		return nil, err
	}

	db, err := database.NewPostgres(&conf.DB)
	if err != nil {
		return nil, err
	}

	rdb, err := database.NewRedisClient(&conf.Redis)
	if err != nil {
		return nil, err
	}

	var conClient *docker.Client
	if conf.Docker.Enable {
		conClient, err = docker.NewClient(conf.Docker.Host)
		if err != nil {
			logrus.Warningf("failed to create docker client, container api disabled: %v", err)
			conf.Docker.Enable = false
		}
	}

	var kubeClient *kubernetes.KubeClient
	if conf.Kubernetes.Enable {
		kubeClient, err = kubernetes.NewClient(&conf.Kubernetes)
		if err != nil {
			logrus.Warnf("failed to create k8s client: %v", err)
			conf.Kubernetes.Enable = false
		}
	}

	repository := repository.NewRepository(db, rdb)
	if conf.DB.Migrate {
		if err := repository.Migrate(); err != nil {
			return nil, err
		}
	}

	if err := authorization.InitAuthorization(db, conf.AuthConfig); err != nil {
		return nil, err
	}

	userService := service.NewUserService(repository.User())
	groupService := service.NewGroupService(repository.Group(), repository.User())
	jwtService := authentication.NewJWTService(conf.Server.JWTSecret)
	oauthManager := oauth.NewOAuthManager(conf.OAuthConfig)

	userController := controller.NewUserController(userService)
	groupController := controller.NewGroupController(groupService)
	authController := controller.NewAuthController(userService, jwtService, oauthManager)
	containerController := controller.NewContainerController(conClient)
	rbacController := controller.NewRbacController()
	kubeController := kubecontroller.NewKubeControllers(kubeClient)
	postController := controller.NewPostController(service.NewPostService(repository.Post()))

	controllers := []controller.Controller{userController, groupController, authController, rbacController, postController}
	if conf.Docker.Enable {
		controllers = append(controllers, containerController)
	}
	if conf.Kubernetes.Enable {
		controllers = append(controllers, kubeController)
	}

	gin.SetMode(conf.Server.ENV)

	e := gin.New()
	e.Use(
		gin.Recovery(),
		rateLimitMiddleware,
		middleware.MonitorMiddleware(),
		middleware.CORSMiddleware(),
		middleware.RequestInfoMiddleware(&request.RequestInfoFactory{APIPrefixes: set.NewString("api")}),
		middleware.LogMiddleware(logger, "/"),
		middleware.AuthenticationMiddleware(jwtService),
		middleware.AuthorizationMiddleware(),
		middleware.TraceMiddleware(),
	)

	e.LoadHTMLFiles("static/terminal.html")

	return &Server{
		engine:          e,
		config:          conf,
		logger:          logger,
		db:              db,
		rdb:             rdb,
		containerClient: conClient,
		kubeClient:      kubeClient,
		controllers:     controllers,
	}, nil
}

type Server struct {
	engine *gin.Engine
	config *config.Config
	logger *logrus.Logger

	db              *gorm.DB
	rdb             *database.RedisDB
	containerClient *docker.Client
	kubeClient      *kubernetes.KubeClient

	controllers []controller.Controller
}

// graceful shutdown
func (s *Server) Run() {
	defer s.Close()

	s.initRouter()

	if s.kubeClient != nil {
		if err := s.kubeClient.StartCache(); err != nil {
			s.logger.Fatalf("failed to start kubernetes cache")
		}
	}

	addr := fmt.Sprintf("%s:%d", s.config.Server.Address, s.config.Server.Port)
	s.logger.Infof("Start server on: %s", addr)

	server := &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Fatalf("Failed to start server, %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.Server.GracefulShutdownPeriod)*time.Second)
	defer cancel()

	ch := <-sig
	s.logger.Infof("Receive signal: %s", ch)

	server.Shutdown(ctx)
}

func (s *Server) Close() {
	s.rdb.Close()
	db, _ := s.db.DB()
	if db != nil {
		db.Close()
	}
	if s.containerClient != nil {
		s.containerClient.Close()
	}
}

func (s *Server) initRouter() {
	root := s.engine

	// register non-resource routers
	root.GET("/", common.WrapFunc(s.getRoutes))
	root.GET("/index", controller.Index)
	root.GET("/healthz", common.WrapFunc(s.Ping))
	root.GET("/version", common.WrapFunc(version.Get))
	root.GET("/metrics", gin.WrapH(promhttp.Handler()))
	root.Any("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
	if gin.Mode() != gin.ReleaseMode {
		root.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	api := root.Group("/api/v1")
	controllers := make([]string, 0, len(s.controllers))
	for _, router := range s.controllers {
		router.RegisterRoute(api)
		controllers = append(controllers, router.Name())
	}
	logrus.Infof("server enabled controllers: %v", controllers)
}

func (s *Server) getRoutes() []string {
	paths := set.NewString()
	for _, r := range s.engine.Routes() {
		if r.Path != "" {
			paths.Insert(r.Path)
		}
	}
	return paths.Slice()
}

type ServerStatus struct {
	Ping  bool `json:"ping"`
	DB    bool `json:"db"`
	Redis bool `json:"redis"`
}

func (s *Server) Ping() *ServerStatus {
	status := &ServerStatus{Ping: true}

	ctx, cannel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cannel()

	db, err := s.db.DB()
	if err == nil {
		err = db.PingContext(ctx)
		if err == nil {
			status.DB = true
		}
	}
	if err != nil {
		logrus.Warnf("check db failed: %v", err)
	}

	if s.rdb != nil {
		_, err := s.rdb.Ping(ctx).Result()
		if err == nil {
			status.Redis = true
		} else {
			logrus.Warnf("check redis failed: %v", err)
		}
	}

	return status
}

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

	_ "weave/docs"
	"weave/pkg/authentication"
	"weave/pkg/authentication/oauth"
	"weave/pkg/authorization"
	"weave/pkg/config"
	"weave/pkg/controller"
	"weave/pkg/database"
	"weave/pkg/library/docker"
	"weave/pkg/middleware"
	"weave/pkg/repository"
	"weave/pkg/service"
	"weave/pkg/utils/request"
	"weave/pkg/utils/set"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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

	conClient, err := docker.NewClient(conf.Server.DockerHost)
	if err != nil {
		logrus.Warningf("failed to create docker client, container api disabled: %v", err)
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
	authContoller := controller.NewAuthController(userService, jwtService, oauthManager)
	containerController := controller.NewContainerController(conClient)
	rbacController := controller.NewRbacController()

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

	// set route
	InitRouter(e, userController, groupController, authContoller, containerController, rbacController)

	return &Server{
		engine:          e,
		config:          conf,
		logger:          logger,
		db:              db,
		rdb:             rdb,
		containerClient: conClient,
	}, nil
}

type Server struct {
	engine *gin.Engine
	config *config.Config
	logger *logrus.Logger

	db              *gorm.DB
	rdb             *database.RedisDB
	containerClient *docker.Client
}

// graceful shutdown
func (s *Server) Run() {
	defer s.Close()

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

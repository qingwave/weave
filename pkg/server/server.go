package server

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	_ "weave/docs"
	"weave/pkg/config"
	"weave/pkg/controller"
	"weave/pkg/database"
	"weave/pkg/middleware"
	"weave/pkg/repository"
	"weave/pkg/service"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New(conf *config.Config, logger *logrus.Logger) (*Server, error) {
	e := gin.New()
	e.Use(
		middleware.LogMiddleware(logger, "/api/v1"),
		gin.Recovery(),
	)

	db, err := database.InitDB(&conf.DBConfig)
	if err != nil {
		return nil, err
	}

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	return &Server{
		engine:         e,
		config:         conf,
		logger:         logger,
		userController: userController,
	}, nil
}

type Server struct {
	engine *gin.Engine
	config *config.Config
	logger *logrus.Logger

	userController *controller.UserController
}

// graceful shutdown
func (s *Server) Run() {
	s.Routers()

	addr := fmt.Sprintf("127.0.0.1:%d", s.config.Port)
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

	ctx, cancel := context.WithTimeout(context.Background(), s.config.GracefulShutdownTime)
	defer cancel()

	ch := <-sig
	s.logger.Infof("Receive signal: %s", ch)
	server.Shutdown(ctx)
}

func (s *Server) Routers() {
	root := s.engine
	root.GET("/", s.Index)
	root.GET("/healthz", s.Healthz)
	root.GET("/metrics", gin.WrapH(promhttp.Handler()))
	root.Any("/debug/pprof/*any")
	root.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := root.Group("/api/v1")
	api.GET("/users", s.userController.List)
	api.POST("/users", s.userController.Create)
	api.GET("/users/:id", s.userController.Get)
	api.PUT("/users/:id", s.userController.Update)
	api.DELETE("/users/:id", s.userController.Delete)
}

// @Summary Index
// @Produce plain
// @Tags index
// @Router / [get]
func (s *Server) Index(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(
		`<html>
	<head>
		<title>Weave Server</title>
	</head>
	<body>
		<h1>Hello Weave</h1>
		<ul>
			<li><a href="/swagger/index.html">swagger</a></li>
			<li><a href="/metrics">metrics</a></li>
			<li><a href="/healthz">healthz</a></li>
	  	</ul>
		<hr>
		<center>Weave/1.0</center>
	</body>
<html>`))
}

// @Summary Healthz
// @Produce json
// @Tags healthz
// @Router /healthz [get]
func (s *Server) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

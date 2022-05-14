package server

import (
	"weave/pkg/controller"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(e *gin.Engine, routers ...controller.Controller) {
	root := e
	registerDefaultRouter(root)

	api := root.Group("/api/v1")
	for _, router := range routers {
		router.RegisterRoute(api)
	}
}

func registerDefaultRouter(root *gin.Engine) {
	root.GET("/", controller.Index)
	root.GET("/healthz", controller.Healthz)
	root.GET("/metrics", gin.WrapH(promhttp.Handler()))
	root.Any("/debug/pprof/*any")
	root.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

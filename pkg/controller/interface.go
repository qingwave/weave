package controller

import "github.com/gin-gonic/gin"

type Controller interface {
	RegisterRoute(*gin.RouterGroup)
}

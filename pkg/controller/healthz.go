package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary Healthz
// @Produce json
// @Tags healthz
// @Success 200 {string}  string    "ok"
// @Router /healthz [get]
func Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

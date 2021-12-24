package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	success = 200
	failed  = 500
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResponse(c *gin.Context, code int, data interface{}, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	NewResponse(c, success, data, "success")
}

func ResponseFailed(c *gin.Context, code int, err error) {
	if code == 0 {
		code = failed
	}
	NewResponse(c, code, nil, err.Error())
}

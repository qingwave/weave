package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResponse(c *gin.Context, code int, data interface{}, msg string) {
	c.JSON(code, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	NewResponse(c, http.StatusOK, data, "success")
}

func ResponseFailed(c *gin.Context, code int, err error) {
	if code == 0 {
		code = http.StatusInternalServerError
	}
	var msg string
	if err != nil {
		msg = err.Error()
		user := GetUser(c)
		var name string
		if user != nil {
			name = user.Name
		}
		logrus.Warnf("url: %s, user: %s, error: %v", c.Request.URL, name, err)
	}
	NewResponse(c, code, nil, msg)
}

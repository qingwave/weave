package common

import (
	"net/http"
	"weave/pkg/model"

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
	if err != nil {
		val, _ := c.Get(UserContextKey)
		user, ok := val.(*model.User)
		var name string
		if ok {
			name = user.Name
		}
		logrus.Warnf("url: %s, user: %s, error: %v", c.Request.URL, name, err)
	}
	NewResponse(c, code, nil, err.Error())
}

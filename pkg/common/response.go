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
	if code == http.StatusUnauthorized && c.Request != nil {
		if val, err := c.Cookie(CookieTokenName); err == nil && val != "" {
			c.SetCookie(CookieTokenName, "", -1, "/", "", true, true)
			c.SetCookie(CookieLoginUser, "", -1, "/", "", true, false)
		}
	}

	var msg string
	if err != nil {
		msg = err.Error()
		user := GetUser(c)
		var name string
		if user != nil {
			name = user.Name
		}
		var url string
		if c.Request != nil {
			url = c.Request.URL.String()
		}
		logrus.Warnf("url: %s, user: %s, error: %v", url, name, msg)
	}
	NewResponse(c, code, nil, msg)
}

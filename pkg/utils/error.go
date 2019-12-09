package utils

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type HttpError struct {
	Code    int    `json:"code"`
	Key     string `json:"error"`
	Message string `json:"message"`
}

func NewHttpError(code int, key string, msg string) *HttpError {
	return &HttpError{
		Code:    code,
		Key:     key,
		Message: msg,
	}
}

// Error makes it compatible with `error` interface.
func (e *HttpError) Error() string {
	return e.Key + ": " + e.Message
}

func HttpErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		key  = "ServerError"
		msg  string
	)
	// 二话不说先打日志
	Logger.Error(err.Error())

	if he, ok := err.(*HttpError); ok {
		// 我们自定的错误
		code = he.Code
		key = he.Key
		msg = he.Message
	} else if ee, ok := err.(*echo.HTTPError); ok {
		// echo 框架的错误
		code = ee.Code
		key = http.StatusText(code)
		msg = key
	} else if err == gorm.ErrRecordNotFound {
		// 我们将 gorm 的没有找到直接返回 404
		code = http.StatusNotFound
		key = "NotFound"
		msg = "没有找到相应记录"
	} else {
		// 剩下的都是500 开了debug显示详细错误
		msg = err.Error()
	}

	// 判断 context 是否已经返回了
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			err := c.NoContent(code)
			if err != nil {
				Logger.Error(err.Error())
			}
		} else {
			err := c.JSON(code, NewHttpError(code, key, msg))
			if err != nil {
				Logger.Error(err.Error())
			}
		}
	}
}

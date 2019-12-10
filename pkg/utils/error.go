package utils

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"qinng.io/weave/pkg/utils/log"
)

var (
	InvalidIDError      = NewHttpError(http.StatusBadRequest, "InvalidID", "请在URL中提供合法的ID")
	NotMatchError       = NewHttpError(http.StatusBadRequest, "NotMatch", "请求ID不匹配")
	NotFoundError       = NewHttpError(http.StatusNotFound, "NotFound", "没有找到相应记录")
	InternalServerError = NewHttpError(http.StatusInternalServerError, "ServerError", "")
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
	var httpError *HttpError
	// 二话不说先打日志
	log.Logger.Error(err.Error())

	if he, ok := err.(*HttpError); ok {
		// 我们自定的错误
		httpError = he
	} else if ee, ok := err.(*echo.HTTPError); ok {
		// echo 框架的错误
		httpError = NewHttpError(ee.Code, http.StatusText(ee.Code), err.Error())
	} else if err == gorm.ErrRecordNotFound {
		// 我们将 gorm 的没有找到直接返回 404
		httpError = NotFoundError
	} else {
		// 剩下的都是500 开了debug显示详细错误
		httpError = InternalServerError
		httpError.Message = err.Error()
	}

	// 判断 context 是否已经返回了
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			err := c.NoContent(httpError.Code)
			if err != nil {
				log.Logger.Error(err.Error())
			}
		} else {
			err := c.JSON(httpError.Code, httpError)
			if err != nil {
				log.Logger.Error(err.Error())
			}
		}
	}
}

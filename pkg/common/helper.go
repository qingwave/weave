package common

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// WrapFunc will wrap func(args ...interface{}) (interface{}, <error>) as a Gin HandlerFunc
func WrapFunc(f interface{}, args ...interface{}) gin.HandlerFunc {
	fn := reflect.ValueOf(f)
	if fn.Type().NumIn() != len(args) {
		panic(fmt.Sprintf("invaild input parameters of function %v", fn.Type()))
	}

	outNum := fn.Type().NumOut()
	if outNum == 0 {
		panic(fmt.Sprintf("invaild output parameters of function %v, at least one, but got %d", fn.Type(), outNum))
	}

	inputs := make([]reflect.Value, len(args))
	for k, in := range args {
		inputs[k] = reflect.ValueOf(in)
	}

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Warnf("panic: %v", err)
				ResponseFailed(c, http.StatusInternalServerError, fmt.Errorf("%v", err))
			}
		}()

		outputs := fn.Call(inputs)
		if len(outputs) > 1 {
			err, ok := outputs[len(outputs)-1].Interface().(error)
			if ok && err != nil {
				ResponseFailed(c, http.StatusInternalServerError, err)
				return
			}
		}
		c.JSON(http.StatusOK, outputs[0].Interface())
	}
}

package common

import (
	"github.com/qingwave/weave/pkg/model"
	"github.com/qingwave/weave/pkg/utils/request"
	"github.com/qingwave/weave/pkg/utils/trace"

	"github.com/gin-gonic/gin"
)

func SetTrace(c *gin.Context, t *trace.Trace) {
	if c == nil || t == nil {
		return
	}

	c.Set(TraceContextKey, t)
}

func GetTrace(c *gin.Context) *trace.Trace {
	if c == nil {
		return nil
	}

	val, ok := c.Get(TraceContextKey)
	if !ok {
		return nil
	}

	trace, ok := val.(*trace.Trace)
	if !ok {
		return nil
	}

	return trace
}

func TraceStep(c *gin.Context, msg string, fields ...trace.Field) {
	trace := GetTrace(c)
	if trace != nil {
		trace.Step(msg, fields...)
	}
}

func SetUser(c *gin.Context, user *model.User) {
	if c == nil || user == nil {
		return
	}

	c.Set(UserContextKey, user)
}

func GetUser(c *gin.Context) *model.User {
	if c == nil {
		return nil
	}

	val, ok := c.Get(UserContextKey)
	if !ok {
		return nil
	}

	user, ok := val.(*model.User)
	if !ok {
		return nil
	}

	return user
}

func SetRequestInfo(c *gin.Context, ri *request.RequestInfo) {
	if c == nil || ri == nil {
		return
	}

	c.Set(RequestInfoContextKey, ri)
}

func GetRequestInfo(c *gin.Context) *request.RequestInfo {
	if c == nil {
		return nil
	}

	val, ok := c.Get(RequestInfoContextKey)
	if !ok {
		return nil
	}

	ri, ok := val.(*request.RequestInfo)
	if !ok {
		return nil
	}

	return ri
}

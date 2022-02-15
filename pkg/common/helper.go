package common

import (
	"weave/pkg/utils/trace"

	"github.com/gin-gonic/gin"
)

func GetTrace(c *gin.Context) *trace.Trace {
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

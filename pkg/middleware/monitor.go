package middleware

import (
	"strconv"
	"time"

	"github.com/qingwave/weave/pkg/metrics"

	"github.com/gin-gonic/gin"
)

func MonitorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		metrics.HTTPInflightRequests.WithLabelValues(method, path).Inc()

		defer func() {
			latency := float64(time.Since(start) / time.Second)
			code := c.Writer.Status()
			metrics.HTTPInflightRequests.WithLabelValues(method, path).Dec()
			metrics.HTTPRequestsTotal.WithLabelValues(method, path, strconv.Itoa(code)).Inc()
			metrics.HTTPRequestsDuration.WithLabelValues(method, path).Observe(latency)
		}()

		c.Next()
	}
}

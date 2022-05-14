package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	hostname, _ = os.Hostname()
)

func LogMiddleware(logger *logrus.Logger, pathPrefix ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		logged := len(pathPrefix) == 0
		for _, prefix := range pathPrefix {
			if strings.HasPrefix(path, prefix) {
				logged = true
				break
			}
		}
		if !logged {
			return
		}

		start := time.Now()

		defer func() {
			latency := time.Since(start)
			statusCode := c.Writer.Status()
			clientIP := c.ClientIP()
			clientUserAgent := c.Request.UserAgent()

			entry := logger.WithFields(logrus.Fields{
				"hostname":   hostname,
				"path":       path,
				"method":     c.Request.Method,
				"statusCode": statusCode,
				"clientIP":   clientIP,
				"userAgent":  clientUserAgent,
			})

			if len(c.Errors) > 0 {
				entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
			} else {
				msg := fmt.Sprintf("[%s %s] %d %v", c.Request.Method, c.Request.URL, statusCode, latency)
				if statusCode >= http.StatusInternalServerError {
					entry.Error(msg)
				} else if statusCode >= http.StatusBadRequest {
					entry.Warn(msg)
				} else {
					entry.Info(msg)
				}
			}
		}()

		c.Next()
	}
}

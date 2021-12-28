package middleware

import (
	"errors"
	"net/http"
	"strings"

	"weave/pkg/common"
	"weave/pkg/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func AuthMiddleware(jwtService *service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			common.ResponseFailed(c, http.StatusForbidden, errors.New("No Authorization header provided"))
			c.Abort()
			return
		}
		token := strings.Fields(auth)
		if len(token) != 2 || token[0] != "Bearer" || token[1] == "" {
			common.ResponseFailed(c, http.StatusForbidden, errors.New("Authorization header invaild"))
			c.Abort()
			return
		}

		user, err := jwtService.ParseToken(token[1])
		if err != nil {
			common.ResponseFailed(c, http.StatusForbidden, errors.New("Authorization header invaild"))
			c.Abort()
			return
		}

		logrus.Infof("auth user %#v", user)

		c.Next()
	}
}

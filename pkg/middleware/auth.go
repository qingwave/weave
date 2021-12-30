package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"weave/pkg/common"
	"weave/pkg/service"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService *service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := getTokenFromAuthorizationHeader(c)
		if token == "" {
			token, _ = getTokenFromCookie(c)
		}

		if token == "" {
			common.ResponseFailed(c, http.StatusUnauthorized, errors.New("Authorization info invaild"))
			c.Abort()
			return
		}

		user, err := jwtService.ParseToken(token)
		if err != nil {
			common.ResponseFailed(c, http.StatusForbidden, errors.New("Authorization failed"))
			c.Abort()
			return
		}

		c.Set(common.UserContextKey, user)

		c.Next()
	}
}

func getTokenFromCookie(c *gin.Context) (string, error) {
	return c.Cookie("token")
}

func getTokenFromAuthorizationHeader(c *gin.Context) (string, error) {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		return "", nil
	}

	token := strings.Fields(auth)
	if len(token) != 2 || strings.ToLower(token[0]) != "bearer" || token[1] == "" {
		return "", fmt.Errorf("Authorization header invaild")
	}

	return token[1], nil
}

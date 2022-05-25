package middleware

import (
	"fmt"
	"net/http"

	"github.com/qingwave/weave/pkg/authorization"
	"github.com/qingwave/weave/pkg/common"
	"github.com/qingwave/weave/pkg/model"

	"github.com/gin-gonic/gin"
)

func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := common.GetUser(c)
		if user == nil {
			user = &model.User{}
		}

		ri := common.GetRequestInfo(c)
		if ri == nil {
			common.ResponseFailed(c, http.StatusBadRequest, fmt.Errorf("failed to get request info"))
			c.Abort()
			return
		}

		if ri.IsResourceRequest {
			resource := ri.Resource
			if !authorization.Enforce(user.Name, ri.Namespace, resource, ri.Name, ri.Verb) {
				common.ResponseFailed(c, http.StatusForbidden, nil)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

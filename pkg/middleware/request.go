package middleware

import (
	"net/http"

	"github.com/qingwave/weave/pkg/common"
	"github.com/qingwave/weave/pkg/utils/request"

	"github.com/gin-gonic/gin"
)

func RequestInfoMiddleware(resolver request.RequestInfoResolver) gin.HandlerFunc {
	return func(c *gin.Context) {
		ri, err := resolver.NewRequestInfo(c.Request)
		if err != nil {
			common.ResponseFailed(c, http.StatusBadRequest, err)
			c.Abort()
			return
		}

		common.SetRequestInfo(c, ri)

		c.Next()
	}
}

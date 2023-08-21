package middleware

import (
	"github.com/cross-space-official/common/consts"
	"github.com/cross-space-official/common/failure"
	"github.com/gin-gonic/gin"
)

func ApiKeyAuthenticator(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Request.Header.Get(consts.ApiKeyHeader)
		if key != apiKey {
			_ = c.Error(failure.GeneratePlainFailure(failure.NotAuthenticatedError, "", ""))
			c.Abort()
			return
		}

		c.Next()
	}
}

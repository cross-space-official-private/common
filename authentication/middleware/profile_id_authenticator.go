package middleware

import (
	"github.com/cross-space-official-private/common/consts"
	"github.com/cross-space-official-private/common/failure"
	"github.com/gin-gonic/gin"
)

func ProfileIdAuthenticator() gin.HandlerFunc {
	return func(c *gin.Context) {
		profileId := c.Request.Header.Get(consts.ProfileIDHeader)
		if len(profileId) == 0 {
			_ = c.Error(failure.GeneratePlainFailure(failure.NotAuthenticatedError, "", ""))
			c.Abort()
			return
		}

		c.Next()
	}
}

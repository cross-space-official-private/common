package middleware

import (
	"github.com/cross-space-official-private/common/utils"
	"github.com/gin-gonic/gin"
)

type RuntimeErrorHandler func(c *gin.Context, err error)

func GinLocalRuntimeErrorHandler(handler RuntimeErrorHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		HandleRuntimeErrors(c, handler)
	}
}

func HandleRuntimeErrors(c *gin.Context, handler RuntimeErrorHandler) {
	err := c.Errors.Last()

	if utils.IsNil(err) {
		return
	}

	handler(c, err)
}

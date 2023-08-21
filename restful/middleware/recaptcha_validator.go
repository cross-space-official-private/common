package middleware

import (
	"github.com/cross-space-official/common/consts"
	"github.com/cross-space-official/common/correlationId/middleware"
	"github.com/cross-space-official/common/logger"
	"github.com/cross-space-official/common/restful"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	RecaptchaTokenHeader = "csp-recaptcha-token"
)

func RecaptchaValidator(actionName string, recaptchaService RecaptchaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Request.Header.Get(consts.ProfileIDHeader)
		if recaptchaService.ShouldSkipByProfileID(c, userID) {
			logger.GetLoggerEntry(c).Info("recaptcha validation skipped")
			c.Next()
			return
		}

		token := c.Request.Header.Get(RecaptchaTokenHeader)
		valid := recaptchaService.AssessToken(c, token, actionName)
		if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				restful.ErrorResponse{
					Code:          "recaptcha_validation_failure",
					Message:       "recaptcha failed validation",
					CorrelationID: middleware.GetCorrelationID(c),
				},
			)
			return
		}

		c.Next()
	}
}

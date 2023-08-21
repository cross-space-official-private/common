package middleware

import (
	"github.com/cross-space-official/common/consts"
	"github.com/cross-space-official/common/correlationId/middleware"
	"github.com/cross-space-official/common/failure"
	"github.com/cross-space-official/common/logger"
	"github.com/cross-space-official/common/ratelimiter"
	"github.com/cross-space-official/common/restful"
	"github.com/gin-gonic/gin"
	"net/http"
)

type KeyGenerator func(c *gin.Context) string

func EndpointRateLimiter(keyGen KeyGenerator, limit ratelimiter.Limit, errorMessage string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyGen(c)

		if len(key) == 0 || limit.Rate == 0 {
			return
		}

		res, err := ratelimiter.GetRateLimiter().AllowByKey(c, key, limit)
		if err != nil {
			logger.GetLoggerEntry(c).
				WithField("internal_error_message", err.Error()).
				Warn("failed to do rate limit filtering")
			return
		}
		if !res.Allowed {
			correlationID := middleware.GetCorrelationID(c)
			profileID := c.Request.Header.Get(consts.ProfileIDHeader)
			logger.GetLoggerEntry(c).
				WithField("correlation_id", correlationID).
				WithField("url", c.Request.URL.Path).
				WithField("ip", c.RemoteIP()).
				WithField("user_id", profileID).
				WithField("limit rate", limit.Rate).
				WithField("limit period", limit.Period).
				WithField("limit key", key).
				WithField("message", errorMessage).
				Warn("rate limit hit")
			c.AbortWithStatusJSON(http.StatusTooManyRequests,
				restful.ErrorResponse{Code: failure.TooManyRequestsError, Message: errorMessage, CorrelationID: correlationID})
			return
		}

		c.Next()
	}
}

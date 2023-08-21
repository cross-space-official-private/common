package middleware

import (
	"context"
	"github.com/cross-space-official-private/common/httpclient"
	"github.com/cross-space-official-private/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	CorrelationHeaderName = "Correlation-ID"
	CorrelationIDKey      = "correlationID"
)

func init() {
	logger.WithField(
		CorrelationIDKey, func(c context.Context) string { return GetCorrelationID(c) })
	httpclient.GetHttpClientFactory().WithHeader(
		CorrelationHeaderName, func(c context.Context) string { return GetCorrelationID(c) })
}

func GinCorrelationIDHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.GetLoggerEntry(c)
		cid := c.Request.Header.Get(CorrelationHeaderName)
		if len(cid) < 1 {
			temp, err := uuid.NewRandom()
			cid = temp.String()
			if err != nil {
				log.Warn("Custom uuid generation failed with error: ", err)
			}
			log.Info("Custom uuid generated: ", cid)
		}

		c.Set(CorrelationIDKey, cid)

		c.Next()
	}
}

func GetCorrelationID(c context.Context) string {
	id, ok := c.Value(CorrelationIDKey).(string)
	if !ok {
		return ""
	}
	return id
}

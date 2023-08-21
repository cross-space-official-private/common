package middleware

import (
	"bytes"
	"github.com/cross-space-official-private/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinRequestResponseLoggerInjector() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.GetLoggerEntry(c)
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		LogRequest(log, c)

		c.Next()

		LogResponse(log, blw)
	}
}

func LogRequest(logger *logrus.Entry, c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)

	logger.WithField("clientIp", c.Request.RemoteAddr).
		WithField("URL", c.Request.RequestURI).
		WithField("method", c.Request.Method).
		WithField("requestBody", string(body)).
		Info("Inbound request")

	c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
}

func LogResponse(logger *logrus.Entry, writer *bodyLogWriter) {
	logger.WithField("httpCode", writer.Status()).
		WithField("responseBody", writer.body.String()).
		Info("Outbound response")
}

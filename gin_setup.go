package common

import (
	"github.com/cross-space-official-private/common/configuration"
	correlationMiddleware "github.com/cross-space-official-private/common/correlationId/middleware"
	"github.com/cross-space-official-private/common/metrics"
	ginMiddleware "github.com/cross-space-official-private/common/restful/middleware"
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"net/http"
)

func SetupAll(e *gin.Engine, handlers ...ginMiddleware.RuntimeErrorHandler) {
	configuration.Init()

	e.Use(gin.Recovery())
	e.Use(gin.Logger())
	e.Use(correlationMiddleware.GinCorrelationIDHandler())
	e.Use(ginMiddleware.GinRequestResponseLoggerInjector())
	e.Use(ginMiddleware.GinFailureHandler())
	for _, handler := range handlers {
		e.Use(ginMiddleware.GinLocalRuntimeErrorHandler(handler))
	}
	e.Use(nrgin.Middleware(metrics.GetApplication()))

	e.HandleMethodNotAllowed = true

	e.GET("/internal/healthcheck", healthcheck)
}

func healthcheck(c *gin.Context) {
	c.String(http.StatusOK, "ping-pong")
}

package middleware

import (
	"fmt"
	"github.com/cross-space-official-private/common/businesserror"
	"net/http"

	"github.com/cross-space-official-private/common/correlationId/middleware"
	"github.com/cross-space-official-private/common/failure"
	"github.com/cross-space-official-private/common/logger"
	"github.com/cross-space-official-private/common/restful"
	"github.com/cross-space-official-private/common/utils"
	"github.com/gin-gonic/gin"
)

func GinFailureHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		HandleErrors(c)
	}
}

func HandleErrors(c *gin.Context) {
	log := logger.GetLoggerEntry(c)
	correlationID := middleware.GetCorrelationID(c)

	err := c.Errors.Last()

	if utils.IsNil(err) {
		return
	}

	fault, ok := err.Err.(businesserror.XSpaceBusinessError)
	if !ok {
		log.Errorf("An error %v occurred for request %s", err, correlationID)
		c.AbortWithStatusJSON(http.StatusInternalServerError, restful.ErrorResponse{Code: err.Error(), CorrelationID: correlationID})
		return
	}

	switch fault.Error() {
	case failure.InvalidArgumentError:
		log.Warnf("An error %#v occurred for request %s", fault, correlationID)
		c.AbortWithStatusJSON(http.StatusBadRequest, restful.ErrorResponse{Code: fault.Error(), Message: fault.Message(), CorrelationID: correlationID})
	case failure.NotAuthenticatedError:
		c.AbortWithStatusJSON(http.StatusUnauthorized, restful.ErrorResponse{Code: fault.Error(), Message: fault.Message(), CorrelationID: correlationID})
	case failure.NoAuthorizationError:
		c.AbortWithStatusJSON(http.StatusForbidden, restful.ErrorResponse{Code: fault.Error(), Message: fault.Message(), CorrelationID: correlationID})
	case failure.ResourceNotFoundError:
		log.Warnf("An error %#v occurred for request %s", fault, correlationID)
		c.AbortWithStatusJSON(http.StatusNotFound, restful.ErrorResponse{Code: fault.Error(), Message: fault.Message(), CorrelationID: correlationID})
	case failure.TooManyRequestsError:
		c.AbortWithStatusJSON(http.StatusTooManyRequests, restful.ErrorResponse{Code: fault.Error(), Message: fault.Message(), CorrelationID: correlationID})
	default:
		log.WithField("correlation_id", correlationID).
			WithField("error", fault.Error()).
			Errorf("An internal error occurred: %#v", fault)
		returnMessage := fault.Message()
		if len(fault.Message()) == 0 {
			returnMessage = fmt.Sprintf("Something wrong happened, correlation id %s", correlationID)
		}

		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			restful.ErrorResponse{
				Code:          fault.Error(),
				CorrelationID: correlationID,
				Message:       returnMessage,
			},
		)
	}
}

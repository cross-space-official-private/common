package middleware

import (
	"encoding/json"
	"errors"
	"github.com/cross-space-official-private/common/correlationId/middleware"
	"github.com/cross-space-official-private/common/failure"
	"github.com/cross-space-official-private/common/restful"
	"github.com/cross-space-official-private/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGinLocalHandler(t *testing.T) {
	// Arrange
	engine := gin.Default()
	engine.Use(GinRequestResponseLoggerInjector())
	engine.Use(middleware.GinCorrelationIDHandler())
	engine.Use(GinFailureHandler())
	engine.Use(GinLocalRuntimeErrorHandler(localAdvisor))
	engine.GET("/test", testLocalHandler)

	// Act
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httptest.NewRequest("GET", "/test", nil))

	raw, _ := ioutil.ReadAll(recorder.Body)
	var response restful.ErrorResponse
	utils.Must(json.Unmarshal(raw, &response))

	// Assert
	assert.Equal(t, recorder.Result().StatusCode, http.StatusBadRequest)
	assert.Equal(t, response.Code, failure.InvalidArgumentError)
	assert.Equal(t, response.Message, "local_error_to_failure")
}

func TestGinLocalHandlerPassThrough(t *testing.T) {
	// Arrange
	engine := gin.Default()
	engine.Use(GinRequestResponseLoggerInjector())
	engine.Use(middleware.GinCorrelationIDHandler())
	engine.Use(GinFailureHandler())
	engine.Use(GinLocalRuntimeErrorHandler(localAdvisor))
	engine.GET("/test", testLocalHandler2)

	// Act
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httptest.NewRequest("GET", "/test", nil))

	raw, _ := ioutil.ReadAll(recorder.Body)
	var response restful.ErrorResponse
	utils.Must(json.Unmarshal(raw, &response))

	// Assert
	assert.Equal(t, recorder.Result().StatusCode, http.StatusInternalServerError)
	assert.Equal(t, response.Code, "local_intercepted")
	assert.Equal(t, response.Message, "local_error_to_response")
}

func testLocalHandler(context *gin.Context) {
	err := errors.New("local_error_to_failure")
	_ = context.Error(err)
}

func testLocalHandler2(context *gin.Context) {
	err := errors.New("local_error_to_response")
	_ = context.Error(err)
}

func localAdvisor(c *gin.Context, err error) {
	correlationID := middleware.GetCorrelationID(c)
	if err.Error() == "local_error_to_failure" {
		_ = c.Error(failure.GeneratePlainFailure(failure.InvalidArgumentError, err.Error(), ""))
	} else if err.Error() == "local_error_to_response" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, restful.ErrorResponse{Code: "local_intercepted", Message: err.Error(), CorrelationID: correlationID})
	}
}

package middleware

import (
	"github.com/cross-space-official/common/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGinCorrelationIDHandler(t *testing.T) {
	// Arrange
	engine := gin.Default()
	engine.Use(GinCorrelationIDHandler())
	engine.GET("/test", testGetCorrelationIDHandler)

	// Act
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httptest.NewRequest("GET", "/test", nil))

	result, _ := ioutil.ReadAll(recorder.Body)

	// Assert
	assert.Equal(t, len(result) > 0, true)
}

func TestGinCorrelationIDWithExistedHandler(t *testing.T) {
	// Arrange
	engine := gin.Default()
	engine.Use(GinCorrelationIDHandler())
	engine.GET("/test", testGetCorrelationIDHandler)

	// Act
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/test", nil)
	request.Header.Set(CorrelationHeaderName, "test")
	engine.ServeHTTP(recorder, request)

	result, _ := ioutil.ReadAll(recorder.Body)

	// Assert
	assert.Equal(t, string(result), "test")
}

func testGetCorrelationIDHandler(c *gin.Context) {
	correlationID := GetCorrelationID(c)

	logger.GetLoggerEntry(c).Info("test")

	c.String(http.StatusOK, correlationID)
}

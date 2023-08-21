package middleware

import (
	"encoding/json"
	"github.com/cross-space-official/common/configuration"
	"github.com/cross-space-official/common/correlationId/middleware"
	"github.com/cross-space-official/common/failure"
	"github.com/cross-space-official/common/ratelimiter"
	"github.com/cross-space-official/common/restful"
	middleware2 "github.com/cross-space-official/common/restful/middleware"
	"github.com/cross-space-official/common/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestEndpointRateLimiter(t *testing.T) {
	// Skipped before find a way to mock redis
	t.Skip()

	// Arrange
	configuration.SetKeyValue("data.redis", map[string]interface{}{
		"host":                   "127.0.0.1",
		"port":                   "6379",
		"username":               "",
		"password":               "",
		"max-idle-connections":   1,
		"idle-timeout":           1,
		"max-active-connections": 1,
	})
	engine := gin.Default()
	engine.Use(middleware.GinCorrelationIDHandler())
	engine.Use(middleware2.GinRequestResponseLoggerInjector())
	engine.Use(middleware2.GinFailureHandler())
	engine.GET("/test", EndpointRateLimiter(
		func(c *gin.Context) string { return "test" }, ratelimiter.Limit{
			Burst:  0,
			Rate:   1,
			Period: time.Minute,
		}, "test_error",
	), testHandler)

	// Act
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httptest.NewRequest("GET", "/test", nil))

	raw, _ := ioutil.ReadAll(recorder.Body)
	var response restful.ErrorResponse
	utils.Must(json.Unmarshal(raw, &response))

	// Assert
	assert.Equal(t, http.StatusTooManyRequests, recorder.Code)
	assert.Equal(t, failure.TooManyRequestsError, response.Code)
}

func testHandler(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{})
}

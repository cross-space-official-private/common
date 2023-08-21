package middleware

import (
	"encoding/json"
	"github.com/cross-space-official/common/correlationId/middleware"
	"github.com/cross-space-official/common/restful"
	"github.com/cross-space-official/common/utils"
	"github.com/go-playground/assert/v2"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cross-space-official/common/businesserror"
	"github.com/cross-space-official/common/failure"
	"github.com/gin-gonic/gin"
	"runtime/debug"
)

func TestGinFailureHandler(t *testing.T) {
	// Arrange
	engine := gin.Default()
	engine.Use(middleware.GinCorrelationIDHandler())
	engine.Use(GinRequestResponseLoggerInjector())
	engine.Use(GinFailureHandler())
	engine.GET("/test", testHandler)

	// Act
	recorder := httptest.NewRecorder()
	engine.ServeHTTP(recorder, httptest.NewRequest("GET", "/test", nil))

	raw, _ := ioutil.ReadAll(recorder.Body)
	var response restful.ErrorResponse
	utils.Must(json.Unmarshal(raw, &response))

	// Assert
	assert.Equal(t, recorder.Result().StatusCode, http.StatusBadRequest)
	assert.Equal(t, response.Code, failure.InvalidArgumentError)
	assert.Equal(t, response.Message, "test_message")
}

func testHandler(context *gin.Context) {
	_, err := testFunc1()
	fault, _ := err.(businesserror.XSpaceBusinessError)

	_ = context.Error(failure.GenerateFailureFromError(failure.InvalidArgumentError, fault))
}

func testFunc1() (*string, error) {
	return testFunc2()
}

func testFunc2() (*string, error) {
	return nil, testFunc3()
}

func testFunc3() error {
	return &serviceDefinedError{"test", "test_message", string(debug.Stack())}
}

type serviceDefinedError struct {
	code       string
	message    string
	stacktrace string
}

func (er *serviceDefinedError) Error() string {
	return er.message
}

func (er *serviceDefinedError) Stacktrace() string {
	return er.stacktrace
}

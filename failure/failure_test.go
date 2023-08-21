package failure

import (
	"errors"
	"github.com/go-playground/assert/v2"
	"github.com/joomcode/errorx"
	"testing"
)

func TestGeneratePlainFailure(t *testing.T) {
	// Action
	failure := GeneratePlainFailure(InvalidArgumentError, "test", "test")

	// Assert
	assert.Equal(t, InvalidArgumentError, failure.ErrorCode)
	assert.Equal(t, "test", failure.Message)
	assert.Equal(t, "test", failure.Stacktrace)
}

type TestError struct {
	code       string
	message    string
	stacktrace string
}

func (er *TestError) Error() string {
	return er.message
}

func (er *TestError) Stacktrace() string {
	return er.stacktrace
}

func TestGenerateFailureFromError(t *testing.T) {
	// Arrange
	internalErr := errors.New("test_error")
	err := &TestError{"test", internalErr.Error(), "test"}

	// Action
	failure := GenerateFailureFromError(InvalidArgumentError, err)

	// Assert
	assert.Equal(t, InvalidArgumentError, failure.ErrorCode)
	assert.Equal(t, "test_error", failure.Message)
	assert.Equal(t, "test", failure.Stacktrace)
}

func TestGenerateFailureFromErrorx(t *testing.T) {
	// Arrange
	errx := callA()

	// Action
	failure := GenerateFailureFromErrorx(InvalidArgumentError, errx)

	// Assert
	assert.Equal(t, InvalidArgumentError, failure.ErrorCode)
	assert.Equal(t, "test_error", failure.Message)
	assert.Equal(t, "test_error\n at github.com/cross-space-official-private/common/failure.callB()\n\t/Users/sli/Development/common/failure/failure_test.go:67\n at github.com/cross-space-official-private/common/failure.callA()\n\t/Users/sli/Development/common/failure/failure_test.go:62\n at github.com/cross-space-official-private/common/failure.TestGenerateFailureFromErrorx()\n\t/Users/sli/Development/common/failure/failure_test.go:50\n at testing.tRunner()\n\t/usr/local/go/src/testing/testing.go:1439\n at runtime.goexit()\n\t/usr/local/go/src/runtime/asm_arm64.s:1263", failure.Stacktrace)
}

func callA() *errorx.Error {
	return callB()
}

func callB() *errorx.Error {
	internalErr := errors.New("test_error")
	return errorx.EnsureStackTrace(internalErr)
}

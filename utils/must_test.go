package utils

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMustThrowsPanic(t *testing.T) {
	assert.Panics(t, func() { Must(functionWithError()) })
}

func TestMustWithoutError(t *testing.T) {
	assert.NotPanics(t, func() { Must(functionWithoutError()) })
}

func functionWithError() error {
	return errors.New("some_err")
}

func functionWithoutError() error {
	return nil
}

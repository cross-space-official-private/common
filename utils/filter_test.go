package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterNotNil(t *testing.T) {
	obj := &foo{1}
	source := []*foo{nil, obj, nil}

	target := FilterNotNil(source)

	assert.True(t, len(target) == 1)
	assert.Equal(t, obj, target[0])
}

func TestFilter(t *testing.T) {
	obj := foo{1}
	source := []foo{obj}
	predictor1 := func(val foo) bool { return val.Value == 1 }
	predictor2 := func(val foo) bool { return val.Value == 2 }

	target1 := Filter(source, predictor1)
	target2 := Filter(source, predictor2)

	assert.True(t, len(target1) == 1)
	assert.Equal(t, obj, target1[0])
	assert.True(t, len(target2) == 0)
}

type foo struct {
	Value int
}

package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap(t *testing.T) {
	obj := foo{1}
	source := []foo{obj}

	target := Map(source, func(val foo) foo { return foo{val.Value + 1} })

	assert.True(t, len(target) == 1)
	assert.Equal(t, 2, target[0].Value)
}

func TestFlatMap(t *testing.T) {
	source := [][]foo{{foo{1}, foo{2}}, {foo{5}}}

	target := FlatMap(source, func(val foo) foo { return foo{val.Value + 1} })

	assert.True(t, len(target) == 3)
	assert.Equal(t, 2, target[0].Value)
	assert.Equal(t, 3, target[1].Value)
	assert.Equal(t, 6, target[2].Value)
}

func TestRemoveDuplicate(t *testing.T) {
	source := []string{"1", "2", "5", "1"}

	target := RemoveDuplicate(source)

	assert.True(t, len(target) == 3)
	assert.True(t, Contains(target, "1"))
	assert.True(t, Contains(target, "2"))
	assert.True(t, Contains(target, "5"))
}

func TestRemoveDuplicateBy(t *testing.T) {
	source := []foo{{1}, {2}, {1}, {2}}

	target := RemoveDuplicateBy(source, func(val foo) int { return val.Value })

	assert.True(t, len(target) == 2)
}

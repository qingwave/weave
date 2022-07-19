package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	s := NewString("key1", "key3", "key2")

	assert.True(t, s.Has("key1"))
	assert.False(t, s.Has("key"))

	assert.True(t, s.HasAll("key1", "key2"))
	assert.True(t, s.HasAny("key", "key1"))

	assert.Equal(t, []string{"key1", "key2", "key3"}, s.Slice())

	s.Insert("key4")
	assert.True(t, s.Has("key4"))

	s.Delete("key4")
	assert.False(t, s.Has("key4"))
}

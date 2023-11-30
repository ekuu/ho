package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	c := New[string, any](WithInterval[string, any](100 * time.Millisecond))
	defer c.Delete("test")

	// 测试 Set 和 Get
	c.Set("test", "value")
	v, ok := c.Get("test")
	assert.True(t, ok)
	assert.Equal(t, "value", v)

	// 测试 Delete
	c.Delete("test")
	_, ok = c.Get("test")
	assert.False(t, ok)

	// 测试 SetTTL
	c.SetTTL("test", "value", time.Millisecond)
	time.Sleep(time.Millisecond * 2)
	_, ok = c.Get("test")
	assert.False(t, ok)

	// 测试 SetExpireAt
	expireAt := time.Now().Add(time.Millisecond)
	c.SetExpireAt("test", "value", expireAt)
	time.Sleep(time.Millisecond * 2)
	_, ok = c.Get("test")
	assert.False(t, ok)

	// 测试 Template
	v, err := c.Template("test", func(k string) (interface{}, error) {
		return "value", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "value", v)

	// 测试自动清理
	c.SetTTL("test-clean", "value", time.Millisecond)
	time.Sleep(time.Millisecond * 2)
	assert.Equal(t, len(c.m), 2)
	time.Sleep(time.Millisecond * 120)
	assert.Equal(t, len(c.m), 1)
}

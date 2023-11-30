package cache

import (
	"sync"
	"time"
)

//go:generate gogen option -n Cache -i m,mu -p _ --with-init
type Cache[K comparable, V any] struct {
	m        map[K]data[V]
	mu       sync.Mutex
	interval time.Duration
	ttl      time.Duration
}

func (c *Cache[K, V]) init() {
	c.m = make(map[K]data[V])
	go c.clear()
}

func (c *Cache[K, V]) clear() {
	if c.interval == 0 {
		return
	}
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		now := time.Now()
		c.mu.Lock()
		if len(c.m) == 0 {
			c.mu.Unlock()
			continue
		}
		for k, v := range c.m {
			if v.expireAt.IsZero() {
				continue
			}
			if v.expireAt.Before(now) {
				c.delete(k)
			}
		}
		c.mu.Unlock()
	}
}

func (c *Cache[K, V]) Set(k K, v V) {
	c.SetTTL(k, v, c.ttl)
}

func (c *Cache[K, V]) SetTTL(k K, v V, ttl time.Duration) {
	d := data[V]{v: v}
	if ttl > 0 {
		d.expireAt = time.Now().Add(ttl)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[k] = d
}

func (c *Cache[K, V]) SetExpireAt(k K, v V, expireAt time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[k] = data[V]{v: v, expireAt: expireAt}
}

func (c *Cache[K, V]) Delete(k K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.delete(k)
}

func (c *Cache[K, V]) delete(k K) {
	delete(c.m, k)
}

func (c *Cache[K, V]) Get(k K) (v V, ok bool) {
	c.mu.Lock()
	d, ok := c.m[k]
	c.mu.Unlock()
	if !ok {
		return v, false
	}
	if !d.expireAt.IsZero() && d.expireAt.Before(time.Now()) {
		c.delete(k)
		return v, false
	}
	return d.v, true
}

func (c *Cache[K, V]) Template(k K, fn func(k K) (V, error)) (v V, err error) {
	if v, ok := c.Get(k); ok {
		return v, nil
	}
	if v, err = fn(k); err == nil {
		c.Set(k, v)
	}
	return
}

type data[V any] struct {
	v        V
	expireAt time.Time // 零值则不过期
}

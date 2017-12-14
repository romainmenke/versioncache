package versioncache

import (
	"context"
	"sync"
	"sync/atomic"
)

type VersionCache struct {
	sync.RWMutex
	objects map[string]*object
}

type object struct {
	sync.RWMutex
	value interface{}
}

func New() *VersionCache {
	return &VersionCache{
		objects: make(map[string]*object),
	}
}

func (c *VersionCache) Version() {
	c.Lock()
	defer c.Unlock()
	c.objects = make(map[string]*object)
}

func (c *VersionCache) Get(key string) interface{} {
	c.RLock()
	o, ok := c.objects[key]
	c.RUnlock()

	if !ok {
		return nil
	}

	o.RLock()
	defer o.RUnlock()

	return o.value
}

func (c *VersionCache) Set(key string, value interface{}) {
	c.RLock()
	o, ok := c.objects[key]
	c.RUnlock()

	if !ok {
		c.Lock()
		defer c.Unlock()
		c.objects[key] = &object{
			value: value,
		}
		return
	}

	o.Lock()
	defer o.Unlock()
	o.value = value
}

func (c *VersionCache) Setter(ctx context.Context, key string) func(interface{}) {
	c.RLock()
	o, ok := c.objects[key]
	c.RUnlock()

	if !ok {
		c.Lock()
		o = &object{}
		c.objects[key] = o
		c.Unlock()
	}

	o.Lock()

	var receivedValue uint64
	objectChan := make(chan interface{}, 1)

	go func() {
		defer o.Unlock()
		defer close(objectChan)
		defer atomic.AddUint64(&receivedValue, 1)

		for {
			select {
			case value := <-objectChan:
				o.value = value
				return
			case <-ctx.Done():
				return
			}
		}
	}()

	return func(value interface{}) {
		if atomic.LoadUint64(&receivedValue) > 0 {
			return
		}

		atomic.AddUint64(&receivedValue, 1)
		objectChan <- value
	}
}

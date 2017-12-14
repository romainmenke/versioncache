package versioncache

import (
	"context"
	"sync"
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
		defer c.Unlock()

		o = &object{}
		c.objects[key] = o
		o.Lock()

		objectChan := make(chan interface{}, 1)

		go func() {
			defer o.Unlock()
			select {
			case value := <-objectChan:
				o.value = value
			case <-ctx.Done():
				break
			}
		}()

		return func(value interface{}) {
			objectChan <- value
		}
	}

	o.Lock()

	objectChan := make(chan interface{}, 1)

	go func() {
		defer o.Unlock()
		select {
		case value := <-objectChan:
			o.value = value
		case <-ctx.Done():
			break
		}
	}()

	return func(value interface{}) {
		objectChan <- value
	}
}

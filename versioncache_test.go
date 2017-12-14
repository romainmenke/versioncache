package versioncache_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/romainmenke/versioncache"
)

func TestNew(t *testing.T) {
	c := versioncache.New()
	if c == nil {
		t.Fatal("unexpected nil")
	}
}

func TestSet(t *testing.T) {
	c := versioncache.New()

	key := fmt.Sprintf("%v", time.Now())

	c.Set(key, int(1))
	c.Set(key, int(2))

	res := c.Get(key)
	if v, ok := res.(int); ok && v == 2 {
		// success
		return
	}

	t.Fatalf("expected : 2, got : %v", res)
}

func TestSetB(t *testing.T) {
	c := versioncache.New()

	wg := &sync.WaitGroup{}

	key := fmt.Sprintf("%v", time.Now())

	for y := 0; y < 100; y++ {
		wg.Add(1)

		go func(v int) {
			defer wg.Done()

			c.Set(key, v)
		}(y)
	}

	for y := 0; y < 5; y++ {
		wg.Add(1)

		go func(v int) {
			defer wg.Done()

			c.Version()
		}(y)
	}

	for y := 0; y < 100; y++ {
		wg.Add(1)

		go func(v int) {
			defer wg.Done()

			res := c.Get(key)
			_ = res
		}(y)
	}

	wg.Wait()
}

func TestWillSet(t *testing.T) {
	c := versioncache.New()

	key := fmt.Sprintf("%v", time.Now())

	set := c.Setter(context.Background(), key)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		res := c.Get(key)
		if v, ok := res.(int); ok && v == 2 {
			// success
			return
		}

		t.Fatalf("expected : 2, got : %v", res)
	}()

	time.Sleep(time.Millisecond * 100)

	set(int(2))

	wg.Wait()
}

func TestWillSetCancel(t *testing.T) {
	c := versioncache.New()

	key := fmt.Sprintf("%v", time.Now())
	c.Set(key, 1)

	ctx, cancel := context.WithCancel(context.Background())
	set := c.Setter(ctx, key)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		res := c.Get(key)
		if v, ok := res.(int); ok && v == 1 {
			// success
			return
		}

		t.Fatalf("expected : 1, got : %v", res)
	}()

	cancel()

	time.Sleep(time.Millisecond * 5)

	set(int(2))

	wg.Wait()

}

func TestWillSetB(t *testing.T) {
	c := versioncache.New()

	key := fmt.Sprintf("%v", time.Now())

	set := c.Setter(context.Background(), key)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		res := c.Get(key)
		if v, ok := res.(int); ok && v == 2 {
			// success
			return
		}

		t.Fatalf("expected : 2, got : %v", res)
	}()

	time.Sleep(time.Millisecond * 100)

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.Set(key, 1)
	}()

	time.Sleep(time.Millisecond * 100)

	set(int(2))

	wg.Wait()
}

func TestWillSetC(t *testing.T) {
	c := versioncache.New()

	key := fmt.Sprintf("%v", time.Now())

	c.Set(key, 1)

	set := c.Setter(context.Background(), key)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		res := c.Get(key)
		if v, ok := res.(int); ok && v == 2 {
			// success
			return
		}

		t.Fatalf("expected : 2, got : %v", res)
	}()

	time.Sleep(time.Millisecond * 100)

	set(int(2))

	wg.Wait()
}

func TestWillSetD(t *testing.T) {
	c := versioncache.New()

	key := fmt.Sprintf("%v", time.Now())

	c.Set(key, 1)

	set := c.Setter(context.Background(), key)

	set(int(2))
	set(int(3))

	res := c.Get(key)
	if v, ok := res.(int); ok && v == 2 {
		// success
		return
	}

	t.Fatalf("expected : 2, got : %v", res)
}

func TestWillSetWithVersion(t *testing.T) {
	c := versioncache.New()

	key := fmt.Sprintf("%v", time.Now())

	set := c.Setter(context.Background(), key)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		res := c.Get(key)
		if v, ok := res.(int); ok && v == 2 {
			// success
			return
		}

		t.Fatalf("expected : 2, got : %v", res)
	}()

	time.Sleep(time.Millisecond * 100)

	c.Version()

	set(int(2))

	res := c.Get(key)
	if v, ok := res.(int); ok && v == 2 {
		t.Fatalf("expected : nil, got : %v", res)
	}

	wg.Wait()
}

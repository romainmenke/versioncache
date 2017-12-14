package versioncache_test

import (
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

	set := c.Setter(key)

	go func() {
		res := c.Get(key)
		if v, ok := res.(int); ok && v == 2 {
			// success
			return
		}

		t.Fatalf("expected : 2, got : %v", res)
	}()

	time.Sleep(time.Millisecond * 100)

	set(int(2))
}

func TestWillSetB(t *testing.T) {
	c := versioncache.New()

	key := fmt.Sprintf("%v", time.Now())

	set := c.Setter(key)

	go func() {
		res := c.Get(key)
		if v, ok := res.(int); ok && v == 2 {
			// success
			return
		}

		t.Fatalf("expected : 2, got : %v", res)
	}()

	time.Sleep(time.Millisecond * 100)

	go c.Set(key, 1)

	time.Sleep(time.Millisecond * 100)

	set(int(2))
}

func TestWillSetC(t *testing.T) {
	c := versioncache.New()

	key := fmt.Sprintf("%v", time.Now())

	c.Set(key, 1)

	set := c.Setter(key)

	go func() {
		res := c.Get(key)
		if v, ok := res.(int); ok && v == 2 {
			// success
			return
		}

		t.Fatalf("expected : 2, got : %v", res)
	}()

	time.Sleep(time.Millisecond * 100)

	set(int(2))
}

func TestWillSetWithVersion(t *testing.T) {
	c := versioncache.New()

	key := fmt.Sprintf("%v", time.Now())

	set := c.Setter(key)

	go func() {
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

}

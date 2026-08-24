package cache

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGetOrComputeCachesResult(t *testing.T) {
	c := New()
	var calls int32

	get := func() (any, error) {
		atomic.AddInt32(&calls, 1)
		return "value", nil
	}

	for i := 0; i < 3; i++ {
		v, err := c.GetOrCompute("k", time.Minute, get)
		if err != nil || v != "value" {
			t.Fatalf("got %v, %v", v, err)
		}
	}
	if calls != 1 {
		t.Fatalf("expected 1 computation, got %d", calls)
	}
}

func TestGetOrComputeExpiry(t *testing.T) {
	c := New()
	var calls int32

	get := func() (any, error) {
		atomic.AddInt32(&calls, 1)
		return "value", nil
	}

	c.GetOrCompute("k", time.Millisecond, get)
	time.Sleep(5 * time.Millisecond)
	c.GetOrCompute("k", time.Millisecond, get)

	if calls != 2 {
		t.Fatalf("expected 2 computations after expiry, got %d", calls)
	}
}

func TestGetOrComputeDoesNotCacheErrors(t *testing.T) {
	c := New()
	boom := errors.New("boom")

	if _, err := c.GetOrCompute("k", time.Minute, func() (any, error) { return nil, boom }); !errors.Is(err, boom) {
		t.Fatalf("expected error, got %v", err)
	}
	// A failed computation must not poison the key.
	if v, err := c.GetOrCompute("k", time.Minute, func() (any, error) { return "ok", nil }); err != nil || v != "ok" {
		t.Fatalf("got %v, %v", v, err)
	}
}

func TestGetOrComputeDeduplicatesConcurrentCalls(t *testing.T) {
	c := New()
	var calls int32
	var wg sync.WaitGroup

	const n = 32
	start := make(chan struct{})
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			v, err := c.GetOrCompute("k", time.Minute, func() (any, error) {
				atomic.AddInt32(&calls, 1)
				time.Sleep(20 * time.Millisecond) // hold the door open
				return "value", nil
			})
			if err != nil || v != "value" {
				t.Errorf("got %v, %v", v, err)
			}
		}()
	}
	close(start)
	wg.Wait()

	if calls != 1 {
		t.Fatalf("expected single computation under concurrency, got %d", calls)
	}
}

func TestInvalidate(t *testing.T) {
	c := New()
	var calls int32
	get := func() (any, error) {
		atomic.AddInt32(&calls, 1)
		return "value", nil
	}

	c.GetOrCompute("k", time.Minute, get)
	c.Invalidate()
	c.GetOrCompute("k", time.Minute, get)

	if calls != 2 {
		t.Fatalf("expected recompute after invalidate, got %d calls", calls)
	}
}

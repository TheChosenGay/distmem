package cache

import (
	"fmt"
	"sync"
	"testing"
)

func TestBaseGetAndSet(t *testing.T) {
	c := NewDefaultCache()

	_, err := c.Get("key")
	if err != ErrKeyNotFound {
		t.Errorf("expected error, got nil")
	}

	err = c.Set("key", "value")
	if err != nil {
		t.Errorf("expected nil, got error: %v", err)
	}

	err = c.Set("key", "value")
	if err != ErrKeyExists {
		t.Errorf("expected key not found error, got nil")
	}

	if err := c.Delete("key"); err != nil {
		t.Errorf("expected nil, got error: %v", err)
	}

	err = c.Set("key2", "value2")
	if err != nil {
		t.Errorf("expected nil, got error: %v", err)
	}

	err = c.Set("key", "value")
	if err != nil {
		t.Errorf("expected nil, got error: %v", err)
	}

	err = c.Clear()
	if err != nil {
		t.Errorf("expected nil, got error: %v", err)
	}

}

func TestThreadSafe(t *testing.T) {
	c := NewDefaultCache()
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			c.Set(fmt.Sprintf("key%d", i), i)
		}(i)
	}
	wg.Wait()
	for i := 0; i < 100; i++ {
		value, err := c.Get(fmt.Sprintf("key%d", i))
		if err != nil {
			t.Errorf("expected nil, got error: %v", err)
		}
		if value != i {
			t.Errorf("expected %d, got %d", i, value)
		}
	}
	if err := c.Clear(); err != nil {
		t.Errorf("expected nil, got error: %v", err)
	}
}

func TestConcurrentGetAndSet(t *testing.T) {
	c := NewDefaultCache()
	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			c.Set(fmt.Sprintf("key%d", i), i)
		}(i)
	}
	wg.Wait()
	// concurrent get test
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			defer wg.Done()
			if value, err := c.Get(fmt.Sprintf("key%d", i)); err != nil || value != i {
				t.Errorf("expected nil, got error: %v", err)
			}
		}(i)
	}

	wg.Wait()
	for i := 0; i < 100; i++ {
		value, err := c.Get(fmt.Sprintf("key%d", i))
		if err != nil {
			t.Errorf("expected nil, got error: %v", err)
		}
		if value != i {
			t.Errorf("expected %d, got %d", i, value)
		}
	}
}

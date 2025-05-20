package sievecache

import (
	"sync"
	"testing"
	"time"
)

func TestSyncSieveCache(t *testing.T) {
	cache, err := NewSync[string, string](100)
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Insert a value
	inserted := cache.Insert("key1", "value1")
	if !inserted {
		t.Error("Expected insert to return true for new key")
	}

	// Read back the value
	val, found := cache.Get("key1")
	if !found || val != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}

	// Check contains key
	if !cache.ContainsKey("key1") {
		t.Error("Expected ContainsKey to return true")
	}

	// Check capacity and length
	if cache.Capacity() != 100 {
		t.Errorf("Expected capacity 100, got %d", cache.Capacity())
	}
	if cache.Len() != 1 {
		t.Errorf("Expected length 1, got %d", cache.Len())
	}

	// Remove a value
	val, found = cache.Remove("key1")
	if !found || val != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}

	if cache.Len() != 0 {
		t.Errorf("Expected length 0, got %d", cache.Len())
	}

	if !cache.IsEmpty() {
		t.Error("Expected IsEmpty to return true")
	}
}

func TestMultithreadedAccess(t *testing.T) {
	cache, _ := NewSync[string, string](100)

	// Add some initial data
	cache.Insert("shared", "initial")

	var wg sync.WaitGroup
	wg.Add(2)

	// Spawn a thread that updates the cache
	go func() {
		defer wg.Done()
		cache.Insert("shared", "updated")
		cache.Insert("thread_only", "thread_value")
	}()

	// Main thread operations
	go func() {
		defer wg.Done()
		cache.Insert("main_only", "main_value")
	}()

	// Wait for goroutines to complete
	wg.Wait()

	// Verify results
	val, _ := cache.Get("shared")
	if val != "updated" {
		t.Errorf("Expected updated, got %v", val)
	}

	val, found := cache.Get("thread_only")
	if !found || val != "thread_value" {
		t.Errorf("Expected thread_value, got %v", val)
	}

	val, found = cache.Get("main_only")
	if !found || val != "main_value" {
		t.Errorf("Expected main_value, got %v", val)
	}

	if cache.Len() != 3 {
		t.Errorf("Expected length 3, got %d", cache.Len())
	}
}

func TestWithLock(t *testing.T) {
	cache, _ := NewSync[string, string](100)

	// Perform multiple operations atomically
	cache.WithLock(func(innerCache *SieveCache[string, string]) {
		innerCache.Insert("key1", "value1")
		innerCache.Insert("key2", "value2")
		innerCache.Insert("key3", "value3")

		// We can check internal state mid-transaction
		if innerCache.Len() != 3 {
			t.Errorf("Expected length 3, got %d", innerCache.Len())
		}
	})

	if cache.Len() != 3 {
		t.Errorf("Expected length 3, got %d", cache.Len())
	}
}

func TestGetMut(t *testing.T) {
	cache, _ := NewSync[string, string](100)
	cache.Insert("key", "value")

	// Modify the value in-place
	modified := cache.GetMut("key", func(value *string) {
		*value = "new_value"
	})

	if !modified {
		t.Error("Expected GetMut to return true for existing key")
	}

	// Verify the value was updated
	val, _ := cache.Get("key")
	if val != "new_value" {
		t.Errorf("Expected new_value, got %v", val)
	}

	// Try to modify a non-existent key
	modified = cache.GetMut("missing", func(_ *string) {
		t.Error("This should not be called")
	})

	if modified {
		t.Error("Expected GetMut to return false for missing key")
	}
}

func TestForEachMethods(t *testing.T) {
	cache, _ := NewSync[string, string](10)
	cache.Insert("key1", "value1")
	cache.Insert("key2", "value2")

	// Test ForEachValue
	cache.ForEachValue(func(value *string) {
		*value = *value + "_updated"
	})

	val, _ := cache.Get("key1")
	if val != "value1_updated" {
		t.Errorf("Expected value1_updated, got %v", val)
	}

	val, _ = cache.Get("key2")
	if val != "value2_updated" {
		t.Errorf("Expected value2_updated, got %v", val)
	}

	// Test ForEachEntry
	cache.ForEachEntry(func(key string, value *string) {
		if key == "key1" {
			*value = *value + "_special"
		}
	})

	val, _ = cache.Get("key1")
	if val != "value1_updated_special" {
		t.Errorf("Expected value1_updated_special, got %v", val)
	}

	val, _ = cache.Get("key2")
	if val != "value2_updated" {
		t.Errorf("Expected value2_updated, got %v", val)
	}
}

func TestDeadlockPrevention(t *testing.T) {
	cache, _ := NewSync[string, int](100)

	// Add some initial data
	cache.Insert("key1", 1)
	cache.Insert("key2", 2)

	var wg sync.WaitGroup
	wg.Add(2)

	// Thread 1: Recursively accesses the cache within GetMut callback
	go func() {
		defer wg.Done()

		cache.GetMut("key1", func(value *int) {
			// This would deadlock with an unsafe implementation!
			// Attempt to get another value while modifying
			val, found := cache.Get("key2")
			if !found || val != 2 {
				t.Errorf("Expected 2, got %v", val)
			}

			// Even modify another value
			cache.Insert("key3", 3)

			*value += 10
		})
	}()

	// Thread 2: Also performs operations that would deadlock with unsafe impl
	go func() {
		defer wg.Done()

		// Sleep to ensure thread1 starts first
		time.Sleep(10 * time.Millisecond)

		// These operations would deadlock if thread1 held a lock during its callback
		cache.Insert("key4", 4)
		val, found := cache.Get("key2")
		if !found || val != 2 {
			t.Errorf("Expected 2, got %v", val)
		}
	}()

	// Both threads should complete without deadlock
	wg.Wait()

	// Verify final state
	val, _ := cache.Get("key1")
	if val != 11 { // 1 + 10
		t.Errorf("Expected 11, got %v", val)
	}

	val, found := cache.Get("key3")
	if !found || val != 3 {
		t.Errorf("Expected 3, got %v", val)
	}

	val, found = cache.Get("key4")
	if !found || val != 4 {
		t.Errorf("Expected 4, got %v", val)
	}
}

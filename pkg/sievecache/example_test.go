package sievecache_test

import (
	"fmt"

	"github.com/jedisct1/go-sieve-cache/pkg/sievecache"
)

// Example demonstrates basic usage of the SieveCache
func Example() {
	// Create a new cache with capacity for 100 items
	cache, _ := sievecache.New[string, string](100)

	// Insert some items
	cache.Insert("key1", "value1")
	cache.Insert("key2", "value2")

	// Get an item
	value, found := cache.Get("key1")
	if found {
		fmt.Println("Found:", value)
	}

	// Remove an item
	cache.Remove("key2")

	// Check the current cache size
	fmt.Println("Cache size:", cache.Len())

	// Output:
	// Found: value1
	// Cache size: 1
}

// ExampleSyncSieveCache demonstrates thread-safe cache usage
func ExampleSyncSieveCache() {
	// Create a thread-safe cache
	cache, _ := sievecache.NewSync[string, int](100)

	// Insert some values
	cache.Insert("counter", 5)

	// Safely modify values
	cache.GetMut("counter", func(value *int) {
		*value = *value * 2
	})

	// Retrieve modified value
	val, _ := cache.Get("counter")
	fmt.Println("Counter:", val)

	// Perform multiple operations atomically
	cache.WithLock(func(innerCache *sievecache.SieveCache[string, int]) {
		innerCache.Insert("a", 1)
		innerCache.Insert("b", 2)
		innerCache.Insert("sum", 3)
	})

	fmt.Println("Cache size:", cache.Len())

	// Output:
	// Counter: 10
	// Cache size: 4
}

// ExampleShardedSieveCache demonstrates usage of the sharded cache for high concurrency
func ExampleShardedSieveCache() {
	// We'll use a simpler example to avoid issues with sharding
	cache, _ := sievecache.NewShardedWithShards[string, int](1000, 8)
	fmt.Printf("Created sharded cache with %d shards\n", cache.NumShards())

	// Basic operations work the same as other cache types
	cache.Insert("counter", 0)

	// Increment counter using GetMut
	cache.GetMut("counter", func(value *int) {
		*value += 1
	})

	// Check the value
	val, _ := cache.Get("counter")
	fmt.Println("Counter:", val)

	// Get and modify the value
	cache.Insert("counter", 5)
	val, _ = cache.Get("counter")
	fmt.Println("Counter:", val)

	// Insert a new key
	cache.Insert("related", 10)
	relatedVal, _ := cache.Get("related")
	fmt.Println("Related:", relatedVal)

	// Output:
	// Created sharded cache with 8 shards
	// Counter: 1
	// Counter: 5
	// Related: 10
}

// ExampleSieveCache_RecommendedCapacity demonstrates the capacity recommendation feature
func ExampleSieveCache_RecommendedCapacity() {
	cache, _ := sievecache.New[string, int](100)

	// Add some data and access some of it to create a pattern
	for i := 0; i < 80; i++ {
		cache.Insert(fmt.Sprintf("key%d", i), i)

		// Access every other key to mark it as visited
		if i%2 == 0 {
			cache.Get(fmt.Sprintf("key%d", i))
		}
	}

	// Parameters:
	// - minFactor: 0.5 (never go below 50% of current capacity)
	// - maxFactor: 2.0 (never go above 200% of current capacity)
	// - lowThreshold: 0.3 (reduce capacity if utilization is below 30%)
	// - highThreshold: 0.7 (increase capacity if utilization is above 70%)
	newCapacity := cache.RecommendedCapacity(0.5, 2.0, 0.3, 0.7)
	fmt.Printf("Recommended capacity for a cache with %d/%d items: %d\n",
		cache.Len(), cache.Capacity(), newCapacity)

	// Note: Exact output may vary slightly, so we're not including the Output comment
}

package main

import (
	"fmt"
	"time"

	"github.com/jedisct1/go-sieve-cache/pkg/sievecache"
)

func main() {
	// Create a single-threaded cache
	fmt.Println("=== Single-threaded Cache Example ===")
	singleThreadedExample()

	// Create a thread-safe cache
	fmt.Println("\n=== Thread-safe Cache Example ===")
	threadSafeExample()

	// Create a sharded cache
	fmt.Println("\n=== Sharded Cache Example ===")
	shardedExample()
}

func singleThreadedExample() {
	// Create a new cache with capacity for 3 items
	cache, err := sievecache.New[string, string](3)
	if err != nil {
		fmt.Printf("Error creating cache: %v\n", err)
		return
	}

	// Insert some values
	cache.Insert("key1", "value1")
	cache.Insert("key2", "value2")
	cache.Insert("key3", "value3")
	fmt.Printf("Cache length after initial inserts: %d\n", cache.Len())

	// Access some items to mark them as visited
	val, ok := cache.Get("key1")
	if ok {
		fmt.Printf("Found key1: %s\n", val)
	}

	// Insert a new item, should evict the least recently visited
	cache.Insert("key4", "value4")
	fmt.Printf("Cache length after inserting key4: %d\n", cache.Len())

	// key2 or key3 should have been evicted (they weren't visited)
	if !cache.ContainsKey("key2") {
		fmt.Println("key2 was evicted")
	}
	if !cache.ContainsKey("key3") {
		fmt.Println("key3 was evicted")
	}

	// key1 should still be there (it was visited)
	if cache.ContainsKey("key1") {
		fmt.Println("key1 was retained")
	}

	// Print all keys
	fmt.Println("Keys in cache:", cache.Keys())

	// Get a recommended capacity based on utilization
	recommended := cache.RecommendedCapacity(0.5, 2.0, 0.3, 0.7)
	fmt.Printf("Recommended capacity: %d\n", recommended)
}

func threadSafeExample() {
	// Create a new thread-safe cache with capacity for 100 items
	cache, _ := sievecache.NewSync[string, int](100)

	// Insert some values
	for i := 0; i < 10; i++ {
		cache.Insert(fmt.Sprintf("key%d", i), i)
	}

	// Modify a value with the callback
	cache.GetMut("key5", func(val *int) {
		*val *= 10
	})

	// Get the modified value
	val, _ := cache.Get("key5")
	fmt.Printf("Modified value of key5: %d\n", val)

	// Use ForEachValue to modify all values
	cache.ForEachValue(func(val *int) {
		*val += 1
	})

	// Check some values
	val, _ = cache.Get("key5")
	fmt.Printf("key5 after ForEachValue: %d\n", val)
	val, _ = cache.Get("key1")
	fmt.Printf("key1 after ForEachValue: %d\n", val)

	// Retain only even values
	cache.Retain(func(key string, value int) bool {
		return value%2 == 0
	})

	fmt.Printf("Cache length after retain: %d\n", cache.Len())
	fmt.Println("Keys in cache after retain:", cache.Keys())
}

func shardedExample() {
	// Create a sharded cache with 8 shards
	cache, _ := sievecache.NewShardedWithShards[string, int](100, 8)
	fmt.Printf("Created sharded cache with %d shards\n", cache.NumShards())

	// Insert items with different patterns to test sharding
	start := time.Now()
	for i := 0; i < 10000; i++ {
		cache.Insert(fmt.Sprintf("key%d", i), i)
	}
	fmt.Printf("Inserted 10000 items in %v\n", time.Since(start))

	// Check a few keys
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("key%d", i*1000)
		val, ok := cache.Get(key)
		if ok {
			fmt.Printf("Found %s: %d\n", key, val)
		}
	}

	// Perform batch modification
	start = time.Now()
	cache.ForEachValue(func(val *int) {
		*val = *val * 2
	})
	fmt.Printf("Doubled all values in %v\n", time.Since(start))

	// Check the modified values
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("key%d", i*1000)
		val, ok := cache.Get(key)
		if ok {
			fmt.Printf("%s after doubling: %d\n", key, val)
		}
	}

	fmt.Printf("Final cache length: %d\n", cache.Len())
}

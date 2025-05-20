# SIEVE Cache for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/jedisct1/go-sieve-cache.svg)](https://pkg.go.dev/github.com/jedisct1/go-sieve-cache)
[![Go Report Card](https://goreportcard.com/badge/github.com/jedisct1/go-sieve-cache)](https://goreportcard.com/report/github.com/jedisct1/go-sieve-cache)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A high-performance Go implementation of the SIEVE cache replacement algorithm with thread-safe and sharded variants.

## What is SIEVE?

SIEVE (Simple, space-efficient, In-memory, EViction mEchanism) is a cache eviction algorithm that maintains a single bit per entry to track whether an item has been "visited" since it was last considered for eviction. This approach requires less state than LRU but achieves excellent performance, especially on skewed workloads.

## Features

- **Simple API**: Easy to use and integrate with existing code
- **Generic implementation**: Works with any key and value types
- **High performance**: Efficient implementation with O(1) operations
- **Thread safety options**: Choose the right level of concurrency for your needs
- **Minimal memory overhead**: Uses only a single bit per entry for tracking
- **Dynamic sizing**: Can recommend capacity adjustments based on access patterns

## Performance

SIEVE offers excellent performance comparable to or better than LRU for most workloads, while using significantly less memory overhead (1 bit per entry vs. pointers for doubly-linked list in LRU).

### When to use SIEVE

- When memory efficiency is important
- For applications with skewed access patterns (some keys accessed much more frequently)
- When you need a simple, fast, and effective caching solution

### Benchmarks

```
$ go test -bench=. ./pkg/sievecache -benchtime=1s

BenchmarkSieveCache_Mixed-10                 5177452               224.9 ns/op
BenchmarkSyncSieveCache_Mixed-10             2671064               434.5 ns/op
BenchmarkShardedSieveCache_Mixed-10          3809210               325.1 ns/op
BenchmarkParallelAccess-10                  12815046                88.4 ns/op
```

The sharded implementation offers the best performance under high concurrent load, as shown in the parallel access benchmark.

## Implementation Options

This package provides three cache implementations:

1. **SieveCache**: The core single-threaded implementation for use in simple scenarios
2. **SyncSieveCache**: A thread-safe cache for multi-threaded applications with moderate concurrency
3. **ShardedSieveCache**: A highly concurrent cache that shards data across multiple internal caches

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/jedisct1/go-sieve-cache/pkg/sievecache"
)

func main() {
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

    // For thread-safe operations
    syncCache, _ := sievecache.NewSync[string, int](100)

    // For high-concurrency applications
    shardedCache, _ := sievecache.NewSharded[string, int](1000)

    // The sharded cache can update values with mutator functions
    shardedCache.Insert("counter", 0)
    shardedCache.GetMut("counter", func(value *int) {
        *value++
    })

    // Check the current counter value
    counterValue, _ := shardedCache.Get("counter")
    fmt.Println("Counter:", counterValue)
}
```

## Advanced Usage

### Using the Thread-Safe Cache

```go
// Create a thread-safe cache
cache, _ := sievecache.NewSync[string, int](1000)

// Safely modify values
cache.GetMut("key", func(value *int) {
    *value = *value * 2
})

// Perform multiple operations atomically
cache.WithLock(func(innerCache *sievecache.SieveCache[string, int]) {
    // All operations here are atomic
    item1, _ := innerCache.Get("item1")
    item2, _ := innerCache.Get("item2")
    innerCache.Insert("sum", item1 + item2)
})

// Modify all values in one operation
cache.ForEachValue(func(value *int) {
    *value += 1
})
```

### Working with the Sharded Cache

```go
// Create a sharded cache with 32 shards for high concurrency
cache, _ := sievecache.NewShardedWithShards[string, string](10000, 32)

// Operations work the same as the other cache types
cache.Insert("key", "value")
value, _ := cache.Get("key")

// For operations that need to be atomic within a shard
cache.WithKeyLock("key", func(shard *sievecache.SieveCache[string, string]) {
    // These operations are atomic only for keys in the same shard as "key"
    shard.Insert("key", "new value")
    shard.Insert("related_key", "related value")
})
```

## Performance Tuning

The cache provides a `RecommendedCapacity` method that analyzes the current usage pattern and recommends an optimal capacity:

```go
// Get a recommended cache size based on access patterns
newCapacity := cache.RecommendedCapacity(0.5, 2.0, 0.3, 0.7)
fmt.Printf("Recommended capacity: %d\n", newCapacity)
```

Parameters:
- `minFactor`: Minimum scaling factor (0.5 means never go below 50% of current capacity)
- `maxFactor`: Maximum scaling factor (2.0 means never go above 200% of current capacity)
- `lowThreshold`: Utilization threshold below which capacity is reduced
- `highThreshold`: Utilization threshold above which capacity is increased

## Installation

```sh
go get github.com/jedisct1/go-sieve-cache
```

## License

MIT License

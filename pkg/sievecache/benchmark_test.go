package sievecache

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

// Fixed parameters for benchmarks
const (
	benchCacheSize  = 10000
	benchKeySize    = 100000 // Number of possible keys
	benchWorkingSet = 20000  // Number of keys in working set
	benchRandSeed   = 42     // Fixed seed for reproducible benchmarks
	benchShardCount = 16     // Number of shards for ShardedSieveCache
)

// generateKeys generates a set of keys for benchmarking
func generateKeys(count int) []string {
	keys := make([]string, count)
	for i := 0; i < count; i++ {
		keys[i] = fmt.Sprintf("key-%d", i)
	}
	return keys
}

// zipfDistribution generates a set of keys following a Zipf distribution
// This simulates a realistic cache access pattern with frequently accessed hot keys
func zipfDistribution(keyCount, sampleCount int, rng *rand.Rand) []string {
	zipf := rand.NewZipf(rng, 1.1, 1.0, uint64(keyCount-1))
	samples := make([]string, sampleCount)

	for i := 0; i < sampleCount; i++ {
		keyIndex := zipf.Uint64()
		samples[i] = fmt.Sprintf("key-%d", keyIndex)
	}

	return samples
}

// benchmarkInsert benchmarks the insertion of keys into a cache
func benchmarkInsert(b *testing.B, c Cache[string, int]) {
	keys := generateKeys(benchKeySize)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		keyIndex := i % benchKeySize
		c.Insert(keys[keyIndex], keyIndex)
	}
}

// benchmarkGet benchmarks the retrieval of keys from a cache
func benchmarkGet(b *testing.B, c Cache[string, int]) {
	// First, populate the cache
	keys := generateKeys(benchKeySize)
	for i := 0; i < benchCacheSize; i++ {
		c.Insert(keys[i%benchKeySize], i)
	}

	// Create a deterministic RNG for reproducible results
	rng := rand.New(rand.NewSource(benchRandSeed))

	// Generate access patterns following a Zipf distribution
	accessPatterns := zipfDistribution(benchKeySize, b.N, rng)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Get(accessPatterns[i])
	}
}

// benchmarkMixed benchmarks a mixed workload of inserts and gets
func benchmarkMixed(b *testing.B, c Cache[string, int]) {
	// Generate a large set of keys
	keys := generateKeys(benchKeySize)

	// Create a deterministic RNG for reproducible results
	rng := rand.New(rand.NewSource(benchRandSeed))

	// Pre-populate cache
	for i := 0; i < benchCacheSize/2; i++ {
		c.Insert(keys[i%benchKeySize], i)
	}

	b.ResetTimer()

	// Mixed workload with 80% reads and 20% writes
	for i := 0; i < b.N; i++ {
		if rng.Intn(100) < 80 {
			// Get operation
			keyIndex := rng.Intn(benchWorkingSet)
			c.Get(keys[keyIndex])
		} else {
			// Insert operation
			keyIndex := rng.Intn(benchKeySize)
			c.Insert(keys[keyIndex], i)
		}
	}
}

// Define Cache interface for benchmarking
type Cache[K comparable, V any] interface {
	Insert(key K, value V) bool
	Get(key K) (V, bool)
}

// Benchmark the base SieveCache implementation
func BenchmarkSieveCache_Insert(b *testing.B) {
	cache, _ := New[string, int](benchCacheSize)
	benchmarkInsert(b, cache)
}

func BenchmarkSieveCache_Get(b *testing.B) {
	cache, _ := New[string, int](benchCacheSize)
	benchmarkGet(b, cache)
}

func BenchmarkSieveCache_Mixed(b *testing.B) {
	cache, _ := New[string, int](benchCacheSize)
	benchmarkMixed(b, cache)
}

// Benchmark the thread-safe SyncSieveCache implementation
func BenchmarkSyncSieveCache_Insert(b *testing.B) {
	cache, _ := NewSync[string, int](benchCacheSize)
	benchmarkInsert(b, cache)
}

func BenchmarkSyncSieveCache_Get(b *testing.B) {
	cache, _ := NewSync[string, int](benchCacheSize)
	benchmarkGet(b, cache)
}

func BenchmarkSyncSieveCache_Mixed(b *testing.B) {
	cache, _ := NewSync[string, int](benchCacheSize)
	benchmarkMixed(b, cache)
}

// Benchmark the sharded implementation
func BenchmarkShardedSieveCache_Insert(b *testing.B) {
	cache, _ := NewShardedWithShards[string, int](benchCacheSize, benchShardCount)
	benchmarkInsert(b, cache)
}

func BenchmarkShardedSieveCache_Get(b *testing.B) {
	cache, _ := NewShardedWithShards[string, int](benchCacheSize, benchShardCount)
	benchmarkGet(b, cache)
}

func BenchmarkShardedSieveCache_Mixed(b *testing.B) {
	cache, _ := NewShardedWithShards[string, int](benchCacheSize, benchShardCount)
	benchmarkMixed(b, cache)
}

// Benchmark the ShardedSieveCache with different numbers of shards
func BenchmarkShardCount(b *testing.B) {
	shardCounts := []int{1, 2, 4, 8, 16, 32, 64}

	for _, shards := range shardCounts {
		b.Run(strconv.Itoa(shards), func(b *testing.B) {
			cache, _ := NewShardedWithShards[string, int](benchCacheSize, shards)
			benchmarkMixed(b, cache)
		})
	}
}

// Benchmark parallel access to ShardedSieveCache
func BenchmarkParallelAccess(b *testing.B) {
	cache, _ := NewShardedWithShards[string, int](benchCacheSize, benchShardCount)
	keys := generateKeys(benchKeySize)

	// Pre-populate cache
	for i := 0; i < benchCacheSize/2; i++ {
		cache.Insert(keys[i%benchKeySize], i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		// Each goroutine has its own RNG
		rng := rand.New(rand.NewSource(benchRandSeed))

		counter := 0
		for pb.Next() {
			counter++
			if rng.Intn(100) < 80 {
				// Get operation
				keyIndex := rng.Intn(benchWorkingSet)
				cache.Get(keys[keyIndex])
			} else {
				// Insert operation
				keyIndex := rng.Intn(benchKeySize)
				cache.Insert(keys[keyIndex], counter)
			}
		}
	})
}

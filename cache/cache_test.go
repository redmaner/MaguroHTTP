package cache

import (
	"encoding/binary"
	"fmt"
	"testing"
	"time"
)

var keys = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

func TestCache(t *testing.T) {
	c := NewCache()

	// test a cache get on an empty cache
	// this shouldn't be possible
	if ok, _, _ := c.GetFast("this should not work", 5, 10); ok {
		t.Fail()
	}

	// Insert keys
	for i, v := range keys {
		c.Set(v, i+1)
	}

	// Update keys
	for i, v := range keys {
		c.Set(v, i*2)
	}

	// Retrieve inserted keys
	for i, v := range keys {

		// We want to expire keys to test expiration
		if i > 23 {
			time.Sleep(3 * time.Second)
		}

		if ok, val := c.Get(v, 5); ok {
			fmt.Printf("Key %s has value %v\n", v, val)
		} else {
			fmt.Printf("Key %s has expired\n", v)
		}
	}
}

func BenchmarkCacheSet(b *testing.B) {
	cache := NewCache()
	b.ResetTimer()
	b.ReportAllocs()
	var key [8]byte
	for i := 0; i < b.N; i++ {
		binary.LittleEndian.PutUint64(key[:], uint64(i))
		cache.Set(string(key[:]), make([]byte, 8))
	}
}

func BenchmarkMapSet(b *testing.B) {
	m := make(map[string][]byte)
	var key [8]byte
	for i := 0; i < b.N; i++ {
		binary.LittleEndian.PutUint64(key[:], uint64(i))
		m[string(key[:])] = make([]byte, 8)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cache := NewCache()
	var key [8]byte
	for i := 0; i < b.N; i++ {
		binary.LittleEndian.PutUint64(key[:], uint64(i))
		cache.Set(string(key[:]), make([]byte, 8))
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		binary.LittleEndian.PutUint64(key[:], uint64(i))
		cache.Get(string(key[:]), 5)
	}
}

func BenchmarkCacheGetFast(b *testing.B) {
	cache := NewCache()
	var key [8]byte
	for i := 0; i < b.N; i++ {
		binary.LittleEndian.PutUint64(key[:], uint64(i))
		cache.Set(string(key[:]), make([]byte, 8))
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		binary.LittleEndian.PutUint64(key[:], uint64(i))
		cache.GetFast(string(key[:]), 5, 35)
	}
}

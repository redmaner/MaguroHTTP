// Copyright 2019 Jake van der Putten.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	if ok, _ := c.Get("this should not work", 900000000000); ok {
		t.Fail()
	}

	// Insert keys
	for i, v := range keys {
		c.Set(v, i+1)
	}

	// Update keys
	for i, v := range keys {
		if i == 23 {
			break
		}
		c.Set(v, i*2)
	}

	// Retrieve inserted keys
	for i, v := range keys {

		// We want to expire keys to test expiration
		if i > 23 {
			time.Sleep(3 * time.Second)
		}

		if ok, val := c.Get(v, 5000000000); ok {
			fmt.Printf("Key %s has value %v\n", v, val)
		} else {
			fmt.Printf("Key %s has expired\n", v)
		}
	}
}

func TestCacheFind(t *testing.T) {
	c := NewCache()

	// Insert keys
	for i, v := range keys {
		c.Set(v, i+1)
	}

	c.Find("z")

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
		cache.Get(string(key[:]), 1000000000)
	}
}

func BenchmarkCacheFind(b *testing.B) {
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
		cache.Find(string(key[:]))
	}
}

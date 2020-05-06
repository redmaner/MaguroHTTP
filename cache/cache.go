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

const (

	// defaultShards is the amount of shards in SpearCache
	defaultShards = 256

	// defaultItems is the amount of items in a single shard
	defaultItems = 1024

	// defaultNoClutter is the amount of items that will always be checked,
	// to prevent cache cluttering
	defaultNoClutter = 8
)

// SpearCache is a preallocated in memory cache. SpearCache uses a ring queue of a fixed length.
// All entries in the cache are appended in time order. Entries of the same key don't get updated,
// but appended instead. When the ring is full, the oldest entries are automatically overwritten.
// A cache get will retrieve the newest appended entry to the queue, if it exist and is not yet expired.
type SpearCache struct {
	shards [defaultShards]*shard
}

// NewCache returns an empty SpearCache
func NewCache() *SpearCache {
	return &SpearCache{}
}

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
	"time"

	"github.com/cespare/xxhash"
)

// Set is used to set a key value pair into SpearCache
// The key should always be a string. The value can be everything.
func (c *SpearCache) Set(key string, value interface{}) {

	// hash the key with xxhash and make the id
	keyHash := xxhash.Sum64([]byte(key))
	id := keyHash & (defaultShards - 1)

	// We make sure the shard exists, if it doesn't we create one
	if c.shards[id] == nil {
		c.shards[id] = newShard()
	}

	// Lock the shard for concurrency safety. We don't use defer to unlock the shard (on purpose)
	c.shards[id].lock.Lock()

	// To prevent cache cluttering we check if the previous 10 entries are equal to key
	// If the key exists in the previous 10 entries before cursor, we update that key instead
	// of appending the key value to the cache.
	for i := 1; i < 10; i++ {

		itid := c.shards[id].cursor - i

		if itid < 0 {
			itid = itid + defaultItems
		}

		// If the key exists we update the value and modtime.
		if c.shards[id].items[itid].key == keyHash {

			c.shards[id].items[itid].value = value
			c.shards[id].items[itid].modTime = uint64(time.Now().UnixNano())

			// Unlock and return
			c.shards[id].lock.Unlock()
			return
		}
	}

	// appendKey, it couldn't be updated.
	c.appendKey(keyHash, id, value)

	// We unlock the shard
	c.shards[id].lock.Unlock()
}

// appendKey is used to append a key to the cache. This is called by both
// set and get commands. This should only be called when a shard is already unlocked.
func (c *SpearCache) appendKey(keyHash uint64, id uint64, value interface{}) {

	// if the cursor of the queue is longer than defaultItems - 1, the cursor is reset to zero
	if c.shards[id].cursor > defaultItems-1 {
		c.shards[id].cursor = 0
	}

	// Key value pairs are appended to the cache. SpearCache only updates an existing key
	// if that key is maximally 10 entries removed from the cursor. This is to prevent
	// cache cluttering. This makes SpearCache set commands very fast.
	// A cache get will always retrieve the latest key value, if the key exists and is not yet expired.
	c.shards[id].items[c.shards[id].cursor] = item{
		key:     keyHash,
		value:   value,
		modTime: uint64(time.Now().UnixNano()),
	}

	// We increase the cursor
	c.shards[id].cursor++

}

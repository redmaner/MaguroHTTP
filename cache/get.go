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

// Get is used to retrieve a key from the cache. Get requires the key and the max age
// of the key in nano seconds. If the key is found true and the value are returned.
// If the key is not found or is expired, false and nil are returned
func (c *SpearCache) Get(key string, maxAge uint64) (bool, interface{}) {

	// hash the key with xxhash and make the id
	keyHash := xxhash.Sum64([]byte(key))
	id := keyHash & (defaultShards - 1)

	// Get current time in nano seconds
	now := uint64(time.Now().UnixNano())

	// We make sure the shard exists, if it doesn't the key isn't stored
	if c.shards[id] == nil {
		return false, nil
	}

	// Lock the shard for concurrency safetey
	c.shards[id].lock.Lock()

	rangeEnd := c.shards[id].cursor - defaultItems - 1
	var itemsParsed int

	// We range over the ring queue
	for i := c.shards[id].cursor; i >= rangeEnd; i-- {
		itemsParsed++
		itemID := i

		if itemID < 0 {
			itemID = itemID + defaultItems
		}
		if itemID >= defaultItems {
			itemID = defaultItems - 1
		}

		key := c.shards[id].items[itemID].key

		// If the key is empty we continue, this is when the queue is empty at the start
		if key == 0 {
			continue
		}

		if now-c.shards[id].items[itemID].modTime > maxAge && itemsParsed < defaultNoClutter {
			break
		}

		// We test if the item is within the maxAge range. If it is not we break.
		// If it is we check if the key matches the key we search.
		if key == keyHash {

			// appendKey
			c.appendKey(key, id, c.shards[id].items[itemID].value)

			// Unlock and return
			c.shards[id].lock.Unlock()
			return true, c.shards[id].items[itemID].value
		}
	}

	c.shards[id].lock.Unlock()

	return false, nil
}

package cache

import (
	"time"

	"github.com/cespare/xxhash"
)

// Set is used to set a key value pair into TinyCache
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

	// appendKey
	c.appendKey(keyHash, id, value)

	// We unlock the shard
	c.shards[id].lock.Unlock()
}

func (c *SpearCache) appendKey(keyHash uint64, id uint64, value interface{}) {

	// if the cursor of the queue is longer than defaultItems - 1, the cursor is reset to zero
	if c.shards[id].cursor > defaultItems-1 {
		c.shards[id].cursor = 0
	}

	// All key value pairs are appended to the queue. All keys are in time order.
	// Keys that already exist in the queue are not updated, but appended. A set in
	// SpearCache is therefore very fast.
	// A cache get will retrieve the latest key, if it exists and is not yet expired.
	c.shards[id].items[c.shards[id].cursor] = item{
		key:     keyHash,
		value:   value,
		modTime: uint32(time.Now().Unix()),
	}

	// We increase the cursor
	c.shards[id].cursor++

}

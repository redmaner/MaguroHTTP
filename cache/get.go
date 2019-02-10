package cache

import (
	"errors"
	"time"

	"github.com/cespare/xxhash"
)

var (
	errCoverOutOfRange = errors.New("Cover should be between 10 and 100")
)

func (c *SpearCache) Get(key string, maxAge uint32) (bool, interface{}) {

	// hash the key with xxhash and make the id
	keyHash := xxhash.Sum64([]byte(key))

	if val := c.getKey(keyHash, maxAge, 100); val != nil {
		return true, val
	}

	return false, nil
}

func (c *SpearCache) GetFast(key string, maxAge uint32, cover int) (bool, interface{}, error) {

	// hash the key with xxhash and make the id
	keyHash := xxhash.Sum64([]byte(key))

	if cover < 10 || cover > 100 {
		return false, nil, errCoverOutOfRange
	}

	if val := c.getKey(keyHash, maxAge, cover); val != nil {
		return true, val, nil
	}

	return false, nil, nil
}

func (c *SpearCache) getKey(keyHash uint64, maxAge uint32, cover int) interface{} {

	// get shard ID
	id := keyHash & (defaultShards - 1)

	now := uint32(time.Now().Unix())

	// We make sure the shard exists, if it doesn't the key isn't stored
	if c.shards[id] == nil {
		return nil
	}

	// Lock the shard for concurrency safetey
	c.shards[id].lock.Lock()

	// we determine the range end and start.
	// A default Get will cover 100% of the queue. This can be very costly.
	// The GetFast command allows to set a smaller coverage below 100%
	rangeEnd := defaultItems * cover / 100
	rangeStart := c.shards[id].cursor - rangeEnd

	// We range over the ring queue
	for i := c.shards[id].cursor; i >= rangeStart; i-- {

		itid := i

		if itid < 0 {
			itid = itid + defaultItems
		}
		if itid >= defaultItems {
			itid = defaultItems - 1
		}

		key := c.shards[id].items[itid].key

		// If the key is empty we continue, this is when the queue is empty at the start
		if key == 0 {
			continue
		}

		if now-c.shards[id].items[itid].modTime > maxAge {
			break
		}

		// We test if the item is within the maxAge range. If it is not we break.
		// If it is we check if the key matches the key we search.
		if key == keyHash {

			// appendKey
			c.appendKey(key, id, c.shards[id].items[itid].value)

			// Unlock and return
			c.shards[id].lock.Unlock()
			return c.shards[id].items[itid].value
		}
	}

	c.shards[id].lock.Unlock()

	return nil
}

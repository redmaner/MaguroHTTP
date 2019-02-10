package cache

import (
	"sync"
)

// TinyCache is divided into shards
type shard struct {
	lock   sync.Mutex
	items  [defaultItems]item
	cursor int
}

// Each shard contains an array (yes array not slice) of item
type item struct {
	modTime uint32
	key     uint64
	value   interface{}
}

func newShard() *shard {
	return &shard{}
}

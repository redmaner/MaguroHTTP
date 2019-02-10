package cache

const (

	// defaultShards is the amount of shards in SpearCache
	defaultShards = 256

	// defaultItems is the amount of items in a single shard
	defaultItems = 1024
)

// SpearCache is a preallocated in memory cache. SpearCache uses a ring queue of a fixed length.
// All entries in the cache are appended in time order. Entries of the same key don't get updated,
// but appended instead. When the ring is full, the oldest entries are automatically overwritten.
// A cache get will retrieve the newest appended entry to the queue, if it exist and is not yet expired.
type SpearCache struct {
	shards [defaultShards]*shard
}

func NewCache() *SpearCache {
	return &SpearCache{}
}

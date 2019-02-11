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
	"sync"
)

// SpearCache is divided into shards. This allows multiple go routines to access,
// read and write the cache concurrently while not being limited by a single lock.
type shard struct {
	lock   sync.Mutex
	items  [defaultItems]item
	cursor int
}

// Each shard contains an array (yes array not slice) of item.
type item struct {
	modTime uint64
	key     uint64
	value   interface{}
}

// newShard returns an empty pointer to a shard
func newShard() *shard {
	return &shard{}
}

package kvstore

import (
	"rdb/storage/b_plus_tree"
)

func (db *KV) Get(key []byte) ([]byte, bool) {
    return db.tree.Get(key)
}
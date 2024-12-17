package kvstore

import (
	"rdb/storage/b_plus_tree"
)

func(db *KV) pageWrite(ptr uint64) []byte{
	if node, ok := db.page.updates[ptr]; ok {
		return node
	}
	node := make([]byte, b_plus_tree.BTREE_PAGE_SIZE)
	copy(node, db.pageReadFile(ptr)) // initialized from the file
	db.page.updates[ptr] = node
	return node
}
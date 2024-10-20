package kvstore

import (
	"rdb/storage/b_plus_tree"
)

func (db *KV) pageRead(ptr uint64) []byte {
    start := uint64(0)
    for _, chunk := range db.mmap.chunks {
        end := start + uint64(len(chunk))/b_plus_tree.BTREE_PAGE_SIZE
        if ptr < end {
            offset := b_plus_tree.BTREE_PAGE_SIZE * (ptr - start)
            return chunk[offset : offset+b_plus_tree.BTREE_PAGE_SIZE]
        }
        start = end
    }
    panic("bad ptr")
}
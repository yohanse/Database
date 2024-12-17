package kvstore

import (
	"rdb/storage/b_plus_tree"
)

// `BTree.get`, read a page

// func pageRead(p *page, n uint64) ([]byte, error)
//     The pageRead function retrieves a page of data from the memory-mapped chunks based on a page 
//     number pointer

func (db *KV) pageReadFile(ptr uint64) []byte {
    start := uint64(0)
    for _, chunk := range db.mmap.chunks {
        // Each chunk represents a portion of the memory-mapped file.
        end := start + uint64(len(chunk))/b_plus_tree.BTREE_PAGE_SIZE
        if ptr < end {
            offset := b_plus_tree.BTREE_PAGE_SIZE * (ptr - start)
            return chunk[offset : offset+b_plus_tree.BTREE_PAGE_SIZE]
        }
        start = end
    }
    panic("bad ptr")
}
package kvstore

import (
	"rdb/storage/b_plus_tree"
)

type KV struct {
	path string
	fd int
	tree b_plus_tree.BTree

    mmap struct {
        total  int      // mmap size, can be larger than the file size
        chunks [][]byte // multiple mmaps, can be non-continuous
    }

    page struct {
        flushed uint64   // database size in number of pages
        temp    [][]byte // newly allocated pages
    }
}
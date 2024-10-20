package kvstore

import (
	"rdb/storage/b_plus_tree"
)

// KV Struct

// path: The file path where the database is stored.

// fd: The file descriptor associated with the database file. This allows for direct I/O operations 
// on the file.

// tree: The B+ tree structure that organizes your data.

// mmap struct:
//     total: The total size of the memory-mapped file. This size can be larger than the actual file 
//     size, which allows for extending the file in the future.

//     chunks: A slice of byte slices ([][]byte). This means that you are mapping the file into memory 
//     in chunks, rather than mapping the entire file continuously. Each chunk represents a part of 
//     the file that has been memory-mapped.

// page struct:
//     flushed: Tracks how many pages of the database have been flushed to disk. This is likely used to 
//     manage which parts of the B+ tree have been persisted.

//     temp: A slice of byte slices ([][]byte). These are newly allocated pages in memory that have not 
//     yet been flushed to disk.


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
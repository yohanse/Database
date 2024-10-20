package rdb

import (
    "fmt"
    "testing"
    "unsafe"
	"rdb/storage/b_plus_tree"
)

// Assuming BNode and BTree are defined similarly to your original code.

type C struct {
    tree  b_plus_tree.BTree
    ref   map[string]string // the reference data
    pages map[uint64]b_plus_tree.BNode  // in-memory pages
}

func newC() *C {
    pages := map[uint64]b_plus_tree.BNode{}
    return &C{
        tree: b_plus_tree.BTree{
            get: func(ptr uint64) []byte {
                node, ok := pages[ptr]
                assert(ok)
                return node
            },
            new: func(node []byte) uint64 {
                assert(BNode(node).nbytes() <= BTREE_PAGE_SIZE)
                ptr := uint64(uintptr(unsafe.Pointer(&node[0])))
                assert(pages[ptr] == nil)
                pages[ptr] = node
                return ptr
            },
            del: func(ptr uint64) {
                assert(pages[ptr] != nil)
                delete(pages, ptr)
            },
        },
        ref:   map[string]string{},
        pages: pages,
    }
}

func (c *C) add(key string, val string) {
    c.tree.Insert([]byte(key), []byte(val))
    c.ref[key] = val // reference data
}

func (c *C) testInsertion() {
    keys := []string{"apple", "banana", "cherry", "date", "elderberry"}
    vals := []string{"red", "yellow", "red", "brown", "purple"}

    for i := range keys {
        c.add(keys[i], vals[i])
    }

    for i := range keys {
        if c.ref[keys[i]] != vals[i] {
            panic(fmt.Sprintf("Insertion failed for key: %s", keys[i]))
        }
    }
}

func (c *C) testDeletion() {
    // Initial insertion of key-value pairs
    c.add("apple", "red")
    c.add("banana", "yellow")

    // Delete key "apple"
    success, err := c.tree.Delete([]byte("apple"))
    if !success || err != nil {
        panic("Deletion failed")
    }

    // Check if the key is removed from the reference map
    _, found := c.ref["apple"]
    if found {
        panic("Deletion not reflected in reference data")
    }
}

func (c *C) testLookup() {
    // Insert key-value pairs
    c.add("apple", "red")
    c.add("banana", "yellow")

    // Look up a key in the tree
    value := c.tree.Lookup([]byte("banana"))
    if string(value) != "yellow" {
        panic("Lookup failed for key: banana")
    }
}

func TestBTreeOperations(t *testing.T) {
    c := newC()

    // Test Insertion
    c.testInsertion()

    // Test Deletion
    c.testDeletion()

    // Test Lookup
    c.testLookup()
}

package btree

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
            Get: func(ptr uint64) []byte {
                node, ok := pages[ptr]
                if !ok{
					panic(fmt.Sprintf("Page not found: %d", ptr))
				}
                return node
            },
            New: func(node []byte) uint64 {
                
                ptr := uint64(uintptr(unsafe.Pointer(&node[0])))
                
                pages[ptr] = node
                return ptr
            },
            Del: func(ptr uint64) {
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
    keys := []string{"apple", "banana", "cherry", "date", "elderberry", "fig"}
    vals := []string{"red", "yellow", "red", "yellow", "red", "yellow"}

    for i := range keys {
        c.add(keys[i], vals[i])
        fmt.Println(keys[i], vals[i])
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

    success1, err1 := c.tree.Delete([]byte("banana"))
    if !success1 || err1 != nil {
        panic("Deletion failed")
    }
    
}

func TestBTreeOperations(t *testing.T) {
    c := newC()

    // // Test Insertion
    // c.testInsertion()

    // Test Deletion
    c.testDeletion()
}

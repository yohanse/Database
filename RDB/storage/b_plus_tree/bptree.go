package b_plus_tree

type BTree struct {
    // pointer (a nonzero page number)
    Root uint64
    // callbacks for managing on-disk pages
    Get func(uint64) []byte // dereference a pointer
    New func([]byte) uint64 // allocate a new page
    Del func(uint64)        // deallocate a page
}
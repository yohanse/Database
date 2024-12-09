package b_plus_tree

type BTree struct {
    // pointer (a nonzero page number)
    Root uint64
    // callbacks for managing on-disk pages
    Get func(uint64) []byte // reads a page from disk
    New func([]byte) uint64 // allocates and writes a new page(copy on write)
    Del func(uint64)        // deallocate a page
}
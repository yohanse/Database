package kvstore

// pageAppend adds a new page (node) to the in-memory storage and returns its index pointer.

func (db *KV) pageAppend(node []byte) uint64 {
    ptr := db.page.flushed + uint64(len(db.page.temp)) // just append
    db.page.temp = append(db.page.temp, node)
    return ptr
}
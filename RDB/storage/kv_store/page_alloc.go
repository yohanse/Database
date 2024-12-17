package kvstore

func(db *KV) pageAllocation(node []byte) uint64 {
	if ptr := db.free.PopHead(); ptr != 0 { // try the free list
		db.page.updates[ptr] = node
		return ptr
	}
	return db.pageAppend(node) // append
}
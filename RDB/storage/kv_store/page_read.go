package kvstore

// `BTree.get`, read a page

// func pageRead(p *page, n uint64) ([]byte, error)
//     The pageRead function retrieves a page of data from the memory-mapped chunks based on a page 
//     number pointer

func (db *KV) pageRead(ptr uint64) []byte {
    if node, ok := db.page.updates[ptr]; ok {
        return node // pending update
        }
    return db.pageReadFile(ptr)
}
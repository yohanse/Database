package kvstore


// func (db *KV) Get(key []byte) ([]byte, bool)
//         This function retrieves the value associated with a given key from the B+ tree structure.

func (db *KV) Get(key []byte) ([]byte, bool) {
    return db.tree.Get(key)
}
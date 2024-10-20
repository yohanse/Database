package kvstore

func (db *KV) Get(key []byte) ([]byte, bool) {
    return db.tree.Get(key)
}
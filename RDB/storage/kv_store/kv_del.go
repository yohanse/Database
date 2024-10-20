package kvstore

func (db *KV) Del(key []byte) (bool, error) {
    deleted, _ := db.tree.Delete(key)
    return deleted, updateFile(db)
}
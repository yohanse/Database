package kvstore

// he Del method in your KV struct is intended to delete a key-value pair from the B+ tree and 
// then update the database file to reflect this change. 
func (db *KV) Del(key []byte) (bool, error) {
    deleted, _ := db.tree.Delete(key)
    return deleted, updateFile(db)
}
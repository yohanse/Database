package kvstore


// The Set function in the KV struct is responsible for inserting a key-value pair into the B+ tree 
// and ensuring that the changes are persisted to the underlying storage (file on disk). 

func (db *KV) Set(key []byte, val []byte) error {
    meta := saveMeta(db) // save the in-memory state (tree root)
	err := db.tree.Insert(key, val)
    if err != nil {
        return err // length limit
    }
    return updateOrRevert(db, meta)
}
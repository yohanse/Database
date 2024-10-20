package kvstore

// updateOrRevert attempts to update the database file. 
// If the update fails, it reverts the in-memory state using the provided metadata 
// and clears temporary data. Returns an error if the update operation fails.

func updateOrRevert(db *KV, meta []byte) error {
    // 2-phase update
    err := updateFile(db)
    // revert on error
    if err != nil {
        // the in-memory states can be reverted immediately to allow reads
        loadMeta(db, meta)
        // discard temporaries
        db.page.temp = db.page.temp[:0]
    }
    return err
}
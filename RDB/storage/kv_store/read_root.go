package kvstore

// readRoot reads the metadata from the memory-mapped area, initializing the KV state for an empty 
// file or loading the root metadata from the first page for non-empty files.

func readRoot(db *KV, fileSize int64) error {
    if fileSize == 0 { // empty file
        db.page.flushed = 1 // the meta page is initialized on the 1st write
        db.free.HeadPage = 1 // the 2nd page
        db.free.TailPage = 1
        return nil
    }
    // read the page
    data := db.mmap.chunks[0]
    loadMeta(db, data)
    // verify the page
    // ...
	// ...
    return nil
}
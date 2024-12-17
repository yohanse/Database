package kvstore

// func (db *KV) Open() error
//     Usage: When a database application starts, it needs to connect to the database file (or create 
//     it if it doesn't exist).

func (db *KV) Open() error { // open or create
    db.tree.Get = db.pageRead   // read a page
    db.tree.New = db.pageAppend // apppend a page
    db.tree.Del = func(uint64) {}

    db.free.Get = db.pageRead      // read a page
    db.free.New = db.pageAppend    // append a page
    db.free.Set = db.pageWrite  
	return nil
}
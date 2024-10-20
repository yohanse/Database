package kvstore

func (db *KV) Open() error { // open or create
    db.tree.Get = db.pageRead   // read a page
    db.tree.New = db.pageAppend // apppend a page
    db.tree.Del = func(uint64) {}
}
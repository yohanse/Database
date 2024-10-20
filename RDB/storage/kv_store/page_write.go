package kvstore

import (
	"golang.org/x/sys/unix"
	"rdb/storage/b_plus_tree"
)

// writePages writes the in-memory pages to the file, extending the mmap if needed and updating the 
// flushed page count.


func writePages(db *KV) error {
    // extend the mmap if needed
    size := (int(db.page.flushed) + len(db.page.temp)) * b_plus_tree.BTREE_PAGE_SIZE
    if err := extendMmap(db, size); err != nil {
        return err
    }
    // write data pages to the file
    offset := int64(db.page.flushed * b_plus_tree.BTREE_PAGE_SIZE)
    if _, err := unix.Pwritev(db.fd, db.page.temp, offset); err != nil {
        return err
    }
    // discard in-memory data
    db.page.flushed += uint64(len(db.page.temp))
    db.page.temp = db.page.temp[:0]
    return nil
}
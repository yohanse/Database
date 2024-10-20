package kvstore

import (
	"fmt"
	"syscall"
)

// updateRoot writes the current metadata of the database to the meta page at offset 0.
// It returns an error if the write operation fails.

func updateRoot(db *KV) error {
    if _, err := syscall.Pwrite(db.fd, saveMeta(db), 0); err != nil {
        return fmt.Errorf("write meta page: %w", err)
    }
    return nil
}
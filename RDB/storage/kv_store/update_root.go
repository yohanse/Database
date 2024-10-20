package kvstore

import (
	"fmt"
	"syscall"
)

func updateRoot(db *KV) error {
    if _, err := syscall.Pwrite(db.fd, saveMeta(db), 0); err != nil {
        return fmt.Errorf("write meta page: %w", err)
    }
    return nil
}
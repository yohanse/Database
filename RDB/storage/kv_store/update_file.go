package kvstore

import (
	"syscall"
)

// updateFile writes new nodes to the disk, synchronizes the file to ensure data integrity,
// updates the root pointer atomically, and flushes changes to make them persistent.

func updateFile(db *KV) error {
    // 1. Write new nodes.
    if err := writePages(db); err != nil {
        return err
    }
    // 2. `fsync` to enforce the order between 1 and 3.
    if err := syscall.Fsync(db.fd); err != nil {
        return err
    }
    // 3. Update the root pointer atomically.
    if err := updateRoot(db); err != nil {
        return err
    }
    // 4. `fsync` to make everything persistent.
    return syscall.Fsync(db.fd)
}
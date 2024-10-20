package kvstore

import (
	"fmt"
	"syscall"
)
// updateOrRevert attempts to update the database file. 
// If the update fails, it reverts the in-memory state using the provided metadata 
// and clears temporary data. Returns an error if the update operation fails.

func updateOrRevert(db *KV, meta []byte) error {
    // ensure the on-disk meta page matches the in-memory one after an error
    if db.failed {
        // write and fsync the previous meta page
        if _, err := syscall.Pwrite(db.fd, meta, 0); err != nil {
            return fmt.Errorf("write previous meta page: %w", err)
        }
        if err := syscall.Fsync(db.fd); err != nil {
            return fmt.Errorf("fsync previous meta page: %w", err)
        }
        db.failed = false // reset the failure state after writing
    }

    // Attempt to update the file with new data
    err := updateFile(db)
    if err != nil {
        // the on-disk meta page is in an unknown state;
        // mark it to be rewritten on later recovery.
        db.failed = true
        return fmt.Errorf("update file failed: %w", err) // Return error after marking failed
    }

    return nil // Successful update
}

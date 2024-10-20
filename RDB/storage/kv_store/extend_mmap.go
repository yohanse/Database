package kvstore

import (
	"fmt"
	"syscall"
)

// extendMmap ensures that the memory-mapped region for the database is extended
// as needed to accommodate the required size. It checks if the current mapped
// size is sufficient; if not, it calculates a new allocation size (doubling
// as necessary) and maps the new memory region into the process's address space.
// This allows for efficient access to database pages without excessive disk I/O.
// The function updates the total mapped size and maintains a record of all mapped
// chunks for future reference.

func extendMmap(db *KV, size int) error {
    if size <= db.mmap.total {
        return nil // enough range
    }
    alloc := max(db.mmap.total, 64<<20) // double the current address space
    for db.mmap.total + alloc < size {
        alloc *= 2 // still not enough?
    }
    chunk, err := syscall.Mmap(
        db.fd, int64(db.mmap.total), alloc,
        syscall.PROT_READ, syscall.MAP_SHARED, // read-only
    )
    if err != nil {
        return fmt.Errorf("mmap: %w", err)
    }
    db.mmap.total += alloc
    db.mmap.chunks = append(db.mmap.chunks, chunk)
    return nil
}
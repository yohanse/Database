package kvstore

import (
	"encoding/binary"
)

func loadMeta(db *KV, data []byte) {
    // Read the signature
    if string(data[:16]) != DB_SIG {
        panic("Invalid database signature")
    }

    // Read the root pointer
    db.tree.Root = binary.LittleEndian.Uint64(data[16:24])

    // Read the page usage count
    db.page.flushed = binary.LittleEndian.Uint64(data[24:])
}
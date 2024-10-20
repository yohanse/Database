package kvstore

import (
	"encoding/binary"
)

func loadMeta(db *KV, data []byte) {
	db.tree.Root = binary.LittleEndian.Uint64(data[16:])
	db.page.flushed = binary.LittleEndian.Uint64(data[24:])
}
package kvstore

import (
	"encoding/binary"
)

// | sig | root_ptr | page_used |
// | 16B |    8B    |     8B    |
func saveMeta(db *KV) []byte {
    var data [32]byte
    copy(data[:16], []byte(DB_SIG))
    binary.LittleEndian.PutUint64(data[16:], db.tree.Root)
    binary.LittleEndian.PutUint64(data[24:], db.page.flushed)
    return data[:]
}
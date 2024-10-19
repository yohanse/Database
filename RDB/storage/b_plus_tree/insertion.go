package rdb

import (
	"encoding/binary"
)

func leafInsert(new BNode, old BNode, idx uint16, key []byte, val []byte) {
	new.setHeaders(BNODE_LEAF, old.nkeys() + 1)
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKV(new, idx, 0, key, val)
	nodeAppendRange(new, old, idx+1, idx, old.nkeys()-idx)
}

func nodeAppendKV(new BNode, idx uint16, ptr uint64, key []byte, val []byte) {
	new.setPtr(idx, ptr)

	pos := new.KvPos(idx)

	binary.LittleEndian.AppendUint16(new[pos:], uint16(len(key)))
	binary.LittleEndian.AppendUint16(new[pos + 2:], uint16(len(val)))

	copy(new[pos + 4:], key)
	copy(new[pos + 4 + uint16(len(key)):], val)

	new.setOffset(idx + 1, new.getOffSet(idx) + 4 + uint16((len(key)+len(val))))
}

func nodeAppendRange(new BNode, old BNode, dstNew uint16, srcOld uint16, n uint16) {
	for i := uint16(0); i < n; i++ {
		new.setPtr(dstNew, old.getPtr(srcOld))

		pos := new.KvPos(dstNew)

		key := old.getKey(srcOld)
		val := old.getVal(srcOld)

		binary.LittleEndian.PutUint16(new[pos+0:], uint16(len(key)))
		binary.LittleEndian.PutUint16(new[pos+2:], uint16(len(val)))
		copy(new[pos+4:], key)
		copy(new[pos+4+uint16(len(key)):], val)

		new.setOffset(dstNew + 1, new.getOffSet(dstNew) + 4 + uint16((len(key)+len(val))))

		dstNew += 1
		srcOld += 1
	}
}

// func leafUpdate(new BNode, old BNode, idx uint16, key []byte, val []byte) {
// 	new.setHeaders(BNODE_LEAF, old.nkeys() + 1)
// 	nodeAppendRange(new, old, 0, 0, idx)
// 	nodeAppendKV(new, idx, 0, key, val)
// 	nodeAppendRange(new, old, idx+1, idx+1, old.nkeys()-idx)
// }
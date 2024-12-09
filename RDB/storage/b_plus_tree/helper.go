package b_plus_tree

import (
	"encoding/binary"
)

// TYPE AND NKEYS
func (node BNode) btype() uint16 {
	return binary.LittleEndian.Uint16(node[0:2])
}

func (node BNode) nkeys() uint16 {
	return binary.LittleEndian.Uint16(node[2:4])
}

func (node BNode) setHeaders(btype uint16, nkeys uint16) {
	binary.LittleEndian.PutUint16(node[0:2], btype)
	binary.LittleEndian.PutUint16(node[2:4], nkeys)
}


// PTRS
func (node BNode) getPtr(idx uint16) uint64{
	if idx >= node.nkeys() {
		panic("index out of bounds")
	}
	pos := HEADER + 8*idx
	return binary.LittleEndian.Uint64(node[pos:])
}

func (node BNode) setPtr(idx uint16, val uint64) {
	if idx >= node.nkeys() {
		panic("index out of bounds")
	}
	pos := HEADER + 8*idx
	binary.LittleEndian.PutUint64(node[pos:], val)
}

// OFFSETS

// offsetPos returns the position of the nth offset relative to the whole node.
func offsetPos(node BNode, idx uint16) uint16 {
	if idx > node.nkeys() {
		panic("index out of range")
	}
	return HEADER + 8*node.nkeys() + 2*(idx - 1)
}

func (node BNode) getOffSet(idx uint16) uint16 {
	if (idx == 0) {
		return 0
	}
	return binary.LittleEndian.Uint16(node[offsetPos(node, idx):])
}

func (node BNode) setOffset(idx uint16, val uint16) {
	binary.LittleEndian.PutUint16(node[offsetPos(node, idx):], val)
}

// KV PAIRS

// kvPos returns the position of the nth KV pair relative to the whole node.
func (node BNode) KvPos(idx uint16) uint16 {
	if idx > node.nkeys() {
		panic("index out of range")
	}
	return HEADER + 8*node.nkeys() + 2*node.nkeys() + node.getOffSet(idx)
}

func (node BNode) getKey(idx uint16) []byte {
	if idx >= node.nkeys() {
		panic("index out of range")
	}

	pos := node.KvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:])
	return node[pos + 4:][:klen]
}

func (node BNode) getVal(idx uint16) []byte {
	if idx >= node.nkeys() {
		panic("index out of range")
	}

	pos := node.KvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:])
	vlen := binary.LittleEndian.Uint16(node[pos + 2:])

	return node[pos + 4 + klen:][:vlen]
}

// It also conveniently returns the node size (used space) with an off-by-one lookup.
func (node BNode) nbytes() uint16 {
	return node.KvPos(node.nkeys())
}
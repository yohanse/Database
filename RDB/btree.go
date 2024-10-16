package rdb

import (
	"encoding/binary"
)

const HEADER = 4
const BTREE_PAGE_SIZE = 4096
const BTREE_MAX_KEY_SIZE = 1000
const BTREE_MAX_VAL_SIZE = 3000

type BNode []byte

type BTree struct {
	root uint64
	get func (uint64) []byte
	new func ([]byte) uint64
	del func (uint64)
}

const (
	BNODE_NODE = 1
	BNODE_LEAF = 2
)

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

func offsetPos(node BNode, idx uint16) uint16 {
	if idx < 1 || idx > node.nkeys() {
		panic("index out of range")
	}

	return HEADER + 8*node.nkeys() + 2*(idx - 1)
}

func (node BNode) getOffSet(idx uint16) uint16 {
	return binary.LittleEndian.Uint16(node[offsetPos(node, idx):])
}

func (node BNode) setOffset(idx uint16, val uint16) {
	binary.LittleEndian.PutUint16(node[offsetPos(node, idx):], val)
}

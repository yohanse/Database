package rdb

import (
	"bytes"
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

func (node BNode) KvPos(idx uint16) uint16 {
	if idx < 1 || idx > node.nkeys() {
		panic("index out of range")
	}
	return HEADER + 8*node.nkeys() + 2*node.nkeys() + node.getOffSet(idx)
}

func (node BNode) getKey(idx uint16) []byte {
	if idx < 1 || idx > node.nkeys() {
		panic("index out of range")
	}

	pos := node.KvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:])
	return node[pos + 4:][:klen]
}

func (node BNode) getVal(idx uint16) []byte {
	if idx < 1 || idx > node.nkeys() {
		panic("index out of range")
	}

	pos := node.KvPos(idx)
	klen := binary.LittleEndian.Uint16(node[pos:])
	vlen := binary.LittleEndian.Uint16(node[pos + 2:])

	return node[pos + 4 + klen:][:vlen]
}

func (node BNode) nbytes() uint16 {
	return node.KvPos(node.nkeys())
}


func nodeLookupLE(node BNode, key []byte) uint16 {
	nkeys := node.nkeys()
	found := uint16(0)

	for i := uint16(1); i < nkeys; i++ {
		cmp := bytes.Compare(node.getKey(i), key)

		if cmp > 0 {
			break
		}
		found = i
	}
	return found
}


func leafInsert(new BNode, old BNode, idx uint16, key []byte, val []byte) {
	new.setHeaders(BNODE_LEAF, old.nkeys() + 1)
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKV(new, idx, 0, key, val)
	nodeAppendRange(new, old, idx+1, idx, old.nkeys()-idx)
}

func leafUpdate(new BNode, old BNode, idx uint16, key []byte, val []byte) {
	new.setHeaders(BNODE_LEAF, old.nkeys() + 1)
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendKV(new, idx, 0, key, val)
	nodeAppendRange(new, old, idx+1, idx+1, old.nkeys()-idx)
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

func nodeReplaceKidN(tree *BTree, new BNode, old BNode, idx uint16, kids ...BNode) {
	inc := uint16(len(kids))

	new.setHeaders(BNODE_NODE, old.nkeys() + inc - 1)
	nodeAppendRange(new, old, 0, 0, idx)
	for i, node := range kids {
		nodeAppendKV(new, idx+uint16(i), tree.new(node), node.getKey(0), nil)
	}
	nodeAppendRange(new, old, idx+inc, idx+1, old.nkeys()-(idx+1))
}

func nodeSplit2(left BNode, right BNode, old BNode) {
	nkeys := old.nkeys()
	mid := nkeys / 2

	nodeAppendRange(left, old, 0, 0, mid)
	nodeAppendRange(right, old, 0, mid, nkeys - mid)
}

func nodeSplit3(old BNode) (uint16, [3]BNode) {
	if old.nbytes() <= BTREE_PAGE_SIZE {
		old = old[:BTREE_PAGE_SIZE]
		return 1, [3]BNode{old}
	}

	left := BNode(make([]byte, 2*BTREE_PAGE_SIZE))
	right := BNode(make([]byte, BTREE_PAGE_SIZE))

	nodeSplit2(left, right, old)

	if left.nbytes() <= BTREE_PAGE_SIZE {
		left = left[:BTREE_PAGE_SIZE]
		return 2, [3]BNode{left, right}
	}

	leftleft := BNode(make([]byte, 2*BTREE_PAGE_SIZE))
	middle := BNode(make([]byte, BTREE_PAGE_SIZE))

	nodeSplit2(leftleft, middle, left)
	return 3, [3]BNode{leftleft, middle, right}
}

func treeInsert(tree *BTree, node BNode, key []byte, val []byte) BNode {
	new := BNode(make([]byte, 2*BTREE_PAGE_SIZE))

	idx := nodeLookupLE(node, key)
	switch node.btype() {
	case BNODE_LEAF:
	if bytes.Equal(key, node.getKey(idx)) {
		leafUpdate(new, node, idx, key, val)
	} else {
		leafInsert(new, node, idx+1, key, val)
	}
	case BNODE_NODE:
		nodeInsert(tree, new, node, idx, key, val)
	default:
		panic("bad node!")
	}
	return new
}

func nodeInsert(tree *BTree, new BNode, node BNode, idx uint16, key []byte, val []byte) {
	kptr := node.getPtr(idx)
	knode := treeInsert(tree, tree.get(kptr), key, val)
	nsplit, split := nodeSplit3(knode)
	tree.del(kptr)
	nodeReplaceKidN(tree, new, node, idx, split[:nsplit]...)
}

func (tree *BTree) Insert(key []byte, val []byte) {
	if tree.root == 0 {
		root := BNode(make([]byte, BTREE_PAGE_SIZE))
		root.setHeaders(BNODE_LEAF, 1)

		nodeAppendKV(root, 0, 0, nil, nil)
		nodeAppendKV(root, 1, 0, key, val)
		tree.root = tree.new(root)
		return 
	}

	node := treeInsert(tree, tree.get(tree.root), key, val)
	nsplit, split := nodeSplit3(node)
	tree.del(tree.root)
	
	if nsplit > 1 {
		root := BNode(make([]byte, BTREE_PAGE_SIZE))
		root.setHeaders(BNODE_NODE, nsplit)
		for i, knode := range split[:nsplit] {
			ptr, key := tree.new(knode), knode.getKey(0)
			nodeAppendKV(root, uint16(i), ptr, key, nil)
		}
		tree.root = tree.new(root)
	} else {
		tree.root = tree.new(split[0])
	}
}

func (tree *BTree) Delete(key []byte) bool {

}

func shouldMerge(tree *BTree, node BNode, idx uint16, updated BNode) (int, BNode){
	if idx > 0 {
		sibling := BNode(tree.get(node.getPtr(idx-1)))
		merged := sibling.nbytes() + updated.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return -1, sibling
		}
	}

	if idx + 1 < node.nkeys() {
		sibling := BNode(tree.get(node.getPtr(idx+1)))
		merged := sibling.nbytes() + updated.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return 1, sibling
		}
	}
	return 0, BNode{}
}
	
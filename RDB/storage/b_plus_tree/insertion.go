package b_plus_tree

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

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
		nodeAppendKV(new, idx+uint16(i), tree.New(node), node.getKey(0), nil)
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
	knode := treeInsert(tree, tree.Get(kptr), key, val)
	nsplit, split := nodeSplit3(knode)
	tree.Del(kptr)
	nodeReplaceKidN(tree, new, node, idx, split[:nsplit]...)
}

func (tree *BTree) Insert(key []byte, val []byte) {
	if tree.Root == 0 {
		fmt.Println("Root is 0")
		root := BNode(make([]byte, BTREE_PAGE_SIZE))
		root.setHeaders(BNODE_LEAF, 2)
		
		nodeAppendKV(root, 0, 0, nil, nil)
		nodeAppendKV(root, 1, 0, key, val)
		tree.Root = tree.New(root)
		return 
	}
	fmt.Println("root", tree.Root)
	node := treeInsert(tree, tree.Get(tree.Root), key, val)
	nsplit, split := nodeSplit3(node)
	tree.Del(tree.Root)
	
	if nsplit > 1 {
		root := BNode(make([]byte, BTREE_PAGE_SIZE))
		root.setHeaders(BNODE_NODE, nsplit)
		for i, knode := range split[:nsplit] {
			ptr, key := tree.New(knode), knode.getKey(0)
			nodeAppendKV(root, uint16(i), ptr, key, nil)
		}
		tree.Root = tree.New(root)
	} else {
		tree.Root = tree.New(split[0])
	}
}
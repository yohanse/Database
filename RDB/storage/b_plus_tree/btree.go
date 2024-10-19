package rdb

import (
	"bytes"
)
type BTree struct {
    // pointer (a nonzero page number)
    root uint64
    // callbacks for managing on-disk pages
    get func(uint64) []byte // dereference a pointer
    new func([]byte) uint64 // allocate a new page
    del func(uint64)        // deallocate a page
}

// returns the first kid node whose range intersects the key. (kid[i] <= key)
func nodeLookupLE(node BNode, key []byte) uint16 {
	right := node.nkeys()
    left := uint16(0)

	for left < right {
		mid := (left + right + 1) / 2
		cmp := bytes.Compare(node.getKey(mid), key)
		if cmp > 0 {
			right = mid - 1
		} else {
			left = mid
		}
	}
	return left
}

func leafDelete(new BNode, old BNode, idx uint16)
// merge 2 nodes into 1
func nodeMerge(new BNode, left BNode, right BNode)
// replace 2 adjacent links with 1
func nodeReplace2Kid(new BNode, old BNode, idx uint16, ptr uint64, key []byte)

func (tree *BTree) Delete(key []byte) bool 

func shouldMerge(tree *BTree, node BNode, idx uint16, updated BNode) (int, BNode){
	if updated.nbytes() > BTREE_PAGE_SIZE/4 {
		return 0, BNode{}
	}

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

func treeDelete(tree *BTree, node BNode, key []byte) BNode {
	new := BNode(make([]byte, 2*BTREE_PAGE_SIZE))

	idx := nodeLookupLE(node, key)
	switch node.btype() {
	case BNODE_LEAF:
		leafDelete(new, node, idx)
	case BNODE_NODE:
		nodeDelete(tree, new, node, idx, key)
	default:
		panic("bad node!")
	}
	return new
}

func nodeDelete(tree *BTree, node BNode, idx uint16, key []byte) BNode {
	kptr := node.getPtr(idx)
	knode := tree.get(kptr)

	updated := treeDelete(tree, knode, key)

	if len(updated) == 0 {
		return BNode{}
	}

	tree.del(kptr)

	new := BNode(make([]byte, BTREE_PAGE_SIZE))

	mergeDir, sibling := shouldMerge(tree, node, idx, updated)

	switch {
		case mergeDir < 0:
			merged := BNode(make([]byte, BTREE_PAGE_SIZE))
			nodeMerge(merged, sibling, updated)
			tree.del(node.getPtr(idx-1))
			nodeReplace2Kid(new, node, idx - 1, tree.new(merged), merged.getKey(0))

		case mergeDir > 0:
			merged := BNode(make([]byte, BTREE_PAGE_SIZE))
			nodeMerge(merged, updated, sibling)
			tree.del(node.getPtr(idx+1))
			nodeReplace2Kid(new, node, idx, tree.new(merged), merged.getKey(0))
		
		case mergeDir == 0 && updated.nkeys() == 0:
			new.setHeaders(BNODE_NODE, 0)
		
		case mergeDir == 0 && updated.nkeys() > 0:
			nodeReplaceKidN(tree, new, node, idx, updated)
	}
	return new
}
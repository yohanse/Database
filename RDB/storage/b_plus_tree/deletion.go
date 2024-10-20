package b_plus_tree

import (
	"fmt"
	"bytes"
)

// remove a key from a leaf node
func leafDelete(new BNode, old BNode, idx uint16) {
	new.setHeaders(BNODE_LEAF, old.nkeys() - 1)
	nodeAppendRange(new, old, 0, 0, idx)
	nodeAppendRange(new, old, idx, idx+1, old.nkeys()-idx-1)
}

func leafCopy(new BNode, old BNode) {
	new.setHeaders(BNODE_LEAF, old.nkeys())
	nodeAppendRange(new, old, 0, 0, old.nkeys())
}

// merge 2 nodes into 1
func nodeMerge(new BNode, left BNode, right BNode) {
	new.setHeaders(BNODE_NODE, left.nkeys() + right.nkeys())
	nodeAppendRange(new, left, 0, 0, left.nkeys())
	nodeAppendRange(new, right, left.nkeys(), 0, right.nkeys())
}

// replace 2 adjacent links with 1
func nodeReplace2Kid(new BNode, old BNode, idx uint16, ptr uint64, key []byte) {
	// Set the header for the new node (node type + number of keys)
    new.setHeaders(BNODE_NODE, old.nkeys()-1)
    
    // Copy all key-value pairs from the old node before the index
    nodeAppendRange(new, old, 0, 0, idx)
    
    // Append the new pointer and key in place of the two old child pointers
    nodeAppendKV(new, idx, ptr, key, nil)
    
    // Copy the remaining key-value pairs after the two links being replaced
    nodeAppendRange(new, old, idx+1, idx+2, old.nkeys()-(idx+2))
}

func shouldMerge(tree *BTree, node BNode, idx uint16, updated BNode) (int, BNode){
	if updated.nbytes() > BTREE_PAGE_SIZE/4 {
		return 0, BNode{}
	}

	if idx > 0 {
		sibling := BNode(tree.Get(node.getPtr(idx-1)))
		merged := sibling.nbytes() + updated.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return -1, sibling
		}
	}

	if idx + 1 < node.nkeys() {
		sibling := BNode(tree.Get(node.getPtr(idx+1)))
		merged := sibling.nbytes() + updated.nbytes() - HEADER
		if merged <= BTREE_PAGE_SIZE {
			return 1, sibling
		}
	}
	return 0, BNode{}
}

// delete a key from the tree
func treeDelete(tree *BTree, node BNode, key []byte) BNode {
	new := BNode(make([]byte, BTREE_PAGE_SIZE))
	idx := nodeLookupLE(node, key)
	fmt.Println(idx, node.nkeys())
	fmt.Println(node.btype())
	switch node.btype() {
		case BNODE_LEAF:
			fmt.Println("Leaf", string(node.getKey(idx)), string(key))
			if idx < node.nkeys() && bytes.Equal(key, node.getKey(idx)) {
				leafDelete(new, node, idx)
			} else {
				leafCopy(new, node)
			}
		case BNODE_NODE:
			nodeDelete(tree, node, idx, key)
		default:
			panic("bad node!")
	}
	return new
}


// delete a key from an internal node; part of the treeDelete()
func nodeDelete(tree *BTree, node BNode, idx uint16, key []byte) BNode {
	kptr := node.getPtr(idx)
	knode := tree.Get(kptr)

	updated := treeDelete(tree, knode, key)

	if len(updated) == 0 {
		return BNode{}
	}

	tree.Del(kptr)

	new := BNode(make([]byte, BTREE_PAGE_SIZE))

	mergeDir, sibling := shouldMerge(tree, node, idx, updated)

	switch {
		case mergeDir < 0:
			merged := BNode(make([]byte, BTREE_PAGE_SIZE))
			nodeMerge(merged, sibling, updated)
			tree.Del(node.getPtr(idx-1))
			nodeReplace2Kid(new, node, idx - 1, tree.New(merged), merged.getKey(0))

		case mergeDir > 0:
			merged := BNode(make([]byte, BTREE_PAGE_SIZE))
			nodeMerge(merged, updated, sibling)
			tree.Del(node.getPtr(idx+1))
			nodeReplace2Kid(new, node, idx, tree.New(merged), merged.getKey(0))
		
		case mergeDir == 0 && updated.nkeys() == 0:
			new.setHeaders(BNODE_NODE, 0)
		
		case mergeDir == 0 && updated.nkeys() > 0:
			nodeReplaceKidN(tree, new, node, idx, updated)
	}
	return new
}

func (tree *BTree) Delete(key []byte) (bool, error) {
    if tree.Root == 0 {
        return false, fmt.Errorf("tree is empty")
    }
    
    updatedRoot := treeDelete(tree, tree.Get(tree.Root), key)
    
    if len(updatedRoot) == 0 {
        root := BNode(make([]byte, BTREE_PAGE_SIZE))
        root.setHeaders(BNODE_LEAF, 0)
        tree.Root = tree.New(root)
    } else {
        // If the updated root is not empty, set it as the new root
        tree.Root = tree.New(updatedRoot)
    }

    return true, nil
}

package b_plus_tree

import "fmt"

func (tree *BTree) PrintAll(node BNode) {
	if node.btype() == BNODE_LEAF {
		nkeys := node.nkeys()
		for i := uint16(0); i < nkeys; i++ {
			fmt.Println("Key: ", string(node.getKey(i)), "Value: ", string(node.getVal(i)))
		}
		return
	}

	nkeys := node.nkeys()
	fmt.Println("Node", node.nbytes())
	for i := uint16(0); i < nkeys; i++ {
		ptr := node.getPtr(i)
		knode := tree.Get(ptr)
		tree.PrintAll(knode)
	}
}

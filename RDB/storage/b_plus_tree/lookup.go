package b_plus_tree

import (
	"bytes"
)

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
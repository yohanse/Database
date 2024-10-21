package linkedlist

import (
	"encoding/binary"
)

// node format:
// | next | pointers | unused |
// |  8B  |   n*8B   |   ...  |

type LNode []byte


func (node LNode) getNext() uint64 {
	return binary.LittleEndian.Uint64(node[:8])
}

func (node LNode) setNext(next uint64) {
	binary.LittleEndian.PutUint64(node[:8], next)
}

func (node LNode) getPtr(idx int) uint64 {
	return binary.LittleEndian.Uint64(node[8+8*idx:])
}
func (node LNode) setPtr(idx int, ptr uint64) {
	binary.LittleEndian.PutUint64(node[8+8*idx:], ptr)
}
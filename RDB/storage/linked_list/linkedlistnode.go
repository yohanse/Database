package linkedlist

import (
	"encoding/binary"
)

// node format:
// | next | pointers | unused |
// |  8B  |   n*8B   |   ...  |

type LNode []byte

func (node LNode) getNext() uint64 {
    return binary.LittleEndian.Uint64(node[:FREE_LIST_HEADER])
}

func (node LNode) setNext(next uint64) {
    binary.LittleEndian.PutUint64(node[:FREE_LIST_HEADER], next)
}

func (node LNode) getPtr(idx int) uint64 {
    offset := FREE_LIST_HEADER + idx*8
    return binary.LittleEndian.Uint64(node[offset:])
}

func (node LNode) setPtr(idx int, ptr uint64) {
    offset := FREE_LIST_HEADER + idx*8
    binary.LittleEndian.PutUint64(node[offset:], ptr)
}
package linkedlist

import (
	"rdb/storage/b_plus_tree"
)

func (fl *FreeList) PushTail(ptr uint64) {
	 // add it to the tail node
	LNode(fl.Set(fl.TailPage)).setPtr(seq2idx(fl.TailSeq), ptr)
    fl.TailSeq++
    // add a new tail node if it's full (the list is never empty)
    if seq2idx(fl.TailSeq) == 0 {
        // try to reuse from the list head
        next, head := flPop(fl) // may remove the head node
        if next == 0 {
            // or allocate a new node by appending
            next = fl.New(make([]byte, b_plus_tree.BTREE_PAGE_SIZE))
        }
        // link to the new tail node
        LNode(fl.Set(fl.TailPage)).setNext(next)
        fl.TailPage = next
        // also add the head node if it's removed
        if head != 0 {
            LNode(fl.Set(fl.TailPage)).setPtr(0, head)
            fl.TailSeq++
        }
    }
}
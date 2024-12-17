package linkedlist

func (fl *FreeList) PushTail(ptr uint64) {
	 // add it to the tail node
	LNode(fl.set(fl.tailPage)).setPtr(seq2idx(fl.tailSeq), ptr)
    fl.tailSeq++
    // add a new tail node if it's full (the list is never empty)
    if seq2idx(fl.tailSeq) == 0 {
        // try to reuse from the list head
        next, head := flPop(fl) // may remove the head node
        if next == 0 {
            // or allocate a new node by appending
            next = fl.new(make([]byte, BTREE_PAGE_SIZE))
        }
        // link to the new tail node
        LNode(fl.set(fl.tailPage)).setNext(next)
        fl.tailPage = next
        // also add the head node if it's removed
        if head != 0 {
            LNode(fl.set(fl.tailPage)).setPtr(0, head)
            fl.tailSeq++
        }
    }
}
}
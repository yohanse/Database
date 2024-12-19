package linkedlist

// remove 1 item from the head node, and remove the head node if empty.
func flPop(fl *FreeList) (ptr uint64, head uint64) {
    if fl.HeadSeq == fl.maxSeq {
        return 0, 0
    }
    node := LNode(fl.Get(fl.HeadPage))
    ptr = node.getPtr(seq2idx(fl.HeadSeq)) // item
    fl.HeadSeq++
    // move to the next one if the head node is empty
    if seq2idx(fl.HeadSeq) == 0 {
        head, fl.HeadPage = fl.HeadPage, node.getNext()
    }
    return
}
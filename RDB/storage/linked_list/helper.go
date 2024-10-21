package linkedlist

func seq2idx(seq uint64) int {
	return int(seq % FREE_LIST_CAP)
}

// make the newly added items available for consumption
func (fl *FreeList) SetMaxSeq() {
    fl.maxSeq = fl.tailSeq
}
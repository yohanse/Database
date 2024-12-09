package linkedlist

func (fl *FreeList) PopHead() uint64 {
	ptr, head := flPop(fl)

	if head != 0 {
		fl.put(head, fl.headPage)
	}
}
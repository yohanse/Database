package linkedlist

type FreeList struct {
    // callbacks for managing on-disk pages
    Get func(uint64) []byte // read a page
    New func([]byte) uint64 // append a new page
    Set func(uint64) []byte // update an existing page
    // persisted data in the meta page
    HeadPage uint64 // pointer to the list head node
    HeadSeq  uint64 // monotonic sequence number to index into the list head
    TailPage uint64
    TailSeq  uint64
    // in-memory states
    maxSeq uint64 // saved `tailSeq` to prevent consuming newly added items
}
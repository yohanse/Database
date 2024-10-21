package linkedlist

import (
	"rdb/storage/b_plus_tree"
)

const FREE_LIST_HEADER = 8
const FREE_LIST_CAP = (b_plus_tree.BTREE_PAGE_SIZE - FREE_LIST_HEADER) / 8
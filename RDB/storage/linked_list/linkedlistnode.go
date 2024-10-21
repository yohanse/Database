package linkedlist

// node format:
// | next | pointers | unused |
// |  8B  |   n*8B   |   ...  |

type LNode []byte
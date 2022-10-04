package cache

// N is a pure data struct representing a BVH tree node. This data struct is
// modified externally.
type N struct {
	// isValid is a private variable which indicates whether or not the
	// current node is used in the tree or not.
	isValid bool

	// id is not mutable by any package other than cache.
	id ID

	Parent ID
	Left   ID
	Right  ID
}

func (n *N) allocateOrLoad(x ID, parent ID, left ID, right ID) *N {
	if n == nil {
		n = &N{}
	}
	n.isValid = true

	n.id = x
	n.Parent = parent
	n.Left = left
	n.Right = right
	return n
}

func (n *N) free() {
	n.isValid = false
}

func (n *N) isAllocated() bool { return n.isValid }

func (n *N) ID() ID { return n.id }

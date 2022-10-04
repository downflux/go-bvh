package cache

// N is a pure data struct representing a BVH tree node. This data struct is
// modified externally.
type N struct {
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
	n.id = x
	n.Parent = parent
	n.Left = left
	n.Right = right
	return n
}

func (n *N) free() {}

func (n *N) ID() ID { return n.id }

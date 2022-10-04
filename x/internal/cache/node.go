package cache

// N is a pure data struct representing a BVH tree node. This data struct is
// modified externally.
type N struct {
	// isAllocated is a private variable which indicates whether or not the
	// current node is used in the tree or not.
	isAllocated bool

	// id is not mutable by any package other than cache.
	id ID

	parent ID
	left   ID
	right  ID
}

func (n *N) allocateOrLoad(x ID, parent ID, left ID, right ID) *N {
	if n == nil {
		n = &N{}
	}
	n.isAllocated = true

	n.id = x
	n.parent = parent
	n.left = left
	n.right = right
	return n
}

func (n *N) free() {
	n.isAllocated = false
}

func (n *N) IsAllocated() bool { return n.isAllocated }
func (n *N) ID() ID            { return n.id }

func (n *N) Parent() ID { return n.parent }
func (n *N) Left() ID   { return n.left }
func (n *N) Right() ID  { return n.right }

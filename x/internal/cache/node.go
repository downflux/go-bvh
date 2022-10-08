package cache

import (
	"fmt"
)

const (
	idSelf int = iota
	idParent
	idLeft
	idRight
)

// N is a pure data struct representing a BVH tree node. This data struct is
// modified externally.
type N struct {
	// isAllocated is a private variable which indicates whether or not the
	// current node is used in the tree or not.
	isAllocated bool

	ids [4]ID
}

func (n *N) allocateOrLoad(x ID, parent ID, left ID, right ID) *N {
	if n == nil {
		n = &N{}
	}
	n.isAllocated = true

	n.ids[idSelf] = x
	n.ids[idParent] = parent
	n.ids[idLeft] = left
	n.ids[idRight] = right
	return n
}

func (n *N) free() {
	n.isAllocated = false
}

func (n *N) IsAllocated() bool { return n.isAllocated }
func (n *N) ID() ID            { return n.ids[idSelf] }

func (n *N) Child(b B) ID {
	if !b.IsValid() {
		panic(fmt.Sprintf("invalid branch %v", b))
	}

	if b == BLeft {
		return n.ids[idLeft]
	}
	return n.ids[idRight]
}
func (n *N) SetChild(b B, x ID) {
	if !b.IsValid() {
		panic(fmt.Sprintf("invalid branch %v", b))
	}

	if b == BLeft {
		n.ids[idLeft] = x
	} else {
		n.ids[idRight] = x
	}
}

func (n *N) Parent() ID { return n.ids[idParent] }
func (n *N) Left() ID   { return n.Child(BLeft) }
func (n *N) Right() ID  { return n.Child(BRight) }

func (n *N) SetParent(x ID) { n.ids[idParent] = x }
func (n *N) SetLeft(x ID)   { n.SetChild(BLeft, x) }
func (n *N) SetRight(x ID)  { n.SetChild(BRight, x) }

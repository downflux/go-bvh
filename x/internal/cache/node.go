package cache

import (
	"fmt"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

const (
	idSelf int = iota
	idParent
	idLeft
	idRight
)

func DebugEqual(n, m *N) bool {
	if n.cache != m.cache {
		return false
	}
	if n.isAllocated != m.isAllocated {
		return false
	}
	if n.ids != m.ids {
		return false
	}
	if n.aabbCacheIsValid != m.aabbCacheIsValid {
		return false
	}
	if n.aabbCacheIsValid && !hyperrectangle.Within(n.aabbCache.R(), m.aabbCache.R()) {
		return false
	}

	return true
}

// N is a pure data struct representing a BVH tree node. This data struct is
// modified externally.
type N struct {
	cache *C

	// isAllocated is a private variable which indicates whether or not the
	// current node is used in the tree or not.
	isAllocated bool

	// aabbCache is a buffer to save bounding box calculations for faster
	// lookups. This cache is invalidated when the node children are
	// updated or when new data is attached to a leaf node. Because the node
	// itself does not track the cache nor the
	aabbCache        hyperrectangle.M
	aabbCacheIsValid bool

	ids [4]ID
}

func (n *N) allocateOrLoad(c *C, x ID, parent ID, left ID, right ID) *N {
	if n == nil {
		n = &N{}
	}
	n.isAllocated = true

	n.cache = c

	n.ids[idSelf] = x
	n.ids[idParent] = parent
	n.ids[idLeft] = left
	n.ids[idRight] = right
	return n
}

func (n *N) free() {
	n.isAllocated = false
	n.aabbCacheIsValid = false
}

func (n *N) IsAllocated() bool { return n.isAllocated }
func (n *N) ID() ID            { return n.ids[idSelf] }

func (n *N) IsRoot() bool {
	_, ok := n.cache.Get(n.ids[idParent])
	return !ok
}

func (n *N) Child(b B) *N {
	if !b.IsValid() {
		panic(fmt.Sprintf("invalid branch %v", b))
	}

	if b == BLeft {
		return n.cache.GetOrDie(n.ids[idLeft])
	}
	return n.cache.GetOrDie(n.ids[idRight])
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

// Branch returns the branch of the input child in relation to the current node.
func (n *N) Branch(x ID) B {
	switch x {
	case n.ids[idLeft]:
		return BLeft
	case n.ids[idRight]:
		return BRight
	default:
		return BInvalid
	}
}

func (n *N) Parent() *N { return n.cache.GetOrDie(n.ids[idParent]) }
func (n *N) Left() *N   { return n.Child(BLeft) }
func (n *N) Right() *N  { return n.Child(BRight) }

func (n *N) SetParent(x ID) { n.ids[idParent] = x }
func (n *N) SetLeft(x ID)   { n.SetChild(BLeft, x) }
func (n *N) SetRight(x ID)  { n.SetChild(BRight, x) }

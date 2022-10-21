package cache

import (
	"fmt"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

const (
	idSelf int = iota
	idParent
	idLeft
	idRight
)

func (n *N) Within(m *N) bool {
	if n.cache != m.cache {
		return false
	}
	if n.isAllocated != m.isAllocated {
		return false
	}
	if n.ids != m.ids {
		return false
	}
	if !hyperrectangle.Within(n.aabbCache.R(), m.aabbCache.R()) {
		return false
	}
	if len(n.dataCache) != len(m.dataCache) {
		return false
	}
	for k, v := range n.dataCache {
		if m.dataCache[k] != v {
			return false
		}
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
	// updated or when new data is attached to a leaf node.
	//
	// N.B.: As with the node relationship IDs, the aabbCache must be set
	// manually by the caller.
	aabbCache hyperrectangle.M

	// dataCache is a buffer for leaf nodes to track AABB child objects. The
	// caller is responsible for tracking the actual AABBs, since the
	// objects may move between nodes during operations.
	//
	// N.B.: This cache is manually set by the caller.
	dataCache map[id.ID]bool

	ids [4]ID
}

func (n *N) allocateOrLoad(c *C, x ID, parent ID, left ID, right ID) *N {
	if n == nil {
		n = &N{
			cache: c,
			aabbCache: hyperrectangle.New(
				vector.V(make([]float64, c.K())),
				vector.V(make([]float64, c.K())),
			).M(),
			dataCache: make(map[id.ID]bool, c.LeafSize()),
		}
		n.ids[idSelf] = x
	}

	if c != n.cache {
		panic("cannot set cache again after allocation")
	}
	if x != n.ids[idSelf] {
		panic(fmt.Sprintf("cannot set node ID again after allocation %v", x))
	}

	n.isAllocated = true

	n.ids[idParent] = parent
	n.ids[idLeft] = left
	n.ids[idRight] = right

	// N.B.: The AABB cache is not guaranteed to be zeroed at the end of the
	// allocation. The caller will usually copy data into this cache
	// directly, so zeroing at this point would be useless.

	return n
}

func (n *N) free() {
	n.isAllocated = false

	// Since the dataCache represents external data (which may be freed
	// outside the cache), we should remove all references to that data when
	// the node is marked invalid.
	for k := range n.dataCache {
		delete(n.dataCache, k)
	}
}

func (n *N) IsAllocated() bool { return n.isAllocated }
func (n *N) ID() ID            { return n.ids[idSelf] }

// AABB returns the bounding box of the node. This bounding box may be mutated
// by the caller.
func (n *N) AABB() hyperrectangle.M { return n.aabbCache }

// Data returns the list of AABBs contained in this node. The node must be a
// leaf node.
//
// This cache may be mutated by the caller.
func (n *N) Data() map[id.ID]bool { return n.dataCache }

func (n *N) IsRoot() bool {
	_, ok := n.cache.Get(n.ids[idParent])
	return !ok
}

// IsLeaf returns if the current node has no valid children.
//
// N.B.: A valid BVH tree must have either both children be valid, or no valid
// children (and contain only data). We are not checking the right child here.
func (n *N) IsLeaf() bool {
	_, ok := n.cache.Get(n.ids[idLeft])
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

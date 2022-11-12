// Package impl is a node implementation with manual garbage collection support.
// This package should only be called by the cache package or for testing
// purposes.
package impl

import (
	"fmt"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache/branch"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

const (
	idSelf int = iota
	idParent
	idLeft
	idRight
)

// A is an node interface. This is our cache back-reference. We are
// defining this as an interface to avoid cyclic imports and to avoid a bloated
// module.
type A interface {
	Get(x cid.ID) (node.N, bool)
	LeafSize() int
	K() vector.D
}

// N is a pure data struct representing a BVH tree node. This data struct is
// modified externally.
type N struct {
	cache A

	// isAllocated is a private variable which indicates whether or not the
	// current node is used in the tree or not.
	isAllocated bool

	// aabbCache is a buffer to save bounding box calculations for faster
	// lookups. This cache is invalidated when the node children are
	// updated or when new data is attached to a leaf node.
	//
	// N.B.: As with the node relationship IDs, the aabbCache must be set
	// manually by the caller.
	//
	// aabbCache may be expanded by some additional factor as compared to
	// the child nodes; this useful for frequently updated trees to reduce
	// object insert and remove churn.
	aabbCache hyperrectangle.M

	// dataCache is a buffer for leaf nodes to track AABB child objects. The
	// caller is responsible for tracking the actual AABBs, since the
	// objects may move between nodes during operations.
	//
	// N.B.: This cache is manually set by the caller.
	//
	// N.B.: A valid node AABB must contain all data objects, and may also
	// be extended in each direction as a buffer. This buffer is useful for
	// minimizing the amount of frivolous tree add / remove operations, per
	// Catto 2019.
	dataCache map[id.ID]struct{}

	heightCache int

	ids [4]cid.ID
}

func New(a A, x cid.ID) *N {
	return &N{
		cache: a,
		aabbCache: hyperrectangle.New(
			vector.V(make([]float64, a.K())),
			vector.V(make([]float64, a.K())),
		).M(),
		dataCache: make(map[id.ID]struct{}, a.LeafSize()),
		ids: [4]cid.ID{
			/* idSelf = */ x,
			/* idParent = */ cid.IDInvalid,
			/* idLeft = */ cid.IDInvalid,
			/* idRight = */ cid.IDInvalid,
		},
	}
}

// Allocate resets an unallocated node with new neighbors.
func (n *N) Allocate(parent cid.ID, left cid.ID, right cid.ID) {
	if n.IsAllocated() {
		panic("cannot re-allocate an existing node")
	}

	n.isAllocated = true

	n.ids[idLeft] = left
	n.ids[idRight] = right
	n.ids[idParent] = parent

	// N.B.: The AABB cache is not guaranteed to be zeroed at the end of the
	// allocation. The caller will usually copy data into this cache
	// directly, so zeroing at this point would be useless.
}

func (n *N) Free() {
	if !n.IsAllocated() {
		return
	}

	n.isAllocated = false

	// Since the dataCache represents external data (which may be freed
	// outside the cache), we should remove all references to that data when
	// the node is marked invalid.
	for k := range n.dataCache {
		delete(n.dataCache, k)
	}
}

func (n *N) IsAllocated() bool { return n != nil && n.isAllocated }

func (n *N) ID() cid.ID {
	if !n.IsAllocated() {
		panic("accessing an unallocated node")
	}

	return n.ids[idSelf]
}

func (n *N) IsRoot() bool {
	if !n.IsAllocated() {
		panic("accessing an unallocated node")
	}

	_, ok := n.cache.Get(n.ids[idParent])
	return !ok
}

func (n *N) Height() int {
	if !n.IsAllocated() {
		panic("accessing an unallocated node")
	}

	return n.heightCache
}

func (n *N) SetHeight(h int) {
	if !n.IsAllocated() {
		panic("accessing an unallocated node")
	}

	n.heightCache = h
}

// IsLeaf returns if the current node has no valid children.
//
// N.B.: A valid BVH tree must have either both children be valid, or no valid
// children (and contain only data). We are not checking the right child here.
func (n *N) IsLeaf() bool {
	if !n.IsAllocated() {
		panic("accessing an unallocated node")
	}

	_, ok := n.cache.Get(n.ids[idLeft])
	return !ok
}

func (n *N) IsFull() bool {
	if !n.IsAllocated() {
		panic("accessing an unallocated node")
	}

	if !n.IsLeaf() {
		panic(fmt.Sprintf("internal node %v does not have a data cache size", n.ID()))
	}
	return len(n.dataCache) >= n.cache.LeafSize()
}

// Leaves returns the list of AABBs contained in this node. The node must be a
// leaf node.
//
// This cache may be mutated by the caller.
//
// N.B.: We make the assumption that every member of the map here is a valid
// object -- in order to remove an object from the node, the caller must instead
// make a delete call, e.g.
//
//	delete(n.Data(), x)
//
// The member cannot be invalidated by setting the corresponding key to some
// null struct (as this is already the default value).
//
// To test for membership, use
//
//	if _, ok := n.Data()[x]; ok { ... }
func (n *N) Leaves() map[id.ID]struct{} {
	if !n.IsAllocated() {
		panic("accessing an unallocated node")
	}

	return n.dataCache
}

// AABB returns the bounding box of the node. This bounding box may be mutated
// by the caller.
func (n *N) AABB() hyperrectangle.M {
	if !n.IsAllocated() {
		panic("accessing an unallocated node")
	}

	return n.aabbCache
}

func (n *N) Child(b branch.B) node.N {
	if !n.IsAllocated() {
		panic("accessing an unallocated node")
	}

	if !b.IsValid() {
		panic(fmt.Sprintf("invalid branch %v", b))
	}

	var id int
	switch b {
	case branch.BLeft:
		id = idLeft
	case branch.BRight:
		id = idRight
	}
	m, ok := n.cache.Get(n.ids[id])
	if !ok {
		return nil
	}
	return m
}

func (n *N) SetChild(b branch.B, x cid.ID) {
	if !n.IsAllocated() {
		panic("accessing an unallocated node")
	}

	if !b.IsValid() {
		panic(fmt.Sprintf("invalid branch %v", b))
	}

	switch b {
	case branch.BLeft:
		n.ids[idLeft] = x
	case branch.BRight:
		n.ids[idRight] = x
	}
}

// Branch returns the branch of the input child in relation to the current node.
func (n *N) Branch(x cid.ID) branch.B {
	if !n.IsAllocated() {
		panic("accessing an unallocated node")
	}

	switch x {
	case n.ids[idLeft]:
		return branch.BLeft
	case n.ids[idRight]:
		return branch.BRight
	default:
		return branch.BInvalid
	}
}

func (n *N) Parent() node.N {
	if !n.IsAllocated() {
		panic("accessing an unallocated node")
	}

	m, ok := n.cache.Get(n.ids[idParent])
	if !ok {
		return nil
	}
	return m
}

func (n *N) Left() node.N  { return n.Child(branch.BLeft) }
func (n *N) Right() node.N { return n.Child(branch.BRight) }

func (n *N) SetParent(x cid.ID) {
	if !n.IsAllocated() {
		panic("accessing an unallocated node")
	}

	n.ids[idParent] = x
}

func (n *N) SetLeft(x cid.ID)  { n.SetChild(branch.BLeft, x) }
func (n *N) SetRight(x cid.ID) { n.SetChild(branch.BRight, x) }

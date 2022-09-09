// Package node is an internal-only node implementation struct, and its
// properties and data points should only be accessed via the operations API in
// the /internal/node/op/ directory.
package node

import (
	"fmt"
	"math"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node/stack"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
	nid "github.com/downflux/go-bvh/internal/node/id"
)

type C struct {
	lookup map[nid.ID]*N
}

func Cache() *C {
	return &C{
		lookup: map[nid.ID]*N{},
	}
}

func (c *C) Allocate() nid.ID {
	var id nid.ID
	for ; id.IsZero(); id = nid.Generate() {
	}
	for _, ok := c.lookup[id]; ok; {
		id = nid.Increment(id)
	}
	c.lookup[id] = nil
	return id
}

func (c *C) Delete(id nid.ID) {
	delete(c.lookup, id)
}

type O struct {
	// Nodes ia a node lookup table. This is to reduce the amount of random
	// memory jumps necessary with a pointer implementation, as well as
	// frustrating double-pointer details, e.g. memory leaks and test
	// verbosity.
	Nodes *C

	// ID represents the internal node ID. If the ID is unspecified, a
	// random ID will be generated.
	ID nid.ID

	// N.B.: The left, right, and parent node IDs may not exist or be
	// allocated in the cache at construct time, as the nodes may be built
	// out of order. If the given IDs still do not exist when calling e.g.
	// n.Left(), a nil value will be returned.

	Left   nid.ID
	Right  nid.ID
	Parent nid.ID

	// Data represents an object to AABB lookup for the leaf node. Note that
	// these AABBs here do not necessarily represent the physical boundaries
	// of the object -- remember that these AABBs could take into account
	// some physical buffer space to reduce the amount of tree mutations
	// that occur, per Catto 2019.
	Data map[id.ID]hyperrectangle.R

	// Size represents how many objects may be added to a leaf node before
	// the node splits.
	Size uint
}

type N struct {
	nodes *C

	id     nid.ID
	left   nid.ID
	right  nid.ID
	parent nid.ID

	data map[id.ID]hyperrectangle.R
	size uint

	aabbCacheIsValid bool
	aabbCache        hyperrectangle.R

	heightCacheIsValid bool
	heightCache        uint
}

func New(o O) *N {
	if uint(len(o.Data)) > o.Size {
		panic(fmt.Sprintf("cannot create node with data exceeding size %v", o.Size))
	}

	// If the user does not specify an ID, automatically allocate one to
	// use.
	id := o.ID
	if id.IsZero() {
		id = o.Nodes.Allocate()
	}

	n := &N{
		nodes:  o.Nodes,
		id:     id,
		left:   o.Left,
		right:  o.Right,
		parent: o.Parent,
		data:   o.Data,
		size:   o.Size,
	}

	if o.Nodes.lookup[n.ID()] != nil {
		panic("cannot add node with duplicate ID")
	}
	o.Nodes.lookup[n.ID()] = n
	return n
}

// Get returns the AABB of the associated object from a leaf node. If the object
// is not found, a false value will be returned.
func (n *N) Get(id id.ID) (hyperrectangle.R, bool) {
	if !n.IsLeaf() {
		panic("cannot get a bounding box from a non-leaf node")
	}

	aabb, ok := n.data[id]
	return aabb, ok
}

func (n *N) Insert(id id.ID, aabb hyperrectangle.R) {
	if _, ok := n.data[id]; ok {
		panic(fmt.Sprintf("cannot insert an existing object ID %v", id))
	}
	if n.IsFull() {
		panic(fmt.Sprintf("cannot insert into a full leaf node"))
	}

	n.data[id] = aabb
	n.invalidateAABBCache()
}

func (n *N) Remove(id id.ID) {
	if _, ok := n.data[id]; !ok {
		panic(fmt.Sprintf("cannot find specified object ID %v", id))
	}

	delete(n.data, id)
	n.invalidateAABBCache()
}

// Return the list of entities in this node.
//
// N.B.: This data is the original copy stored in the node -- the caller must
// not mutate this. We are returning the original data because the caller is the
// library itself, and to save on an extra memory allocation.
func (n *N) Data() map[id.ID]hyperrectangle.R { return n.data }

func (n *N) Cache() *C  { return n.nodes }
func (n *N) Size() uint { return n.size }

func (n *N) invalidateHeightCache() {
	if !n.heightCacheIsValid {
		return
	}
	n.heightCacheIsValid = false
	if !n.IsRoot() {
		n.Parent().invalidateHeightCache()
	}
}

func (n *N) invalidateAABBCache() {
	// Since invalidateAABBCache is called recursively up towards the root,
	// and AABB is calculated towards the leaf, if the cache is invalid at
	// some node, we are guaranteed all nodes above the current node are
	// also marked with an invalid cache. Skipping the tree iteration here
	// can reduce the complexity by a factor of O(log N) if we are
	// traveling up the tree anyway in some other algorithm.
	if !n.aabbCacheIsValid {
		return
	}

	n.aabbCacheIsValid = false
	if !n.IsRoot() {
		n.Parent().invalidateAABBCache()
	}
}

// Return the node ID.
func (n *N) ID() nid.ID { return n.id }

func (n *N) Left() *N   { return n.nodes.lookup[n.left] }
func (n *N) Right() *N  { return n.nodes.lookup[n.right] }
func (n *N) Parent() *N { return n.nodes.lookup[n.parent] }

func (n *N) SetLeft(m *N) {
	n.left = m.ID()
	n.invalidateHeightCache()
	n.invalidateAABBCache()
}
func (n *N) SetRight(m *N) {
	n.right = m.ID()
	n.invalidateHeightCache()
	n.invalidateAABBCache()
}
func (n *N) SetParent(m *N) {
	// We may be attempting to set the node as the new root.
	var id nid.ID
	if m != nil {
		id = m.ID()
		n.invalidateHeightCache()
		n.invalidateAABBCache()
	}
	n.parent = id
}

func (n *N) Root() *N {
	var p *N
	for p = n; !p.IsRoot(); p = p.Parent() {
	}
	return p
}

// BroadPhase returns a list of object IDs which intersect with the query
// hyperrectangle.
//
// Further refinement should be done by the caller to check if the objects
// actually collide.
func (n *N) BroadPhase(q hyperrectangle.R) []id.ID {
	open := stack.New(make([]*N, 0, 128))
	open.Push(n)

	ids := stack.New(make([]id.ID, 0, 128))

	for m, ok := open.Pop(); ok; m, ok = open.Pop() {
		if m.IsLeaf() {
			for id, h := range m.data {
				if !bhr.Disjoint(q, h) {
					ids.Push(id)
				}
			}
		} else {
			open.Push(m.Left())
			open.Push(m.Right())
		}
	}

	return ids.Data()
}

func (n *N) IsLeaf() bool  { return n.left.IsZero() && n.right.IsZero() }
func (n *N) IsRoot() bool  { return n.Parent() == nil }
func (n *N) IsFull() bool  { return n.IsLeaf() && uint(len(n.data)) >= n.Size() }
func (n *N) IsEmpty() bool { return n.IsLeaf() && len(n.data) == 0 }
func (n *N) Height() uint {
	if n.heightCacheIsValid {
		return n.heightCache
	}

	n.heightCacheIsValid = true
	if n.IsLeaf() {
		n.heightCache = 1
	} else {
		n.heightCache = 1 + uint(math.Max(float64(n.Left().Height()), float64(n.Right().Height())))
	}
	return n.heightCache
}
func (n *N) AABB() hyperrectangle.R {
	if n.aabbCacheIsValid {
		return n.aabbCache
	}

	n.aabbCacheIsValid = true
	if n.IsLeaf() {
		if n.IsEmpty() {
			panic("AABB is not defined for an empty leaf node")
		}
		rs := make([]hyperrectangle.R, 0, len(n.data))
		for _, aabb := range n.data {
			rs = append(rs, aabb)
		}
		n.aabbCache = bhr.AABB(rs)
	} else {
		n.aabbCache = bhr.Union(n.Left().AABB(), n.Right().AABB())
	}

	return n.aabbCache
}

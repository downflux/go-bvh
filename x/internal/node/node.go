// Package node is an internal-only node implementation struct, and its
// properties and data points should only be accessed via the operations API in
// the /internal/node/op/ directory.
package node

import (
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/x/hyperrectangle"
	nid "github.com/downflux/go-bvh/x/internal/node/id"
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
	for ; id.IsZero(); id = nid.Generate() {}
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
}

type N struct {
	nodes *C

	id     nid.ID
	left   nid.ID
	right  nid.ID
	parent nid.ID

	data map[id.ID]hyperrectangle.R

	aabbCacheIsValid bool
	aabbCache        hyperrectangle.R
}

func New(o O) *N {
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

func (n *N) Cache() *C { return n.nodes }

func (n *N) InvalidateAABBCache() {
	// Since InvalidateAABBCache is called recursively up towards the root,
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
		n.Parent().InvalidateAABBCache()
	}
}

func (n *N) ID() nid.ID { return n.id }

func (n *N) Left() *N   { return n.nodes.lookup[n.left] }
func (n *N) Right() *N  { return n.nodes.lookup[n.right] }
func (n *N) Parent() *N { return n.nodes.lookup[n.parent] }

func (n *N) SetLeft(m *N)   { n.left = m.ID() }
func (n *N) SetRight(m *N)  { n.right = m.ID() }
func (n *N) SetParent(m *N) {
	// We may be attempting to set the node as the new root.
	var id nid.ID
	if m != nil {
		id = m.ID()
	}
	n.parent = id
}

func (n *N) Root() *N {
	if n.Parent() == nil {
		return n
	}
	return n.Parent().Root()
}

// BroadPhase returns a list of object IDs which intersect with the query
// hyperrectangle.
//
// Further refinement should be done by the caller to check if the objects
// actually collide.
func (n *N) BroadPhase(q hyperrectangle.R) []id.ID {
	if bhr.Disjoint(q, n.AABB()) {
		return nil
	}

	if n.IsLeaf() {
		ids := make([]id.ID, 0, len(n.data))
		for id, h := range n.data {
			if !bhr.Disjoint(q, h) {
				ids = append(ids, id)
			}
		}
		return ids
	}
	l := make(chan []id.ID)
	r := make(chan []id.ID)
	go func(ch chan<- []id.ID) {
		defer close(ch)
		ch <- n.Left().BroadPhase(q)
	}(l)
	go func(ch chan<- []id.ID) {
		defer close(ch)
		ch <- n.Right().BroadPhase(q)
	}(r)

	return append(<-l, <-r...)
}

func (n *N) IsLeaf() bool { return len(n.data) > 0 }
func (n *N) IsRoot() bool { return n.Parent() == nil }
func (n *N) AABB() hyperrectangle.R {
	if n.aabbCacheIsValid {
		return n.aabbCache
	}

	n.aabbCacheIsValid = true
	if n.IsLeaf() {
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

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
	id := nid.Generate()
	for _, ok := c.lookup[id]; ok; {
		id = nid.Increment(id)
	}
	c.lookup[id] = nil
	return id
}

type O struct {
	Nodes *C

	ID     nid.ID
	Left   nid.ID
	Right  nid.ID
	Parent nid.ID

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
	if id.IsNil() {
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
func (n *N) SetParent(m *N) { n.parent = m.ID() }

func (n *N) Root() *N {
	if n.Parent() == nil {
		return n
	}
	return n.Parent().Root()
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

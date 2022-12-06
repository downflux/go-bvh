package node

import (
	"github.com/downflux/go-bvh/y/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

// N is a tree node struct. Note that this tree does not have a way to get its
// parent -- the path to the node must be constructed via a search operation.
type N struct {
	id ID

	k         vector.D
	n         int
	tolerance float64

	aabbCache     hyperrectangle.M
	aabbCostCache float64

	heightCache int

	children map[ID]*N
	isLeaf   bool
	leaf     id.ID
}

func (n *N) ID() ID       { return n.id }
func (n *N) IsLeaf() bool { return n.isLeaf }
func (n *N) Height() int  { return n.heightCache }

// AABB returns the bounding box of the node. External callers must not mutate
// this struct into its mutable version; the caller will need to call SetAABB()
// explicitly.
func (n *N) AABB() hyperrectangle.R { return n.aabbCache.R() }

func (n *N) Leaf() id.ID {
	if !n.isLeaf {
		panic("cannot get leaf of a non-leaf node")
	}

	return n.leaf
}

func (n *N) SetLeaf(x id.ID) {
	if len(n.children) != 0 {
		panic("cannot set leaf of a non-leaf node")
	}

	n.isLeaf = true
	n.leaf = x
}

// Children returns the child nodes of a given node instance. The nodes here may
// be mutated.
func (n *N) Children() map[ID]*N { return n.children }

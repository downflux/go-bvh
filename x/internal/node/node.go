package node

import (
	"fmt"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
)

type Branch int

func (b Branch) Sibling() Branch { return b ^ 1 }
func (b Branch) IsValid() bool {
	return b == BranchLeft || b == BranchRight
}

const (
	BranchLeft Branch = iota
	BranchRight

	BranchInvalid
)

type N struct {
	id       cache.ID
	parent   cache.ID
	children [2]cache.ID

	roNodes  cache.RO[*N]
	roLeaves map[cache.ID]map[id.ID]hyperrectangle.R

	heightCache        int
	heightCacheIsValid bool

	aabbCache        hyperrectangle.R
	aabbCacheIsValid bool
}

type O struct {
	Nodes  *cache.C[*N]
	Leaves map[cache.ID]map[id.ID]hyperrectangle.R

	Parent cache.ID
	Left   cache.ID
	Right  cache.ID
}

func New(o O) *N {
	n := &N{
		roNodes:  o.Nodes,
		roLeaves: o.Leaves,

		parent:   o.Parent,
		children: [2]cache.ID{o.Left, o.Right},
	}
	x := o.Nodes.Insert(n)

	n.id = x
	return n
}

// AABB returns the bounding box of the subtree. We assume the input node is
// valid.
func (n *N) AABB() hyperrectangle.R {
	x := n.ID()

	if n.aabbCacheIsValid {
		return n.aabbCache
	}

	n.aabbCacheIsValid = true

	if n.IsLeaf() {
		if len(n.roLeaves[x]) == 0 {
			panic(fmt.Sprintf("AABB is not defined for an empty leaf node %v", x))
		}

		rs := make([]hyperrectangle.R, 0, len(n.roLeaves[x]))
		for _, aabb := range n.roLeaves[x] {
			rs = append(rs, aabb)
		}

		bhr.AABBBuf(rs, n.aabbCache)
	} else {
		bhr.UnionBuf(
			n.Left().AABB(),
			n.Right().AABB(),
			n.aabbCache,
		)
	}

	return n.aabbCache
}

// Height returns the subtree height. We assume the input node is valid.
//
// Leaf nodes have a height of 0.
func (n *N) Height() int {
	if n.heightCacheIsValid {
		return n.heightCache
	}

	n.heightCacheIsValid = true

	if n.IsLeaf() {
		n.heightCache = 0
	} else {
		h := n.Left().Height()
		if g := n.Right().Height(); g > h {
			h = g
		}
		n.heightCache = 1 + h
	}

	return n.heightCache
}

func (n *N) ID() cache.ID { return n.id }

func (n *N) IsRoot() bool { return !n.parent.IsValid() }
func (n *N) IsLeaf() bool {
	return !n.children[BranchLeft].IsValid() && !n.children[BranchRight].IsValid()
}

func (n *N) Parent() *N {
	m, ok := n.roNodes.Get(n.parent)
	if !ok {
		return nil
	}

	return m
}

func (n *N) Branch(child cache.ID) Branch {
	if n.children[BranchLeft] == child {
		return BranchLeft
	}
	if n.children[BranchRight] == child {
		return BranchRight
	}
	return BranchInvalid
}

// Child is a convenience function for programatic tree explorations -- instead
// of calling
//
//	n.Left()
//
// we can instead call
//
//	n.Child(BranchLeft)
func (n *N) Child(b Branch) *N {
	if !b.IsValid() {
		panic(fmt.Sprintf("invalid branch option %v", b))
	}

	m, ok := n.roNodes.Get(n.children[b])
	if !ok {
		return nil
	}

	return m
}

func (n *N) Left() *N  { return n.Child(BranchLeft) }
func (n *N) Right() *N { return n.Child(BranchRight) }

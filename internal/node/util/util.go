package util

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	nid "github.com/downflux/go-bvh/internal/node/id"
)

func Interval(min, max float64) hyperrectangle.R {
	return *hyperrectangle.New(*vector.New(min), *vector.New(max))
}

type T struct {
	Data  map[nid.ID]map[id.ID]hyperrectangle.R
	Nodes map[nid.ID]N
	Root  nid.ID
	Size  uint
}

type N struct {
	Left   nid.ID
	Right  nid.ID
	Parent nid.ID
}

func New(t T) *node.N {
	c := node.Cache()

	var r *node.N
	for id, n := range t.Nodes {
		m := node.New(node.O{
			Nodes: c,

			ID:     id,
			Left:   n.Left,
			Right:  n.Right,
			Parent: n.Parent,

			Data: t.Data[id],
			Size: t.Size,
		})
		if m.ID() == t.Root {
			r = m
		}
	}
	return r
}

func Equal(a *node.N, b *node.N) bool {
	if a == nil && b == nil {
		return true
	}
	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}
	if (a.IsLeaf() && !b.IsLeaf()) || (!a.IsLeaf() && b.IsLeaf()) {
		return false
	}
	if !cmp.Equal(a.AABB(), b.AABB(), cmp.AllowUnexported(hyperrectangle.R{})) {
		return false
	}
	if a.Height() != b.Height() {
		return false
	}

	if a.IsLeaf() && b.IsLeaf() {
		return cmp.Equal(a, b,
			cmpopts.IgnoreFields(
				node.N{},
				"nodes",
				"id",
				"parent",
				"left",
				"right",
			),
			cmp.AllowUnexported(node.N{}, node.C{}, hyperrectangle.R{}),
		)
	}
	// Check the nodes match orientation. This is purely for debugging
	// purposes -- in real life, left and right node swaps do not matter.
	return Equal(a.Left(), b.Left()) && Equal(a.Right(), b.Right())
}

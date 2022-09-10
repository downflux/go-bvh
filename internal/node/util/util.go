package util

import (
	"log"
	"math"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/heuristic"
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

// MaxImbalance returns the maximum imbalance within the tree.
func MaxImbalance(n *node.N) uint {
	if n.IsLeaf() {
		return 0
	}

	return uint(math.Max(
		math.Max(
			float64(MaxImbalance(n.Left())),
			float64(MaxImbalance(n.Right())),
		),
		math.Abs(float64(n.Left().Height())-float64(n.Right().Height())),
	))
}

// Cost returns the total tree heuristic.
func Cost(n *node.N) float64 {
	if n.IsLeaf() {
		return heuristic.H(n.AABB())
	}
	return heuristic.H(n.AABB()) + Cost(n.Left()) + Cost(n.Right())
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

func PreOrder(n *node.N, f func(n *node.N)) {
	f(n)
	if !n.IsLeaf() {
		PreOrder(n.Left(), f)
		PreOrder(n.Right(), f)
	}
}

func Log(l *log.Logger, n *node.N) {
	PreOrder(n, func(n *node.N) {
		l.Printf(
			"node NID: %v, L: %v, R: %v, P: %v, Data: %v\n",
			n.ID(),
			func() nid.ID {
				if !n.IsLeaf() {
					return n.Left().ID()
				}
				return 0
			}(),
			func() nid.ID {
				if !n.IsLeaf() {
					return n.Right().ID()
				}
				return 0
			}(),
			func() nid.ID {
				if !n.IsRoot() {
					return n.Parent().ID()
				}
				return 0
			}(),
			func() map[id.ID]hyperrectangle.R {
				if n.IsLeaf() {
					return n.Data()
				}
				return nil
			}(),
		)
	})
}

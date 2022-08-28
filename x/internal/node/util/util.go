package util

import (
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Interval(min, max float64) hyperrectangle.R {
	return *hyperrectangle.New(*vector.New(min), *vector.New(max))
}

type NodeID uint64

type T struct {
	Data  map[NodeID][]node.D
	Nodes map[NodeID]N
	Root  NodeID
}

type N struct {
	Left  NodeID
	Right NodeID
}

func New(t T) *node.N {
	if len(t.Data[t.Root]) > 0 {
		return node.New(node.O{
			Data: t.Data[t.Root],
		})
	}
	return node.New(node.O{
		Left: New(T{
			Data:  t.Data,
			Nodes: t.Nodes,
			Root:  t.Nodes[t.Root].Left,
		}),
		Right: New(T{
			Data:  t.Data,
			Nodes: t.Nodes,
			Root:  t.Nodes[t.Root].Right,
		}),
	})
}

func Equal(a *node.N, b *node.N) bool {
	if (a.IsLeaf() && !b.IsLeaf()) || (!a.IsLeaf() && b.IsLeaf()) {
		return false
	}
	if !cmp.Equal(a.AABB(), b.AABB(), cmp.AllowUnexported(hyperrectangle.R{})) {
		return false
	}

	if a.IsLeaf() && b.IsLeaf() {
		return cmp.Equal(
			a,
			b,
			cmpopts.IgnoreFields(
				node.N{},
				"left",
				"right",
				"parent",
			),
			cmp.AllowUnexported(node.N{}, hyperrectangle.R{}),
		)
	}
	return (Equal(a.Left(), b.Left()) && Equal(a.Right(), b.Right())) || (Equal(a.Left(), b.Right()) && Equal(a.Right(), b.Left()))
}

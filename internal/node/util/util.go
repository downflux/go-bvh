package util

import (
	"math"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/heuristic"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
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

func k(data map[nid.ID]map[id.ID]hyperrectangle.R) vector.D {
	for _, leaf := range data {
		for _, aabb := range leaf {
			return aabb.Min().Dimension()
		}
	}
	return 0
}

func New(t T) *node.N {
	if len(t.Data) == 0 {
		panic("cannot create a root node with no data")
	}

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
			K:    k(t.Data),
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

// SAH returns the surface area heuristic as defined in Macdonald and Booth
// 1990.
//
// The total heuristic is comprised of three separate components -- the cost of
// the internal nodes, the cost of the leaves, and the cost of testing for
// intersections. We use track these via ci, cl, and co respectively.
//
// Per Aila et al., a "normal" SAH value is around 100.
func SAH(n *node.N) float64 {
	var ci, cl, co float64
	PreOrder(n, func(n *node.N) {
		if !n.IsLeaf() {
			ci += heuristic.H(n.AABB())
		} else {
			cl += heuristic.H(n.AABB())
			co += heuristic.H(n.AABB()) * float64(len(n.Data()))
		}
	})
	return (1.2*ci + 1.0*cl + 0*co) / heuristic.H(n.Root().AABB())
}

// OverlapPenalty generates a heuristic for checking how many nodes in the tree
// have overlapping children, and as such, will require additional node
// expansions.
//
// A lower value is better.
func OverlapPenalty(n *node.N) float64 {
	var p float64
	PreOrder(n, func(n *node.N) {
		if !n.IsLeaf() && !bhr.Disjoint(n.Left().AABB(), n.Right().AABB()) {
			p += float64(n.Height()) * float64(n.Height())
		}
	})
	return math.Sqrt(p)
}

// BalancePenalty checks for how balanced the tree is.
//
// A lower value is better.
func BalancePenalty(n *node.N) float64 {
	depth := n.Height()
	layers := make([]int, n.Height()+1)
	PreOrder(n, func(n *node.N) {
		layers[depth-n.Height()] += 1
	})

	var p float64
	for d, count := range layers {
		n := math.Pow(2.0, float64(d))
		// A perfect binary tree will have 2 ** d nodes per layer. We
		// penalize a tree if its total node count at this layer strays
		// from the ideal.
		p += math.Abs((n - float64(count))) / n
	}
	return p
}

func AverageSize(n *node.N) float64 {
	var sizes []float64
	PreOrder(n, func(n *node.N) {
		if n.IsLeaf() {
			sizes = append(sizes, float64(len(n.Data())))
		}
	})

	var sum float64
	for _, s := range sizes {
		sum += s
	}

	return sum / float64(len(sizes))
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
	if !hyperrectangle.Within(a.AABB(), b.AABB()) {
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

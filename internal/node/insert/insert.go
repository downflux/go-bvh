package insert

import (
	"math"

	"github.com/downflux/go-bvh/internal/heuristic"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
)

func Insert(n *node.N, id point.ID, bound hyperrectangle.R) {
	// A new leaf node will necessitate creating a new parent p, which will
	// replace the current node's parent, as well as the cost of the new
	// node b itself.
	//
	//     p
	//    / \
	//   b   n
	create := 2 * heuristic.Heuristic(bhr.Union(bound, n.Bound))

	if n.Leaf() || create < h(n, bound) {
		m := &node.N{
			Parent: n.Parent,
			Left: &node.N{
				ID:    id,
				Bound: bound,
			},
			Right: n,
			Bound: bhr.Union(bound, n.Bound),
		}

		// Update the original parent to point to the newly created
		// internal node.
		if n.Parent != nil {
			if n.Parent.Left == n {
				n.Parent.Left = m
			} else {
				n.Parent.Right = m
			}
		}

		// Update the current nodes.
		n.Parent = m
		m.Left.Parent = m
	} else {
		if h(n.Left, bound) < h(n.Right, bound) {
			Insert(n.Left, id, bound)
		} else {
			Insert(n.Right, id, bound)
		}
	}
}
func h(n *node.N, bound hyperrectangle.R) heuristic.H {
	if n.Leaf() {
		// A new leaf node will necessitate creating a new parent p,
		// which will replace the current node's parent, as well as the
		// cost of the new node b itself.
		//
		//     p
		//    / \
		//   b   n
		return 2 * heuristic.Heuristic(bhr.Union(bound, n.Bound))
	}

	delta := func(n *node.N) heuristic.H {
		return heuristic.Heuristic(bhr.Union(bound, n.Bound)) - heuristic.Heuristic(n.Bound)
	}

	// deferred is the additional heuristic penalty caused by potentially
	// expanding the current node's bounds.
	deferred := 2 * delta(n)

	return heuristic.H(
		math.Min(
			float64(delta(n.Left)+deferred),
			float64(delta(n.Right)+deferred),
		),
	)
}

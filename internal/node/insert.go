package node

import (
	"github.com/downflux/go-bvh/internal/heuristic"
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-pq/pq"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
)

func Insert(root *N, id point.ID, bound hyperrectangle.R) {
	// Find best new sibling for the new leaf.
	q := pq.New[*N](0)
	q.Push(root, -float64(heuristic.Heuristic(bhr.Union(root.Bound, bound))))

	candidate := root

	for n := q.Pop(); n != nil || !q.Empty(); n = q.Pop() {
		direct := heuristic.Heuristic(bhr.Union(bound, n.Bound))

		inherited := heuristic.H(0.0)
		for m := n.Parent; m.Parent != nil; m = m.Parent {
			inherited += heuristic.Heuristic(bhr.Union(bound, m.Bound)) - heuristic.Heuristic(bound)
		}

		if float64(direct+inherited) < -q.Priority() {
			candidate = n
		}

		h := heuristic.Heuristic(bound) + (heuristic.Heuristic(bhr.Union(bound, n.Bound)) - heuristic.Heuristic(n.Bound))
		if float64(h) < -q.Priority() {
			q.Push(n.Left, -float64(h))
			q.Push(n.Right, -float64(h))
		}
	}

	// Create new parent.
	parent := &N{
		Parent: candidate.Parent,
		Left: &N{
			ID:    id,
			Bound: bound,
		},
		Right: candidate,
		Bound: bhr.Union(bound, candidate.Bound),
	}
	if candidate.Parent != nil {
		if candidate.Parent.Left == candidate {
			candidate.Parent.Left = parent
		} else {
			candidate.Parent.Right = parent
		}
	}
	candidate.Parent = parent
	parent.Left.Parent = parent

	// Refit AABBs.
	for m := parent.Parent; m.Parent != nil; m = m.Parent {
		m.Bound = bhr.Union(bound, m.Bound)
	}
}

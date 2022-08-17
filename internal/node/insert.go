package node

import (
	"fmt"

	"github.com/downflux/go-bvh/internal/allocation"
	"github.com/downflux/go-bvh/internal/heuristic"
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-pq/pq"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
)

// Insert adds the given point into the tree. If a new node is created, it will
// be created with a new index.
//
// Insert will return the new root.
func Insert(nodes allocation.C[*N], root *N, id point.ID, bound hyperrectangle.R) *N {
	if root == nil {
		nid := nodes.Allocate()
		n := New(O{
			ID:    id,
			Index: nid,
			Bound: bound,
		})
		if err := nodes.Insert(nid, n); err != nil {
			panic(fmt.Sprintf("cannot insert node: %s", err))
		}
		return n
	}

	// Find best new sibling for the new leaf.
	q := pq.New[*N](0)
	q.Push(root, -float64(heuristic.Heuristic(bhr.Union(root.Bound(), bound))))

	candidate := root

	for n := q.Pop(); n != nil || !q.Empty(); n = q.Pop() {
		direct := heuristic.Heuristic(bhr.Union(bound, n.Bound()))

		inherited := heuristic.H(0.0)
		for m := Parent(nodes, n); Parent(nodes, m) != nil; m = Parent(nodes, m) {
			inherited += heuristic.Heuristic(bhr.Union(bound, m.Bound())) - heuristic.Heuristic(bound)
		}

		if float64(direct+inherited) < -q.Priority() {
			candidate = n
		}

		h := heuristic.Heuristic(bound) + (heuristic.Heuristic(bhr.Union(bound, n.Bound())) - heuristic.Heuristic(n.Bound()))
		if float64(h) < -q.Priority() {
			q.Push(Left(nodes, n), -float64(h))
			q.Push(Right(nodes, n), -float64(h))
		}
	}

	// Create new parent.
	pid := nodes.Allocate()
	nid := nodes.Allocate()
	n := New(O{
		ID: id,

		Index:  nid,
		Parent: pid,

		Bound: bound,
	})
	p := New(O{
		Index:  pid,
		Parent: Parent(nodes, candidate).Index(),
		Left:   nid,
		Right:  candidate.Index(),

		Bound: bhr.Union(bound, candidate.Bound()),
	})
	nodes.Insert(nid, n)
	nodes.Insert(pid, p)
	if Parent(nodes, candidate) != nil {
		if Left(nodes, Parent(nodes, candidate)) == candidate {
			Parent(nodes, candidate).left = pid
		} else {
			Parent(nodes, candidate).right = pid
		}
	}
	candidate.parent = pid
	Left(nodes, p).parent = pid

	// Refit AABBs.
	var m *N
	for m = Parent(nodes, p); Parent(nodes, m) != nil; m = Parent(nodes, m) {
		m.bound = bhr.Union(bound, m.Bound())
	}

	return m
}

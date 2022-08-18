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

// candidate finds the node to which an object with the given bound will be
// inserted. This is based on the branch-and-bound algorithm (Catto 2019).
func candidate(nodes allocation.C[*N], root *N, bound hyperrectangle.R) *N {
	q := pq.New[*N](0)
	q.Push(root, -float64(heuristic.Heuristic(bhr.Union(root.Bound(), bound))))

	c := root

	for q.Len() > 0 {
		n := q.Pop()

		direct := heuristic.Heuristic(bhr.Union(bound, n.Bound()))

		inherited := heuristic.H(0.0)
		for m := Parent(nodes, n); m != nil && Parent(nodes, m) != nil; m = Parent(nodes, m) {
			inherited += heuristic.Heuristic(bhr.Union(bound, m.Bound())) - heuristic.Heuristic(bound)
		}

		if float64(direct+inherited) < -q.Priority() {
			c = n
		}

		h := heuristic.Heuristic(bound) + (heuristic.Heuristic(bhr.Union(bound, n.Bound())) - heuristic.Heuristic(n.Bound()))
		if float64(h) < -q.Priority() {
			q.Push(Left(nodes, n), -float64(h))
			q.Push(Right(nodes, n), -float64(h))
		}
	}
	return c
}

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
	c := candidate(nodes, root, bound)

	// Create new parent.
	pid := nodes.Allocate()
	nid := nodes.Allocate()
	n := New(O{
		ID: id,

		Index:  nid,
		Parent: pid,

		Bound: bound,
	})
	var aid allocation.ID
	if Parent(nodes, c) != nil {
		aid = Parent(nodes, c).Index()
	}
	p := New(O{
		Index:  pid,
		Parent: aid,
		Left:   nid,
		Right:  c.Index(),

		Bound: bhr.Union(bound, c.Bound()),
	})
	nodes.Insert(nid, n)
	nodes.Insert(pid, p)
	if Parent(nodes, c) != nil {
		if Left(nodes, Parent(nodes, c)) == c {
			Parent(nodes, c).left = pid
		} else {
			Parent(nodes, c).right = pid
		}
	}
	c.parent = pid
	Left(nodes, p).parent = pid

	// Walk back up the tree refitting AABBs and applying rotations.
	var m *N
	for m = p; Parent(nodes, m) != nil; m = Parent(nodes, m) {
		m.bound = bhr.Union(bound, m.Bound())

		// TODO(minkezhang): Apply rotation.
	}

	return m
}

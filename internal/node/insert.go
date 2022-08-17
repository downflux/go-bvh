package node

import (
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
func Insert(root *N, id point.ID, bound hyperrectangle.R) *N {
	if root == nil {
		ns := Nodes(make(map[Index]*N))
		nid := ns.Allocate()
		n := New(O{
			ID: id,

			Nodes: ns,
			Index: nid,

			Bound: bound,
		})
		n.nodes.Insert(nid, n)
		return n
	}

	// Find best new sibling for the new leaf.
	q := pq.New[*N](0)
	q.Push(root, -float64(heuristic.Heuristic(bhr.Union(root.Bound(), bound))))

	candidate := root

	for n := q.Pop(); n != nil || !q.Empty(); n = q.Pop() {
		direct := heuristic.Heuristic(bhr.Union(bound, n.Bound()))

		inherited := heuristic.H(0.0)
		for m := n.Parent(); m.Parent() != nil; m = m.Parent() {
			inherited += heuristic.Heuristic(bhr.Union(bound, m.Bound())) - heuristic.Heuristic(bound)
		}

		if float64(direct+inherited) < -q.Priority() {
			candidate = n
		}

		h := heuristic.Heuristic(bound) + (heuristic.Heuristic(bhr.Union(bound, n.Bound())) - heuristic.Heuristic(n.Bound()))
		if float64(h) < -q.Priority() {
			q.Push(n.Left(), -float64(h))
			q.Push(n.Right(), -float64(h))
		}
	}

	// Create new parent.
	pid := root.nodes.Allocate()
	nid := root.nodes.Allocate()
	n := New(O{
		ID: id,

		Nodes:  root.nodes,
		Index:  nid,
		Parent: pid,

		Bound: bound,
	})
	p := New(O{
		Nodes:  root.nodes,
		Index:  pid,
		Parent: candidate.Parent().Index(),
		Left:   nid,
		Right:  candidate.Index(),

		Bound: bhr.Union(bound, candidate.Bound()),
	})
	root.nodes.Insert(nid, n)
	root.nodes.Insert(pid, p)
	if candidate.Parent() != nil {
		if candidate.Parent().Left() == candidate {
			candidate.Parent().left = pid
		} else {
			candidate.Parent().right = pid
		}
	}
	candidate.parent = pid
	p.Left().parent = pid

	// Refit AABBs.
	var m *N
	for m = p.Parent(); m.Parent() != nil; m = m.Parent() {
		m.bound = bhr.Union(bound, m.Bound())
	}

	return m
}

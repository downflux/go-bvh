package insert

import (
	"fmt"

	"github.com/downflux/go-bvh/internal/allocation"
	"github.com/downflux/go-bvh/internal/heuristic"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-pq/pq"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
)

// candidate finds the node to which an object with the given bound will be
// inserted. This is based on the branch-and-bound algorithm (Catto 2019).
func candidate(nodes allocation.C[*node.N], root *node.N, bound hyperrectangle.R) *node.N {
	q := pq.New[*node.N](0)

	// Note that the priority queue is a max-heap, so we will need to flip
	// the heuristic signs.
	q.Push(root, -float64(heuristic.Direct(root, bound)))
	c := root

	for q.Len() > 0 {
		n := q.Pop()

		inherited := heuristic.Inherited(nodes, n, bound)

		// Check if the current node is a better insertion candidate.
		if float64(heuristic.Direct(n, bound)+inherited) < -q.Priority() {
			c = n
		}

		// Append queue children to the queue if the lower bound for
		// inserting into the child is less than the current minimum
		// (i.e. there's room for optimization).
		//
		// Note that the inherited heuristic is the same between the
		// left and right children.
		inherited += heuristic.Heuristic(bhr.Union(bound, n.Bound())) - heuristic.Heuristic(bound)
		h := heuristic.Heuristic(bound) + inherited
		if float64(h) < -q.Priority() {
			q.Push(node.Left(nodes, n), -float64(h))
			q.Push(node.Right(nodes, n), -float64(h))
		}
	}
	return c
}

// Insert adds the given point into the tree. If a new node is created, it will
// be created with a new index.
//
// Insert will return the new root.
func Insert(nodes allocation.C[*node.N], root *node.N, id point.ID, bound hyperrectangle.R) *node.N {
	if root == nil {
		nid := nodes.Allocate()
		n := node.New(node.O{
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
	n := node.New(node.O{
		ID: id,

		Index:  nid,
		Parent: pid,

		Bound: bound,
	})
	var aid allocation.ID
	if node.Parent(nodes, c) != nil {
		aid = node.Parent(nodes, c).Index()
	}
	p := node.New(node.O{
		Index:  pid,
		Parent: aid,
		Left:   nid,
		Right:  c.Index(),

		Bound: bhr.Union(bound, c.Bound()),
	})
	nodes.Insert(nid, n)
	nodes.Insert(pid, p)
	if node.Parent(nodes, c) != nil {
		if node.Left(nodes, node.Parent(nodes, c)) == c {
			node.Parent(nodes, c).SetLeft(pid)
		} else {
			node.Parent(nodes, c).SetRight(pid)
		}
	}
	c.SetParent(pid)
	node.Left(nodes, p).SetParent(pid)

	// Walk back up the tree refitting AABBs and applying rotations.
	var m *node.N
	for m = p; node.Parent(nodes, m) != nil; m = node.Parent(nodes, m) {
		m.SetBound(bhr.Union(bound, m.Bound()))

		// TODO(minkezhang): Apply rotation.
	}

	return m
}

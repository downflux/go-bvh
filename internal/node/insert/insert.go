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

// Execute adds the given point into the tree. If a new node is created, it will
// be created with a new index.
//
// This function will return the new root.
func Execute(nodes allocation.C[*node.N], root allocation.ID, id point.ID, bound hyperrectangle.R) allocation.ID {
	if _, ok := nodes[root]; !ok {
		nid := nodes.Allocate()
		n := node.New(node.O{
			ID:    id,
			Index: nid,
			Bound: bound,
		})
		if err := nodes.Insert(nid, n); err != nil {
			panic(fmt.Sprintf("cannot insert node: %s", err))
		}
		return n.Index()
	}

	// Find best new sibling for the new leaf.
	cid := findCandidate(nodes, root, bound)

	// Create new parent.
	pid := createParent(nodes, cid, id, bound)

	// Walk back up the tree refitting AABBs and applying rotations, and
	// find the new root.
	var m *node.N
	for m = nodes[pid]; node.Parent(nodes, m) != nil; m = node.Parent(nodes, m) {
		m.SetBound(bhr.Union(bound, m.Bound()))

		// TODO(minkezhang): Apply rotation.
	}

	return m.Index()
}

// createParent creates a new parent node for a candidate r. This parent will
// have have the old node and a newly-created node with the given bounds.
//
// This function will modify the allocation table as a side-effect.
func createParent(nodes allocation.C[*node.N], rid allocation.ID, id point.ID, bound hyperrectangle.R) allocation.ID {
	r := nodes[rid]

	pid := nodes.Allocate()
	lid := nodes.Allocate()

	l := node.New(node.O{
		ID: id,

		Index:  lid,
		Parent: pid,

		Bound: bound,
	})

	var aid allocation.ID
	if node.Parent(nodes, r) != nil {
		aid = node.Parent(nodes, r).Index()
	}

	p := node.New(node.O{
		Index:  pid,
		Parent: aid,
		Left:   lid,
		Right:  rid,

		Bound: bhr.Union(bound, r.Bound()),
	})
	nodes.Insert(lid, l)
	nodes.Insert(pid, p)
	if node.Parent(nodes, r) != nil {
		if node.Left(nodes, node.Parent(nodes, r)) == r {
			node.Parent(nodes, r).SetLeft(pid)
		} else {
			node.Parent(nodes, r).SetRight(pid)
		}
	}
	r.SetParent(pid)
	node.Left(nodes, p).SetParent(pid)

	return pid
}

// findCandidate finds the node to which an object with the given bound will be
// inserted. This is based on the branch-and-bound algorithm (Catto 2019).
func findCandidate(nodes allocation.C[*node.N], rid allocation.ID, bound hyperrectangle.R) allocation.ID {
	root, ok := nodes[rid]
	if !ok {
		panic("cannot find candidate for an empty root node")
	}

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
		//
		// N.B.: The full expression for the node expansion heuristic is
		//
		//   inherited += heuristic.Heuristic(bhr.Union(bound, n.Bound())) - heuristic.Heuristic(bound)
		//   h := heuristic.Heuristic(bound) + inherited
		//
		// Note that the bounding heuristic cancels out.
		h := inherited + heuristic.Heuristic(bhr.Union(bound, n.Bound()))
		if float64(h) < -q.Priority() {
			q.Push(node.Left(nodes, n), -float64(h))
			q.Push(node.Right(nodes, n), -float64(h))
		}
	}
	return c.Index()
}

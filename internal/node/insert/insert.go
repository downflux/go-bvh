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

type Inserter struct {
	a allocation.C[*node.N]
}

// Execute adds the given point into the tree. If a new node is created, it will
// be created with a new index.
//
// This function will return the new root.
func (c Inserter) Execute(root allocation.ID, id point.ID, bound hyperrectangle.R) allocation.ID {
	if _, ok := c.a[root]; !ok {
		nid := c.a.Allocate()
		n := node.New(node.O{
			ID:    id,
			Index: nid,
			Bound: bound,
		})
		if err := c.a.Insert(nid, n); err != nil {
			panic(fmt.Sprintf("cannot insert node: %s", err))
		}
		return n.Index()
	}

	// Find best new sibling for the new leaf.
	cid := findCandidate(c, root, bound)

	// Create new parent.
	pid := c.createParent(cid, id, bound)

	// Walk back up the tree refitting AABBs and applying rotations, and
	// find the new root.
	var m *node.N
	for m = c.a[pid]; node.Parent(c.a, m) != nil; m = node.Parent(c.a, m) {
		m.SetBound(bhr.Union(bound, m.Bound()))

		// TODO(minkezhang): Apply rotation.
	}

	return m.Index()
}

// createParent creates a new parent node for a candidate r. This parent will
// have have the old node and a newly-created node with the given bounds.
//
// This function will modify the allocation table as a side-effect.
func (c Inserter) createParent(rid allocation.ID, id point.ID, bound hyperrectangle.R) allocation.ID {
	r := c.a[rid]

	pid := c.a.Allocate()
	lid := c.a.Allocate()

	l := node.New(node.O{
		ID: id,

		Index:  lid,
		Parent: pid,

		Bound: bound,
	})

	var aid allocation.ID
	if node.Parent(c.a, r) != nil {
		aid = node.Parent(c.a, r).Index()
	}

	p := node.New(node.O{
		Index:  pid,
		Parent: aid,
		Left:   lid,
		Right:  rid,

		Bound: bhr.Union(bound, r.Bound()),
	})
	c.a.Insert(lid, l)
	c.a.Insert(pid, p)
	if node.Parent(c.a, r) != nil {
		if node.Left(c.a, node.Parent(c.a, r)) == r {
			node.Parent(c.a, r).SetLeft(pid)
		} else {
			node.Parent(c.a, r).SetRight(pid)
		}
	}
	r.SetParent(pid)
	node.Left(c.a, p).SetParent(pid)

	return pid
}

// findCandidate finds the node to which an object with the given bound will be
// inserted. This is based on the branch-and-bound algorithm (Catto 2019).
func findCandidate(inserter Inserter, rid allocation.ID, bound hyperrectangle.R) allocation.ID {
	root, ok := inserter.a[rid]
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

		inherited := heuristic.Inherited(inserter.a, n, bound)

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
			q.Push(node.Left(inserter.a, n), -float64(h))
			q.Push(node.Right(inserter.a, n), -float64(h))
		}
	}
	return c.Index()
}

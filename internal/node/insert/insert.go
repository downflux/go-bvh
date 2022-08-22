package insert

import (
	"fmt"

	"github.com/downflux/go-bvh/internal/allocation"
	"github.com/downflux/go-bvh/internal/allocation/id"
	"github.com/downflux/go-bvh/internal/heuristic"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/rotate"
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-pq/pq"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
)

// Execute adds the given point into the tree. If a new node is created, it will
// be created with a new index.
//
// This function will return the new root.
func Execute(nodes allocation.C[*node.N], root id.ID, id point.ID, bound hyperrectangle.R) id.ID {
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
	cid := findSibling(nodes, root, bound)

	// Create new parent.
	pid := createParent(nodes, cid, id, bound)

	// Rebalance the tree up to the root.
	rotate.Execute(nodes, pid)

	return pid
}

// createParent creates a new parent node for a candidate r. This parent will
// have have the old node and a newly-created node with the given bounds. The
// sibling ID rid must exist and refer to a non-nil node.
//
// This function will modify the allocation table as a side-effect.
//
// This function returns the id.ID of the newly-created parent.
func createParent(nodes allocation.C[*node.N], rid id.ID, i point.ID, bound hyperrectangle.R) id.ID {
	r := nodes[rid]
	if r == nil {
		panic(fmt.Sprintf("given right id.ID does not exist: %v", rid))
	}

	pid := nodes.Allocate()
	lid := nodes.Allocate()

	l := node.New(node.O{
		ID: i,

		Index:  lid,
		Parent: pid,

		Bound: bound,
	})

	var aid id.ID
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

	// Walk back up the tree refitting AABBs.
	for n := nodes[pid]; n != nil; n = node.Parent(nodes, n) {
		n.SetBound(bhr.Union(bound, n.Bound()))
	}

	return pid
}

// findSibling finds the node to which an object with the given bound will
// siblings. A new parent node will be created above both the sibling and the
// input bound. This is based on the branch-and-bound algorithm (Catto 2019).
func findSibling(nodes allocation.C[*node.N], rid id.ID, bound hyperrectangle.R) id.ID {
	root, ok := nodes[rid]
	if !ok {
		panic("cannot find candidate for an empty root node")
	}

	q := pq.New[*node.N](0)

	// Note that the priority queue is a max-heap, so we will need to flip
	// the heuristic signs.
	q.Push(root, -float64(heuristic.Direct(root, bound)))

	c := root
	d := -q.Priority()

	for q.Len() > 0 {
		n := q.Pop()

		inherited := heuristic.Inherited(nodes, n, bound)

		// Check if the current node is a better insertion candidate.
		actual := float64(heuristic.Direct(n, bound) + inherited)
		if actual < d {
			c = n
			d = actual
		}

		if !node.Leaf(nodes, n) {
			// Append queue children to the queue if the lower bound for
			// inserting into the child is less than the current minimum
			// (i.e. there's room for optimization).
			//
			// Note that the inherited heuristic is the same between the
			// left and right children.
			inherited += heuristic.Heuristic(bhr.Union(bound, n.Bound())) - heuristic.Heuristic(n.Bound())
			estimated := heuristic.Heuristic(bound) + inherited

			if float64(estimated) < d {
				q.Push(node.Left(nodes, n), -float64(estimated))
				q.Push(node.Right(nodes, n), -float64(estimated))
			}
		}
	}
	return c.Index()
}

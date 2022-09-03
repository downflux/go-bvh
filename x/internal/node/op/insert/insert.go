package insert

import (
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-bvh/x/internal/node/op/insert/heuristic"
	"github.com/downflux/go-bvh/x/internal/node/op/insert/insert"
	"github.com/downflux/go-bvh/x/internal/node/op/rotate"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-pq/pq"
)

// Execute adds a new node with the given data into the tree. The returned node
// is the new root of the tree.
func Execute(root *node.N, x id.ID, aabb hyperrectangle.R) *node.N {
	var c *node.C
	if root == nil {
		c = node.Cache()
	} else {
		c = root.Cache()
	}

	m := node.New(node.O{
		Nodes: c,
		Data: map[id.ID]hyperrectangle.R{
			x: aabb,
		},
	})
	if root == nil {
		return m
	}

	return rotate.Execute(insert.Execute(sibling(root, aabb), m))
}

// sibling finds the node to which an object with the given bound will siblings.
// A new parent node will be created above both the sibling and the input bound.
// This is based on the branch-and-bound algorithm (Catto 2019).
//
// The hyperrectangle input is the AABB of the new prospective node.
func sibling(root *node.N, aabb hyperrectangle.R) *node.N {
	if root == nil {
		panic("cannot find sibling candidate for an empty root node")
	}

	q := pq.New[*node.N](0)

	// Note that the priority queue is a max-heap, so we will need to flip
	// the heuristic signs.
	q.Push(root, -heuristic.Actual(root, aabb))

	c := root
	d := -q.Priority()

	for q.Len() > 0 {
		n := q.Pop()

		// Check if the current node is a better insertion candidate.
		if actual := heuristic.Actual(n, aabb); actual < d {
			c = n
			d = actual
		}

		if !n.IsLeaf() {
			// Append queue children to the queue if the lower bound for
			// inserting into the child is less than the current minimum
			// (i.e. there's room for optimization).
			//
			// Note that the inherited heuristic is the same between the
			// left and right children.
			if estimate := heuristic.Estimate(n, aabb); estimate < d {
				q.Push(n.Left(), -estimate)
				q.Push(n.Right(), -estimate)
			}
		}
	}

	return c
}

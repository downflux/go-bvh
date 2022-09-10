// Package pq returns a candidate sibling node for the given tree. This
// implementation is based directly on the Catto 2019 slides, which evaluates
// the iterative bounding boxes and puts them into a priority queue.
package pq

import (
	"math"

	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/op/insert/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-pq/pq"
)

// Execute finds the node to which an object with the given bound will siblings.
// If the returned sibling is a leaf node, then the caller should attempt to
// insert the hyperrectangle into the returned node. If, however, the returned
// node is an internal node, then the caller should instead create a new node
// with the returned node as a direct sibling. This is based on the
// branch-and-bound algorithm from Catto 2019.
//
// See https://en.wikipedia.org/wiki/Branch_and_bound for more information on
// the generic algorithm.
//
// The hyperrectangle input is the AABB of the new prospective node.
func Execute(root *node.N, aabb hyperrectangle.R) *node.N {
	if root == nil {
		panic("cannot find sibling candidate for an empty root node")
	}

	q := pq.New[*node.N](0)

	// Note that the priority queue is a max-heap, so we will need to flip
	// the heuristic signs.
	q.Push(root, -heuristic.B(root, aabb))

	var c *node.N
	f := math.Inf(0)

	for q.Len() > 0 {
		n := q.Pop()

		// Check if the current node is a better insertion candidate.
		if actual := heuristic.F(n, aabb); actual < f {
			c = n
			f = actual
		}

		if !n.IsLeaf() {
			// Append queue children to the queue if the lower bound for
			// inserting into the child is less than the current minimum
			// (i.e. there's room for optimization).
			//
			// Note that the inherited heuristic is the same between the
			// left and right children.
			if estimate := heuristic.B(n, aabb); estimate < f {
				q.Push(n.Left(), -estimate)
				q.Push(n.Right(), -estimate)
			}
		}
	}

	return c
}

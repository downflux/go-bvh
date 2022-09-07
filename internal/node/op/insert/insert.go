package insert

import (
	"math"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/op/insert/heuristic"
	"github.com/downflux/go-bvh/internal/node/op/insert/insert"
	"github.com/downflux/go-bvh/internal/node/op/insert/split"
	"github.com/downflux/go-bvh/internal/node/op/rotate"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-pq/pq"
)

// Execute adds a new node with the given data into the tree. The returned node
// is the newly-created node.
func Execute(root *node.N, size int, x id.ID, aabb hyperrectangle.R) *node.N {
	var c *node.C
	if root == nil {
		c = node.Cache()
	} else {
		c = root.Cache()
	}

	if root == nil {
		return node.New(node.O{
			Nodes: c,
			Data: map[id.ID]hyperrectangle.R{
				x: aabb,
			},
			Size: size,
		})
	}

	// m is the newly-created leaf node containing the input data.
	var m *node.N

	s := sibling(root, aabb)
	if s.IsLeaf() {
		if !s.IsFull() {
			m = s
		} else {
			m = split.Execute(s, split.RandomPartition)
		}
		m.Insert(x, aabb)
	} else {
		m = node.New(node.O{
			Nodes: c,
			Data: map[id.ID]hyperrectangle.R{
				x: aabb,
			},
			Size: size,
		})
	}

	// Add a shared parent between the sibling and newly-created node. Note
	// that we will skip this step if the sibling was a non-full leaf.
	if s != m {
		insert.Execute(s, m)
	}

	// m is now linked to the correct parent; we need to balance the tree.
	if !m.IsRoot() {
		rotate.Execute(m.Parent())
	}

	return m
}

// sibling finds the node to which an object with the given bound will siblings.
// A new parent node will be created above both the sibling and the input bound.
// This is based on the branch-and-bound algorithm (Catto 2019).
//
// See https://en.wikipedia.org/wiki/Branch_and_bound.
//
// The hyperrectangle input is the AABB of the new prospective node.
func sibling(root *node.N, aabb hyperrectangle.R) *node.N {
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

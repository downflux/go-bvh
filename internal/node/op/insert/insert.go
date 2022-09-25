package insert

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/balance"
	"github.com/downflux/go-bvh/internal/node/op/insert/insert"
	// "github.com/downflux/go-bvh/internal/node/op/insert/sibling"
	sibling "github.com/downflux/go-bvh/internal/node/op/insert/sibling/greedy"
	"github.com/downflux/go-bvh/internal/node/op/insert/split"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

var (
	// FindSibling is a function which returns the nearest node for a given
	// AABB insertion candidate.
	FindSibling = sibling.Execute
)

// Execute adds a new node with the given data into the tree. The returned node
// is the newly-created node.
func Execute(root *node.N, size uint, x id.ID, aabb hyperrectangle.R) *node.N {
	if root == nil {
		return node.New(node.O{
			Nodes: node.Cache(),
			Data: map[id.ID]hyperrectangle.R{
				x: aabb,
			},
			Size: size,
			K:    aabb.Min().Dimension(),
		})
	}

	// m is the newly-created leaf node containing the input data.
	var m *node.N

	s := FindSibling(root, aabb)
	// If a leaf is returned, we should attempt to insert the object into
	// this leaf if possible -- the reasoning here is that the overall
	// heuristic for inserting into a leaf is lower than creating a new
	// leaf.
	if s.IsLeaf() {
		if !s.IsFull() {
			m = s
		} else {
			// If the leaf is full, we will create a new leaf, and
			// move some of the existing elements of s into the new
			// leaf. This is so we can minimize the total surface
			// area heuristic between the two nodes.
			m = split.Execute(s, split.RandomPartition)
		}
		m.Insert(x, aabb)
	} else {
		m = node.New(node.O{
			Nodes: root.Cache(),
			Data: map[id.ID]hyperrectangle.R{
				x: aabb,
			},
			Size: size,
			K:    s.K(),
		})
	}

	// Add a shared parent between the sibling and newly-created node. Note
	// that we will skip this step if the sibling was a non-full leaf.
	if s != m {
		insert.Execute(s, m)
	}

	// m is now linked to the correct parent; we need to balance the tree.
	balance.Execute(m)

	return m
}

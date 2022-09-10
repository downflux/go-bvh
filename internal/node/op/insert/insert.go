package insert

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/op/insert/insert"
	"github.com/downflux/go-bvh/internal/node/op/insert/sibling/pq"
	"github.com/downflux/go-bvh/internal/node/op/insert/split"
	"github.com/downflux/go-bvh/internal/node/op/rotate"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

// Execute adds a new node with the given data into the tree. The returned node
// is the newly-created node.
func Execute(root *node.N, size uint, x id.ID, aabb hyperrectangle.R) *node.N {
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

	s := pq.Execute(root, aabb)
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

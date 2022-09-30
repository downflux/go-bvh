package remove

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/op/remove/remove"
	"github.com/downflux/go-bvh/internal/node/op/rotate/balance"
)

// Execute removes a leaf node from a BVH node. The returned node is the new root.
func Execute(n *node.N, id id.ID) *node.N {
	if !n.IsLeaf() {
		panic("cannot remove an internal node directly")
	}

	n.Remove(id)
	if n.IsEmpty() {
		m := remove.Execute(n)
		if m != nil && !m.IsRoot() {
			balance.Execute(m)
		}
		if m != nil {
			return m.Root()
		}
		return m
	}
	return n.Root()
}

// Package node is an internal-only node implementation struct, and its
// properties and data points should only be accessed via the operations API in
// the /internal/node/op/ directory.
package node

import (
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/x/hyperrectangle"
)

type D struct {
	ID   id.ID
	AABB hyperrectangle.R
}

type N struct {
	parent *N
	left   *N
	right  *N

	data             []D
	aabbCacheIsValid bool
	aabbCache        hyperrectangle.R
}

func (n *N) InvalidateAABBCache() {
	n.aabbCacheIsValid = false
}

func (n *N) clean() {
	n.parent = nil
	n.left = nil
	n.right = nil
}

func (n *N) Left() *N { return n.left }
func (n *N) SetLeft(m *N) {
	n.Left().clean()

	n.left = m
	m.parent = n
}

func (n *N) Right() *N { return n.right }
func (n *N) SetRight(m *N) {
	n.Right().clean()

	n.right = m
	m.parent = n
}

func (n *N) Parent() *N { return n.parent }

func (n *N) Root() *N {
	if n.Parent() == nil {
		return n
	}
	return n.Parent().Root()
}

func (n *N) IsLeaf() bool { return len(n.data) > 0 }
func (n *N) IsRoot() bool { return n.Parent() == nil }
func (n *N) AABB() hyperrectangle.R {
	if n.aabbCacheIsValid {
		return n.aabbCache
	}

	n.aabbCacheIsValid = true
	if n.IsLeaf() {
		rs := make([]hyperrectangle.R, len(n.data))
		for i := 0; i < len(n.data); i++ {
			rs[i] = n.data[i].AABB
		}
		n.aabbCache = bhr.AABB(rs)
	} else {
		n.aabbCache = bhr.Union(n.Left().AABB(), n.Right().AABB())
	}

	return n.aabbCache
}

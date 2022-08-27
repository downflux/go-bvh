// Package node is an internal-only node implementation struct, and its
// properties and data points should only be accessed via the operations API in
// the /internal/node/op/ directory.
package node

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/x/hyperrectangle"
)

type D[T comparable] struct {
	ID   T
	AABB hyperrectangle.R
}

type N[T comparable] struct {
	parent *N[T]
	left   *N[T]
	right  *N[T]

	data             []D[T]
	aabbCacheIsValid bool
	aabbCache        hyperrectangle.R
}

func (n *N[T]) InvalidateAABBCache() {
	n.aabbCacheIsValid = false
}

func (n *N[T]) clean() {
	n.parent = nil
	n.left = nil
	n.right = nil
}

func (n *N[T]) Left() *N[T] { return n.left }
func (n *N[T]) SetLeft(m *N[T]) {
	n.Left().clean()

	n.left = m
	m.SetParent(n)
}

func (n *N[T]) Right() *N[T] { return n.right }
func (n *N[T]) SetRight(m *N[T]) {
	n.Right().clean()

	n.right = m
	m.SetParent(n)
}

func (n *N[T]) Parent() *N[T] { return n.parent }

func (n *N[T]) Root() *N[T] {
	if n.Parent() == nil {
		return n
	}
	return n.Parent().Root()
}

func (n *N[T]) Leaf() bool { return len(n.data) > 0 }
func (n *N[T]) AABB() hyperrectangle.R {
	if n.aabbCacheIsValid {
		return n.aabbCache
	}

	n.aabbCacheIsValid = true
	if n.Leaf() {
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

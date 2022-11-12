// Package node defines a node interface. This is used by both the cache and
// node implementations to avoid cyclic imports.
package node

import (
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache/branch"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

func Equal(n N, m N) bool {
	if n == nil && m == nil {
		return true
	}

	if n == nil && m != nil {
		return false
	}

	if n.ID() != m.ID() {
		return false
	}

	if n.IsRoot() != m.IsRoot() {
		return false
	}

	if n.Height() != m.Height() {
		return false
	}

	if !n.IsRoot() {
		if (n.Parent() == nil && m.Parent() != nil) || (n.Parent() != nil && m.Parent() == nil) {
			return false
		}

		if n.Parent() != nil {
			if n.Parent().ID() != m.Parent().ID() {
				return false
			}
		}
	}

	if n.IsLeaf() != m.IsLeaf() {
		return false
	}

	if !n.IsLeaf() {
		if !Equal(n.Left(), m.Left()) || !Equal(n.Right(), m.Right()) {
			return false
		}
	}

	if n.IsLeaf() {
		if len(n.Leaves()) != len(m.Leaves()) {
			return false
		}

		for k := range n.Leaves() {
			if _, ok := m.Leaves()[k]; !ok {
				return false
			}
		}
	}

	if !hyperrectangle.Within(n.AABB().R(), m.AABB().R()) {
		return false
	}

	return true
}

type N interface {
	ID() cid.ID

	IsRoot() bool

	Height() int
	SetHeight(h int)

	IsLeaf() bool
	IsFull() bool
	Leaves() map[id.ID]struct{}

	AABB() hyperrectangle.M

	Child(b branch.B) N
	SetChild(b branch.B, x cid.ID)

	Branch(x cid.ID) branch.B

	Parent() N
	Left() N
	Right() N

	SetParent(x cid.ID)
	SetLeft(x cid.ID)
	SetRight(x cid.ID)
}

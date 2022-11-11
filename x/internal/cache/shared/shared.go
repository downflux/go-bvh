// Package shared defines a node interface. This is used by both the cache and
// node implementations to avoid cyclic imports.
package shared

import (
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache/branch"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

func Equal(n N, m N) bool {
	if n.ID() != m.ID() {
		return false
	}

	if n.IsRoot() != m.IsRoot() {
		return false
	}

	if n.IsLeaf() != m.IsRoot() {
		return false
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

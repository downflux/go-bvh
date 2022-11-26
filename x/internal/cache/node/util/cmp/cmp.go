package cmp

import (
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-geometry/epsilon"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

var (
	DefaultF = F{
		Height:    true,
		AABB:      true,
		Heuristic: true,
	}
)

// F is a comparison filter used to optionally exclude comparing cached
// values. This is useful for specific types of testing where these cache values
// do not matter.
type F struct {
	Height    bool
	AABB      bool
	Heuristic bool
}

func Equal(n node.N, m node.N) bool { return DefaultF.Equal(n, m) }

func (f F) Equal(n node.N, m node.N) bool {
	if n == nil && m == nil {
		return true
	}

	if (n == nil && m != nil) || (n != nil && m == nil) {
		return false
	}

	if n.ID() != m.ID() {
		return false
	}

	if n.IsRoot() != m.IsRoot() {
		return false
	}

	if f.Height && n.Height() != m.Height() {
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
		if !f.Equal(n.Left(), m.Left()) || !f.Equal(n.Right(), m.Right()) {
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

	if f.AABB && !hyperrectangle.Within(n.AABB().R(), m.AABB().R()) {
		return false
	}

	if f.Heuristic && !epsilon.Within(n.Heuristic(), m.Heuristic()) {
		return false
	}

	return true
}

// Package node defines a node interface. This is used by both the cache and
// node implementations to avoid cyclic imports.
package node

import (
	"math"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache/branch"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

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

// SetHeight will update the height of a node. The input node must have valid
// and up-to-date child nodes.
func SetHeight(n N) {
	if n.IsLeaf() {
		n.SetHeight(0)
	} else {
		n.SetHeight(1 + int(
			math.Max(
				float64(n.Left().Height()),
				float64(n.Right().Height()),
			),
		))
	}
}

// SetAABB updates a node's AABB with the bounding boxes of its children. For a
// leaf node, this bounding box will have a buffer of some given expansion
// factor.
//
// The input node must be valid and up-to-date.
func SetAABB(n N, data map[id.ID]hyperrectangle.R, epsilon float64) {
	if !n.IsLeaf() {
		n.AABB().Copy(n.Left().AABB().R())
		n.AABB().Union(n.Right().AABB().R())
		return
	}

	// TODO(minkezhang): Improve performance by concurrently updating the
	// final AABB.
	var initialized bool
	var k vector.D
	for x := range n.Leaves() {
		if !initialized {
			n.AABB().Copy(data[x])
			k = data[x].Min().Dimension()
		} else {
			n.AABB().Union(data[x])
		}
	}
	n.AABB().Scale(math.Pow(epsilon, 1/float64(k)))
}

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

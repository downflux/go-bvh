// Package node is an internal-only node implementation struct, and its
// properties and data points should only be accessed via the operations API in
// the /internal/node/op/ directory.
package node

import (
	"fmt"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/x/hyperrectangle"
)

type D struct {
	ID   id.ID
	AABB hyperrectangle.R
}

type O struct {
	Left  *N
	Right *N
	Data  []D
}

type N struct {
	parent *N
	left   *N
	right  *N

	data             []D
	aabbCacheIsValid bool
	aabbCache        hyperrectangle.R
}

func Validate(n *N) error {
	if (n.left != nil && n.right == nil) || (n.right != nil && n.left == nil) {
		return fmt.Errorf("node has mismatching child nodes")
	}
	if n.left != nil && n.right != nil && len(n.data) > 0 {
		return fmt.Errorf("non-leaf node contains data")
	}
	if n.left == nil && n.right == nil && len(n.data) == 0 {
		return fmt.Errorf("leaf node contains no data")
	}
	if (n.left != nil && n.left.parent != n) || (n.right != nil && n.right.parent != n) {
		return fmt.Errorf("node child does not link to the parent")
	}
	if n.parent != nil && n.parent.left != n && n.parent.right != n {
		return fmt.Errorf("node parent does not link to child")
	}
	return nil
}

func New(o O) *N {
	n := &N{
		left:  o.Left,
		right: o.Right,
		data:  o.Data,
	}
	if n.left != nil {
		n.left.parent = n
	}
	if n.right != nil {
		n.right.parent = n
	}

	if err := Validate(n); err != nil {
		panic(fmt.Sprintf("cannot construct node: %v", err))
	}
	return n
}

func (n *N) InvalidateAABBCache() {
	n.aabbCacheIsValid = false
	if !n.IsRoot() {
		n.Parent().InvalidateAABBCache()
	}
}

func (n *N) Swap(m *N) {
	if !n.IsRoot() {
		if n.Parent().Left() == n {
			n.Parent().left = m
		} else {
			n.Parent().right = m
		}
		n.InvalidateAABBCache()
	}
	if !m.IsRoot() {
		if m.Parent().Left() == m {
			m.Parent().left = n
		} else {
			m.Parent().right = n
		}
		m.InvalidateAABBCache()
	}
	n.parent, m.parent = m.parent, n.parent
}

func (n *N) Left() *N   { return n.left }
func (n *N) Right() *N  { return n.right }
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

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

type O struct {
	Left  *N
	Right *N
	Data  map[id.ID]hyperrectangle.R
}

type N struct {
	parent *N
	left   *N
	right  *N

	data             map[id.ID]hyperrectangle.R
	aabbCacheIsValid bool
	aabbCache        hyperrectangle.R
}

func Validate(n *N) error {
	if n == nil {
		return nil
	}

	if (n.Left() != nil && n.Right() == nil) || (n.Right() != nil && n.Left() == nil) {
		return fmt.Errorf("node has mismatching child nodes")
	}
	if n.Left() != nil && n.Right() != nil && len(n.Data()) > 0 {
		return fmt.Errorf("non-leaf node contains data")
	}
	if n.Left() == nil && n.Right() == nil && len(n.Data()) == 0 {
		return fmt.Errorf("leaf node contains no data")
	}
	if (n.Left() != nil && n.Left().Parent() != n) || (n.Right() != nil && n.Right().Parent() != n) {
		return fmt.Errorf("child node does not link to the parent")
	}
	if n.Parent() != nil && n.Parent().Left() != n && n.Parent().Right() != n {
		return fmt.Errorf("parent node does not link to child")
	}
	return nil
}

func New(o O) *N {
	n := &N{
		left:  o.Left,
		right: o.Right,
		data:  o.Data,
	}
	if !n.IsLeaf() {
		n.Left().parent = n
		n.Right().parent = n
	}

	if err := Validate(n); err != nil {
		panic(fmt.Sprintf("cannot construct node: %v", err))
	}
	return n
}

func (n *N) InvalidateAABBCache() {
	// Since InvalidateAABBCache is called recursively up towards the root,
	// and AABB is calculated towards the leaf, if the cache is invalid at
	// some node, we are guaranteed all nodes above the current node are
	// also marked with an invalid cache. Skipping the tree iteration here
	// can reduce the complexity by a factor of O(log N) if we are
	// traveling up the tree anyway in some other algorithm.
	if !n.aabbCacheIsValid {
		return
	}

	n.aabbCacheIsValid = false
	if !n.IsRoot() {
		n.Parent().InvalidateAABBCache()
	}
}

// IsAncestor returns true if the current node n is an ancestor of m.
func (n *N) IsAncestor(m *N) bool {
	if n == m {
		return true
	}
	if m.IsRoot() {
		return false
	}
	return n.IsAncestor(m.Parent())
}

func (n *N) Insert(m *N) {
	if m == nil {
		panic("cannot insert an empty node")
	}
	if !m.IsRoot() {
		panic("cannot insert an internal node")
	}

	n.left, n.right = m, New(O{
		Left:  n.left,
		Right: n.right,
		Data:  n.data,
	})
	n.right.parent = n
	n.left.parent = n

	if err := Validate(n.Left()); err != nil {
		panic(fmt.Errorf("cannot insert node: %v", err))
	}
	if err := Validate(n.Right()); err != nil {
		panic(fmt.Errorf("cannot insert node: %v", err))
	}
}

func (n *N) Swap(m *N) {
	if m == nil {
		panic("cannot swap with empty node")
	}

	if n.IsAncestor(m) || m.IsAncestor(n) {
		panic("cannot swap a child node with its ancestor")
	}

	if !n.IsRoot() {
		p := n.Parent()
		if p.Left() == n {
			p.left = m
		} else {
			p.right = m
		}
		p.InvalidateAABBCache()
	}
	if !m.IsRoot() {
		p := m.Parent()
		if p.Left() == m {
			p.left = n
		} else {
			p.right = n
		}
		p.InvalidateAABBCache()
	}
	n.parent, m.parent = m.Parent(), n.Parent()
	if err := Validate(n); err != nil {
		panic(fmt.Errorf("could not swap nodes: %v", err))
	}
	if err := Validate(m); err != nil {
		panic(fmt.Errorf("could not swap nodes: %v", err))
	}
}

func (n *N) Data() map[id.ID]hyperrectangle.R { return n.data }
func (n *N) Left() *N                         { return n.left }

/*
func (n *N) SetLeft(m *N) {
	if m == nil {
		panic("cannot set an empty node as a child")
	}

	if !m.IsRoot() {
		panic("cannot add an internal node as a child")
	}

	if n.IsLeaf() {
		panic("cannot directly set a child on a leaf node")
	}

	q := n.Left()

	// Avoid memory leaks -- since n can no longer access the old
	// child, ensure the old child cannot access n.
	q.parent = nil

	n.left = m
	m.parent = n

	if err := Validate(q); err != nil {
		panic(fmt.Errorf("could not set child node: %v", err))
	}
	if err := Validate(n); err != nil {
		panic(fmt.Errorf("could not set child node: %v", err))
	}
	if err := Validate(m); err != nil {
		panic(fmt.Errorf("could not set child node: %v", err))
	}
}
*/

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
		rs := make([]hyperrectangle.R, 0, len(n.data))
		for _, aabb := range n.data {
			rs = append(rs, aabb)
		}
		n.aabbCache = bhr.AABB(rs)
	} else {
		n.aabbCache = bhr.Union(n.Left().AABB(), n.Right().AABB())
	}

	return n.aabbCache
}

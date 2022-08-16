package node

import (
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type N struct {
	id point.ID

	parent *N
	left   *N
	right  *N
	bound hyperrectangle.R
}

func (n *N) Bound() hyperrectangle.R { return n.bound }
func (n *N) Leaf() bool          { return n.left == nil && n.right == nil }

func (n *N) ID() point.ID { return n.id }
func (n *N) L() *N        { return n.left }
func (n *N) R() *N        { return n.right }
func (n *N) Parent() *N   { return n.parent }

func (n *N) Move(id point.ID, offset vector.V) bool { return false }
func (n *N) Remove(id point.ID) bool                { return false }
func (n *N) Insert(id point.ID)                     {}

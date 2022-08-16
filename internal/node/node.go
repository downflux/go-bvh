package node

import (
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type O struct {
	ID     point.ID
	Parent *N
	Left   *N
	Right  *N

	Bound hyperrectangle.R
}

type N struct {
	id     point.ID
	parent *N
	left   *N
	right  *N
	bound  hyperrectangle.R
}

func New(o O) *N {
	return &N{
		id:     o.ID,
		parent: o.Parent,
		left:   o.Left,
		right:  o.Right,
		bound:  o.Bound,
	}
}

func (n *N) Bound() hyperrectangle.R { return n.bound }

func (n *N) Leaf() bool   { return n.Left() == nil && n.Right() == nil }
func (n *N) Left() *N     { return n.left }
func (n *N) Right() *N    { return n.right }
func (n *N) Parent() *N   { return n.parent }
func (n *N) ID() point.ID { return n.id }

func (n *N) Move(id point.ID, offset vector.V) bool { return false }
func (n *N) Remove(id point.ID) bool                { return false }

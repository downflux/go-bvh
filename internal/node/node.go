package node

import (
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type O struct {
	LeafSize int

	IDs    []point.ID
	Parent *N
	Left   *N
	Right  *N
	Bound  hyperrectangle.R
}

type N struct {
	o O
}

func New(o O) *N {
	return &N{
		o: o,
	}
}

func (n *N) Bound() hyperrectangle.R { return n.o.Bound }
func (n *N) Leaf() bool              { return n.Left() == nil && n.Right() == nil }

func (n *N) IDs() []point.ID { return n.o.IDs }
func (n *N) Left() *N        { return n.o.Left }
func (n *N) Right() *N       { return n.o.Right }
func (n *N) Parent() *N      { return n.o.Parent }

func (n *N) Move(id point.ID, offset vector.V) bool { return false }
func (n *N) Remove(id point.ID) bool                { return false }
func (n *N) Insert(id point.ID)                     {}

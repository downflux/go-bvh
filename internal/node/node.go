package node

import (
	"github.com/downflux/go-bvh/internal/allocation"
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type O struct {
	ID point.ID

	Index  allocation.ID
	Parent allocation.ID
	Left   allocation.ID
	Right  allocation.ID

	Bound hyperrectangle.R
}

type N struct {
	id point.ID

	index  allocation.ID
	parent allocation.ID
	left   allocation.ID
	right  allocation.ID

	bound hyperrectangle.R
}

func New(o O) *N {
	return &N{
		id: o.ID,

		index:  o.Index,
		parent: o.Parent,
		left:   o.Left,
		right:  o.Right,

		bound: o.Bound,
	}
}

func (n *N) Bound() hyperrectangle.R { return n.bound }

func (n *N) Index() allocation.ID { return n.index }
func (n *N) ID() point.ID         { return n.id }

func (n *N) Move(id point.ID, offset vector.V) bool { return false }
func (n *N) Remove(id point.ID) bool                { return false }

func Leaf(c allocation.C[*N], n *N) bool { return Left(c, n) == nil && Right(c, n) == nil }
func Left(c allocation.C[*N], n *N) *N   { return c[n.left] }
func Right(c allocation.C[*N], n *N) *N  { return c[n.right] }
func Parent(c allocation.C[*N], n *N) *N { return c[n.parent] }

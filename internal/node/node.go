package node

import (
	"github.com/downflux/go-bvh/internal/allocation"
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-bvh/internal/allocation/id"
)

type O struct {
	ID point.ID

	Index  id.ID
	Parent id.ID
	Left   id.ID
	Right  id.ID

	Bound hyperrectangle.R
}

type N struct {
	id point.ID

	index  id.ID
	parent id.ID
	left   id.ID
	right  id.ID

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

func (n *N) Index() id.ID { return n.index }
func (n *N) ID() point.ID         { return n.id }

func (n *N) SetParent(id id.ID)      { n.parent = id }
func (n *N) SetLeft(id id.ID)        { n.left = id }
func (n *N) SetRight(id id.ID)       { n.right = id }
func (n *N) SetBound(bound hyperrectangle.R) { n.bound = bound }

func Leaf(c allocation.C[*N], n *N) bool { return Left(c, n) == nil && Right(c, n) == nil }
func Left(c allocation.C[*N], n *N) *N   { return c[n.left] }
func Right(c allocation.C[*N], n *N) *N  { return c[n.right] }
func Parent(c allocation.C[*N], n *N) *N { return c[n.parent] }

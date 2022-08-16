package node

import (
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type N struct {
	ID     point.ID
	Parent *N
	Left   *N
	Right  *N
	Bound  hyperrectangle.R
}

func (n *N) Leaf() bool { return n.Left == nil && n.Right == nil }

func (n *N) Move(id point.ID, offset vector.V) bool { return false }
func (n *N) Remove(id point.ID) bool                { return false }

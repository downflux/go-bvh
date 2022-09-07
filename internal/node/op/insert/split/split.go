package split

import (
	"math/rand"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

func midpoint(aabb hyperrectangle.R, axis vector.D) float64 {
	return (aabb.Max().X(axis) - aabb.Min().X(axis)) / 2.0
}

// Execute creates a new node to be inserted into the tree. The returned node is
// a leaf that is guaranteed to have capacity for a new data point.
func Execute(n *node.N) *node.N {
	if n == nil {
		panic("cannot split an empty node")
	}
	if !n.IsLeaf() {
		panic("cannot split a non-leaf node")
	}
	if len(n.Data()) == 0 {
		panic("cannot split a leaf node with no data")
	}

	var pivotIsSet bool
	var axis vector.D
	var pivot float64

	m := node.New(node.O{
		Nodes: n.Cache(),
		Data:  map[id.ID]hyperrectangle.R{},
		Size:  n.Size(),
	})

	for id, aabb := range n.Data() {
		if !pivotIsSet {
			axis = vector.D(rand.Intn(int(aabb.Min().Dimension())))
			pivot = midpoint(aabb, axis)
		}
		// At least one point is kept in n, ensuring then that m is not
		// full.
		if midpoint(aabb, axis) > pivot {
			m.Insert(id, aabb)
		}
	}
	for id := range m.Data() {
		n.Remove(id)
	}
	return m
}

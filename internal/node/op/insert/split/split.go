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

type P func(data map[id.ID]hyperrectangle.R) (map[id.ID]hyperrectangle.R, map[id.ID]hyperrectangle.R)

func partition(data map[id.ID]hyperrectangle.R, axis vector.D, pivot float64) (map[id.ID]hyperrectangle.R, map[id.ID]hyperrectangle.R) {
	a := map[id.ID]hyperrectangle.R{}
	b := map[id.ID]hyperrectangle.R{}

	for id, aabb := range data {
		// At least one point is kept in n, ensuring then that m is not
		// full.
		if midpoint(aabb, axis) > pivot {
			b[id] = aabb
		} else {
			a[id] = aabb
		}
	}
	return a, b
}

func RandomPartition(data map[id.ID]hyperrectangle.R) (map[id.ID]hyperrectangle.R, map[id.ID]hyperrectangle.R) {
	var axis vector.D
	var pivot float64

	for _, aabb := range data {
		axis = vector.D(rand.Intn(int(aabb.Min().Dimension())))
		pivot = midpoint(aabb, axis)
		break
	}

	return partition(data, axis, pivot)
}

// Execute creates a new node to be inserted into the tree. The returned node is
// a leaf that is guaranteed to have capacity for a new data point.
func Execute(n *node.N, p P) *node.N {
	if n == nil {
		panic("cannot split an empty node")
	}
	if !n.IsLeaf() {
		panic("cannot split a non-leaf node")
	}
	if len(n.Data()) == 0 {
		panic("cannot split a leaf node with no data")
	}

	data := n.Data()
	_, b := p(data)
	for id := range b {
		n.Remove(id)
	}

	return node.New(node.O{
		Nodes: n.Cache(),
		Data:  b,
		Size:  n.Size(),
	})
}

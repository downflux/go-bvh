// Package node defines a node interface. This is used by both the cache and
// node implementations to avoid cyclic imports.
package node

import (
	"fmt"
	"math"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache/branch"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

type N interface {
	ID() cid.ID

	IsRoot() bool

	Height() int
	SetHeight(h int)

	IsLeaf() bool
	IsFull() bool
	Leaves() map[id.ID]struct{}

	AABB() hyperrectangle.M

	Heuristic() float64
	SetHeuristic(a float64)

	Child(b branch.B) N
	SetChild(b branch.B, x cid.ID)

	Branch(x cid.ID) branch.B

	Parent() N
	Left() N
	Right() N

	SetParent(x cid.ID)
	SetLeft(x cid.ID)
	SetRight(x cid.ID)
}

// SetHeight will update the height of a node. The input node must have valid
// and up-to-date child nodes.
func SetHeight(n N) {
	if n.IsLeaf() {
		n.SetHeight(0)
	} else {
		n.SetHeight(1 + int(
			math.Max(
				float64(n.Left().Height()),
				float64(n.Right().Height()),
			),
		))
	}
}

// SetAABB updates a node's AABB with the bounding boxes of its children. For a
// leaf node, this bounding box will have a buffer of some given expansion
// factor.
//
// The input node must be valid and up-to-date.
func SetAABB(n N, data map[id.ID]hyperrectangle.R, tolerance float64) {
	if tolerance < 1 {
		panic(fmt.Sprintf("cannot set expansion factor to be less than the AABB size"))
	}

	if !n.IsLeaf() {
		n.AABB().Copy(n.Left().AABB().R())
		n.AABB().Union(n.Right().AABB().R())
		n.SetHeuristic(heuristic.H(n.AABB().R()))
		return
	}

	// TODO(minkezhang): Improve performance by concurrently updating the
	// final AABB.
	var initialized bool
	var k vector.D
	for x := range n.Leaves() {
		if !initialized {
			initialized = true
			n.AABB().Copy(data[x])
			k = data[x].Min().Dimension()
		} else {
			n.AABB().Union(data[x])
		}
	}

	epsilon := math.Pow(tolerance, 1/float64(k))
	for i := vector.D(0); i < k; i++ {
		d := n.AABB().Max().X(i) - n.AABB().Min().X(i)
		offset := (epsilon*d - d) / 2
		n.AABB().Min().SetX(i, n.AABB().Min().X(i)-offset)
		n.AABB().Max().SetX(i, n.AABB().Max().X(i)-offset)
	}
	n.AABB().Scale(epsilon)
	n.SetHeuristic(heuristic.H(n.AABB().R()))
}

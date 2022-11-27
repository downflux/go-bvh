// Package node defines a node interface. This is used by both the cache and
// node implementations to avoid cyclic imports.
package node

import (
	"fmt"
	"math"
	"sync"

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

func Union(data map[id.ID]hyperrectangle.R, xs ...id.ID) hyperrectangle.R {
	k := data[xs[0]].Min().Dimension()
	buf := hyperrectangle.New(
		vector.V(make([]float64, k)),
		vector.V(make([]float64, k)),
	).M()

	if len(xs) < 8 {
		var initialized bool
		for _, x := range xs {
			if !initialized {
				initialized = true
				buf.Copy(data[x])
			} else {
				buf.Union(data[x])
			}
		}
		return buf.R()
	}
	lbuf := hyperrectangle.New(
		vector.V(make([]float64, k)),
		vector.V(make([]float64, k)),
	).M()
	rbuf := hyperrectangle.New(
		vector.V(make([]float64, k)),
		vector.V(make([]float64, k)),
	).M()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		lbuf.Copy(Union(data, xs[:len(xs)/2]...))
	}()
	go func() {
		defer wg.Done()
		rbuf.Copy(Union(data, xs[len(xs)/2+1:]...))
	}()
	wg.Wait()

	buf.Copy(lbuf.R())
	buf.Union(rbuf.R())
	return buf.R()
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

	xs := make([]id.ID, 0, len(n.Leaves()))
	for x := range n.Leaves() {
		xs = append(xs, x)
	}

	n.AABB().Copy(Union(data, xs...))
	k := n.AABB().Min().Dimension()

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

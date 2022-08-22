package rotate

import (
	"github.com/downflux/go-bvh/hyperrectangle"
	"github.com/downflux/go-bvh/internal/allocation"
	"github.com/downflux/go-bvh/internal/allocation/id"
	"github.com/downflux/go-bvh/internal/heuristic"
	"github.com/downflux/go-bvh/internal/node"
)

// Execute finds the configuration of the children and grandchildren which
// minimizes the local heuristic, and rotates the tree accordingly. See Catto
// 2019 for more information. As we are applyiing potential rotations on the
// descendents of the input node, neither the bounds of the parent nor its own
// parent change.
//
// We refer to the local subtree using the same labels as the Catto
// presentation, i.e.
//
//	     A
//	    / \
//	   /   \
//	  B     C
//	 / \   / \
//	D   E F   G
//
// We will attempt to apply a rotation on the node pairs
//
//	BF, BG, CD, and CE
//
// and choose the configuration with the minimal "surface area" heuristic.
//
// N.B.: The Box2D implementation uses pure tree heights as the balancing
// heuristic.
//
// TODO(minkezhang): Investigate relative performance characteristics of
// choosing between these two balancing modes.
func Execute(nodes allocation.C[*node.N], aid id.ID) {
	var a, b, c, d, e, f, g *node.N

	a = nodes[aid]
	if node.Leaf(nodes, a) {
		return
	}

	b = node.Left(nodes, a)
	c = node.Right(nodes, a)

	if b != nil {
		d = node.Left(nodes, b)
		e = node.Right(nodes, b)
	}
	if c != nil {
		f = node.Left(nodes, c)
		g = node.Right(nodes, c)
	}

	// We are using the BF rotation syntax as explored in the Catto
	// presentation. Here, we want to check if the b node should be swapped
	// with the f node (i.e. from b to f). The c node is therefore the
	// sibling of the b node, and g is the sibling of the f node (i.e. the
	// other child of c).
	type r struct {
		b *node.N

		c *node.N
		// Note that if f is non-nil, then both the parent c and sibling
		// g are also non-nil.
		f *node.N
		g *node.N
	}
	rs := []r{
		{b: b, c: c, f: f, g: g},
		{b: b, c: c, f: g, g: f},
		{b: c, c: b, f: e, g: d},
		{b: c, c: b, f: d, g: e},
	}

	var optimal r
	h := heuristic.Heuristic(b.Bound()) + heuristic.Heuristic(c.Bound())
	for _, rotation := range rs {
		if rotation.b != nil && rotation.f != nil {
			swapped := heuristic.Heuristic(rotation.f.Bound()) + heuristic.Heuristic(
				hyperrectangle.Union(
					rotation.b.Bound(),
					rotation.g.Bound(),
				),
			)
			if swapped < h {
				optimal = rotation
				h = swapped
			}
		}
	}

	// Apply a rotation if there is a configuration which has a lower
	// surface area heuristic.
	if optimal.b != nil && optimal.f != nil {
		if node.Left(nodes, a) == optimal.b {
			a.SetLeft(optimal.f.Index())
		} else {
			a.SetRight(optimal.f.Index())
		}

		if node.Left(nodes, optimal.c) == optimal.f {
			optimal.c.SetLeft(optimal.b.Index())
		} else {
			optimal.c.SetRight(optimal.b.Index())
		}

		optimal.b.SetParent(optimal.c.Index())
		optimal.f.SetParent(a.Index())
		optimal.c.SetBound(hyperrectangle.Union(optimal.b.Bound(), optimal.g.Bound()))
	}

	if p := node.Parent(nodes, a); p != nil {
		Execute(nodes, p.Index())
	}
}

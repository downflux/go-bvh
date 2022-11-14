package bvh

import (
	"math"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type rotationType int

const (
	rotationTypeNull rotationType = iota
	rotationTypeBF
	rotationTypeDF
)

// rotate returns a valid rotation which will not unbalance our BVH tree.
//
// The list of rotations to be considered is taken from Kopta et al. 2012.
//
//	     A
//	    / \
//	   /   \
//	  B     C
//	 / \   / \
//	D   E F   G
//
// Here, we want to create an optimal rotation for the A subtree. Per Kopta et
// al., we will consider the following node swaps --
//
// B -> F
// B -> G
// C -> D
// C -> E
//
// These rotations are also the same ones mentioned in the Catto 2019 slides.
//
// Kopta also considers the following rotations --
//
// D -> F
// D -> G
//
// Because the BVH tree is invariant under reflection, we do not consider the
// redundant swaps E -> F and E -> G.
//
// For each of these swaps, we expect the balance of A to potentially change
// (e.g. if we are lifting too shallow a tree), and the total heuristic
//
//	H(A.L) + H(A.R)
//
// to change. WLOG we use the A.L notation here instead of B, as B may have been
// swapped.
func rotate(a node.N, data map[id.ID]hyperrectangle.R, epsilon float64) node.N {
	if a.Height() < 2 {
		return a
	}

	b := a.Left()
	c := a.Right()

	r := &r{}

	// Cache the heuristic values for faster computes.
	bh := heuristic.H(b.AABB().R())
	ch := heuristic.H(c.AABB().R())

	opt := bh + ch

	k := a.AABB().Min().Dimension()
	// Generate some virtual node buffers.
	bbuf := hyperrectangle.New(
		vector.V(make([]float64, k)),
		vector.V(make([]float64, k)),
	).M()
	cbuf := hyperrectangle.New(
		vector.V(make([]float64, k)),
		vector.V(make([]float64, k)),
	).M()

	var d, e, f, g node.N
	if !c.IsLeaf() {
		f = c.Left()
		g = c.Right()

		if h, ok := checkBF(b, f, g, opt, bbuf); ok {
			opt = h
			r.src = b
			r.dest = f
			r.t = rotationTypeBF
		}

		if h, ok := checkBF(b, g, f, opt, bbuf); ok {
			opt = h
			r.src = b
			r.dest = g
			r.t = rotationTypeBF
		}
	}

	if !b.IsLeaf() {
		d = b.Left()
		e = b.Right()

		if h, ok := checkBF(c, d, e, opt, cbuf); ok {
			opt = h
			r.src = c
			r.dest = d
			r.t = rotationTypeBF
		}

		if h, ok := checkBF(c, e, d, opt, cbuf); ok {
			opt = h
			r.src = c
			r.dest = e
			r.t = rotationTypeBF
		}
	}

	if !b.IsLeaf() && !c.IsLeaf() {
		if h, ok := checkDF(d, e, f, g, opt, bbuf, cbuf); ok {
			opt = h
			r.src = d
			r.dest = f
			r.t = rotationTypeDF

		}
		if h, ok := checkDF(d, e, g, f, opt, bbuf, cbuf); ok {
			opt = h
			r.src = d
			r.dest = g
			r.t = rotationTypeDF
		}
	}

	switch r.t {
	case rotationTypeBF:
		swap(r.src, r.dest)

		node.SetAABB(r.src.Parent(), data, epsilon)
		node.SetHeight(r.src.Parent())

		node.SetAABB(a, data, epsilon)
		node.SetHeight(a)
	case rotationTypeDF:
		swap(r.src, r.dest)

		node.SetAABB(r.src.Parent(), data, epsilon)
		node.SetHeight(r.src.Parent())

		node.SetAABB(r.dest.Parent(), data, epsilon)
		node.SetHeight(r.dest.Parent())

		node.SetAABB(a, data, epsilon)
		node.SetHeight(a)
	}

	return a
}

type r struct {
	t    rotationType
	src  node.N
	dest node.N
}

// checkBF checks if a potential B -> F rotation will generate a more efficient
// local subtree than the existing configuration.
//
// checkBF will return the new lower bound heuristic and a true value if the
// input configuration generates a better configuration.
//
//	  A
//	 / \
//	B   C
//	   / \
//	  F   G
func checkBF(b node.N, f node.N, g node.N, opt float64, bbuf hyperrectangle.M) (float64, bool) {
	// Simulate a B -> F rotation, which mutates the c node.
	bbuf.Copy(b.AABB().R())
	bbuf.Union(g.AABB().R())
	height := math.Max(float64(b.Height()), float64(g.Height())) + 1
	balanced := math.Abs(height-float64(f.Height())) <= 1

	h := heuristic.H(f.AABB().R()) + heuristic.H(bbuf.R())
	if h < opt && balanced {
		return h, true
	}
	return opt, false
}

func checkDF(d node.N, e node.N, f node.N, g node.N, opt float64, bbuf hyperrectangle.M, cbuf hyperrectangle.M) (float64, bool) {
	bbuf.Copy(f.AABB().R())
	bbuf.Union(e.AABB().R())
	bheight := math.Max(float64(f.Height()), float64(e.Height())) + 1

	cbuf.Copy(d.AABB().R())
	cbuf.Union(g.AABB().R())
	cheight := math.Max(float64(d.Height()), float64(g.Height())) + 1
	balanced := math.Abs(bheight-cheight) <= 1

	h := heuristic.H(bbuf.R()) + heuristic.H(cbuf.R())
	if h < opt && balanced {
		return h, true
	}
	return opt, false
}

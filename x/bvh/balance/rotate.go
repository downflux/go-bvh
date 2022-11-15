package balance

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

// r is a mutable rotation struct which keeps track of the current optimal local
// subtree rotation.
type r struct {
	t    rotationType
	src  node.N
	dest node.N
}

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
// to change.
func rotate(x node.N, data map[id.ID]hyperrectangle.R, epsilon float64) node.N {
	if x.Height() < 2 {
		return x
	}

	r := optimal(x)

	switch r.t {
	case rotationTypeBF:
		swap(r.src, r.dest)

		// By construction, r.src is the b node, and therefore the
		// shallower node before the swap. Now that the nodes are
		// swapped, r.src is the deeper node.
		node.SetAABB(r.src.Parent(), data, epsilon)
		node.SetHeight(r.src.Parent())

		node.SetAABB(x, data, epsilon)
		node.SetHeight(x)
	case rotationTypeDF:
		swap(r.src, r.dest)

		// A D -> F rotation means both nodes are of the same depth --
		// both their parents need to be updated in addition to the
		// local root node.
		node.SetAABB(r.src.Parent(), data, epsilon)
		node.SetHeight(r.src.Parent())

		node.SetAABB(r.dest.Parent(), data, epsilon)
		node.SetHeight(r.dest.Parent())

		node.SetAABB(x, data, epsilon)
		node.SetHeight(x)
	}

	return x
}

func optimal(a node.N) *r {
	r := &r{}

	if a.Height() < 2 {
		return r
	}

	b := a.Left()
	c := a.Right()

	opt := heuristic.H(b.AABB().R()) + heuristic.H(c.AABB().R())

	// Generate a virtual node buffer.
	k := a.AABB().Min().Dimension()
	buf := hyperrectangle.New(
		vector.V(make([]float64, k)),
		vector.V(make([]float64, k)),
	).M()

	var d, e, f, g node.N
	if !c.IsLeaf() {
		f = c.Left()
		g = c.Right()

		if h, ok := checkBF(b, f, g, opt, buf); ok {
			opt = h
			r.src = b
			r.dest = f
			r.t = rotationTypeBF
		}

		if h, ok := checkBF(b, g, f, opt, buf); ok {
			opt = h
			r.src = b
			r.dest = g
			r.t = rotationTypeBF
		}
	}

	if !b.IsLeaf() {
		d = b.Left()
		e = b.Right()

		if h, ok := checkBF(c, d, e, opt, buf); ok {
			opt = h
			r.src = c
			r.dest = d
			r.t = rotationTypeBF
		}

		if h, ok := checkBF(c, e, d, opt, buf); ok {
			opt = h
			r.src = c
			r.dest = e
			r.t = rotationTypeBF
		}
	}

	if !b.IsLeaf() && !c.IsLeaf() {
		if h, ok := checkDF(d, e, f, g, opt, buf); ok {
			opt = h
			r.src = d
			r.dest = f
			r.t = rotationTypeDF

		}
		if h, ok := checkDF(d, e, g, f, opt, buf); ok {
			opt = h
			r.src = d
			r.dest = g
			r.t = rotationTypeDF
		}
	}

	return r
}

// merge simulates the results if the input nodes are set as siblings of one
// another.
func merge(l node.N, r node.N, buf hyperrectangle.M) (int, bool, float64) {
	buf.Copy(l.AABB().R())
	buf.Union(r.AABB().R())

	height := int(math.Abs(float64(l.Height() - r.Height())))
	return height, height <= 1, heuristic.H(buf.R())
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
func checkBF(b node.N, f node.N, g node.N, opt float64, buf hyperrectangle.M) (float64, bool) {
	if _, balanced, h := merge(b, g, buf); balanced && h < opt {
		return opt, false
	}
	return opt, false
}

func checkDF(d node.N, e node.N, f node.N, g node.N, opt float64, buf hyperrectangle.M) (float64, bool) {
	bheight, bbalanced, bh := merge(f, e, buf)
	cheight, cbalanced, ch := merge(d, g, buf)
	balanced := math.Abs(float64(bheight-cheight)) <= 1

	h := bh + ch

	if bbalanced && cbalanced && balanced && h < opt {
		return h, true
	}
	return opt, false
}

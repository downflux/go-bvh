package balance

import (
	"math"

	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/op/unsafe"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type rtype int

const (
	rtypeUnknown rtype = iota
	rtypeBF
	rtypeDF
)

// Rotate returns a valid rotation which will not unbalance our BVH tree.
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
func Rotate(x node.N) node.N {
	if x.Height() < 2 {
		return x
	}

	var b, c, d, e, f, g node.N
	b, c = x.Left(), x.Right()
	if !b.IsLeaf() {
		d, e = b.Left(), b.Right()
	}
	if !c.IsLeaf() {
		f, g = c.Left(), c.Right()
	}

	buf := hyperrectangle.New(
		vector.V(make([]float64, x.AABB().Min().Dimension())),
		vector.V(make([]float64, x.AABB().Min().Dimension())),
	).M()

	h := b.Heuristic() + c.Heuristic()
	r := &struct {
		source node.N
		target node.N

		rtype rtype
	}{}

	if !c.IsLeaf() {
		if ok, i := mergeBF(b, f, g, buf); ok && i < h {
			h = i
			r.source = b
			r.target = f
			r.rtype = rtypeBF
		}
		if ok, i := mergeBF(b, g, f, buf); ok && i < h {
			h = i
			r.source = b
			r.target = g
			r.rtype = rtypeBF
		}
	}
	if !b.IsLeaf() {
		if ok, i := mergeBF(c, d, e, buf); ok && i < h {
			h = i
			r.source = c
			r.target = d
			r.rtype = rtypeBF
		}
		if ok, i := mergeBF(c, e, d, buf); ok && i < h {
			h = i
			r.source = c
			r.target = e
			r.rtype = rtypeBF
		}
	}
	if !b.IsLeaf() && !c.IsLeaf() {
		if ok, i := mergeDF(d, e, f, g, buf); ok && i < h {
			h = i
			r.source = d
			r.target = f
			r.rtype = rtypeDF
		}
		if ok, i := mergeDF(d, e, g, f, buf); ok && i < h {
			h = i
			r.source = d
			r.target = g
			r.rtype = rtypeDF
		}
	}

	switch r.rtype {
	case rtypeBF:
		unsafe.Swap(r.source, r.target)

		node.SetAABB(r.source.Parent(), nil, 1)
		node.SetHeight(r.source.Parent())

		// The AABB of the local root has not changed, so skip
		// re-calculating the bounding box.
		node.SetHeight(x)
	case rtypeDF:
		unsafe.Swap(r.source, r.target)

		node.SetAABB(r.source.Parent(), nil, 1)
		node.SetAABB(r.target.Parent(), nil, 1)
		node.SetHeight(r.source.Parent())
		node.SetHeight(r.target.Parent())

		node.SetHeight(x)
	}

	return x
}

// mergeBF checks the cost due to merging the B node with the F node (i.e. a
// descendent of the C node). It returns true if the configuration preserves
// tree balance, and the calculated heuristic cost of the two child nodes (i.e.
// F and the pseudonode consisting of B and G).
func mergeBF(b node.N, f node.N, g node.N, buf hyperrectangle.M) (bool, float64) {
	lheight := f.Height()
	rheight := g.Height()
	if rheight < b.Height() {
		rheight = b.Height()
	}
	rheight += 1

	if dh := lheight - rheight; !(dh > -2 && dh < 2) {
		return false, math.Inf(1)
	}

	buf.Copy(b.AABB().R())
	buf.Union(g.AABB().R())

	h := f.Heuristic() + heuristic.H(buf.R())
	return true, h
}

func mergeDF(d node.N, e node.N, f node.N, g node.N, buf hyperrectangle.M) (bool, float64) {
	lheight := e.Height()
	if lheight < f.Height() {
		lheight = f.Height()
	}
	lheight += 1
	rheight := g.Height()
	if rheight < d.Height() {
		rheight = d.Height()
	}
	rheight += 1

	if dh := lheight - rheight; !(dh > -2 && dh < 2) {
		return false, math.Inf(1)
	}

	buf.Copy(e.AABB().R())
	buf.Union(f.AABB().R())

	h := heuristic.H(buf.R())

	buf.Copy(d.AABB().R())
	buf.Union(g.AABB().R())

	h += heuristic.H(buf.R())

	return true, h
}

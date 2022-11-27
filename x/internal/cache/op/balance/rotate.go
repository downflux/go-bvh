package balance

import (
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/op/balance/pseudonode"
	"github.com/downflux/go-bvh/x/internal/cache/op/unsafe"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type rotation struct {
	b pseudonode.I
	c pseudonode.I

	source node.N
	target node.N
}

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
	if x.Height() < 1 {
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

	rotations := []rotation{}
	if !b.IsLeaf() {
		rotations = append(rotations, rotation{
			b:      pseudonode.New(c, e, buf),
			c:      d,
			source: c,
			target: d,
		}, rotation{
			b:      pseudonode.New(d, c, buf),
			c:      e,
			source: c,
			target: e,
		})
	}
	if !c.IsLeaf() {
		rotations = append(rotations, rotation{
			b:      f,
			c:      pseudonode.New(b, g, buf),
			source: b,
			target: f,
		}, rotation{
			b:      g,
			c:      pseudonode.New(b, f, buf),
			source: b,
			target: g,
		})
	}
	if !b.IsLeaf() && !c.IsLeaf() {
		rotations = append(rotations, rotation{
			b:      pseudonode.New(f, e, buf),
			c:      pseudonode.New(d, g, buf),
			source: d,
			target: f,
		}, rotation{
			b:      pseudonode.New(g, e, buf),
			c:      pseudonode.New(f, d, buf),
			source: d,
			target: g,
		})
	}

	h := b.Heuristic() + c.Heuristic()
	opt := rotation{}

	for _, r := range rotations {
		if g := r.b.Heuristic() + r.c.Heuristic(); balanced(r.b, r.c) && g < h {
			opt = r
			h = g
		}
	}

	if opt != (rotation{}) {
		unsafe.Swap(opt.source, opt.target)

		if x.ID() != opt.source.Parent().ID() {
			node.SetAABB(opt.source.Parent(), nil, 1)
			node.SetHeight(opt.source.Parent())
		}

		if x.ID() != opt.target.Parent().ID() {
			node.SetAABB(opt.target.Parent(), nil, 1)
			node.SetHeight(opt.target.Parent())
		}
	}

	return x
}

func balanced(a, b pseudonode.I) bool {
	return a.Height()-b.Height() < 2 && a.Height()-b.Height() > -2
}

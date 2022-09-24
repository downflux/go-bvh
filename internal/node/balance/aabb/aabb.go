// Package aabb returns a valid rotation which will not unbalance our BVH tree.
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
// FOr each of these swaps, we expect the balance of A to potentially change
// (e.g. if we are lifting too shallow a tree), and the total heuristic
//
//	H(A.L) + H(A.R)
//
// to change. WLOG we use the A.L notation here instead of B, as B may have been
// swapped.
package aabb

import (
	"math"

	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
)

var (
	_ ni = &node.N{}
	_ ni = p{}
)

type ni interface {
	AABB() hyperrectangle.R
	Height() uint
}

func balanced(a, b ni) bool { return math.Abs(float64(a.Height())-float64(b.Height())) < 2 }

// p is a pseudo-node.
type p struct {
	l, r *node.N
}

func (p p) AABB() hyperrectangle.R { return bhr.Union(p.l.AABB(), p.r.AABB()) }
func (p p) Height() uint           { return uint(math.Max(float64(p.l.Height()), float64(p.r.Height()))) + 1 }

type candidate struct {
	b, c        ni
	src, target *node.N
}

func (c candidate) Check(h float64) bool {
	return balanced(b, c)
}

type R struct {
	Src    *node.N
	Target *node.N
}

type t struct {
	b, c, d, e, f, g *node.N

	// bh is the height of the b node. We pre-calculate these for convenience.
	bh, ch, dh, eh, fh, gh float64
}

// generate returns a set of valid swaps on a subtree which fulfills the AVL
// balance guarantee.
func generate(n *node.N) []R {
	if n.IsLeaf() {
		return nil
	}

	subtree := &t{
		b: n.Left(),
		c: n.Right(),

		bh: float64(n.Left().Height()),
		ch: float64(n.Right().Height()),
	}
	if !n.Left().IsLeaf() {
		subtree.d = n.Left().Left()
		subtree.e = n.Left().Right()

		subtree.dh = float64(subtree.d.Height())
		subtree.eh = float64(subtree.e.Height())
	}
	if !n.Right().IsLeaf() {
		subtree.f = n.Right().Left()
		subtree.g = n.Right().Right()

		subtree.fh = float64(subtree.f.Height())
		subtree.gh = float64(subtree.g.Height())
	}

	rs := make([]R, 0, 8)

	if !subtree.c.IsLeaf() {
		// Simulate the B -> F rotation. In this case, C's  height (and
		// AABB) may change, as its new children will be B and G. We
		// need to ensure A is still balanced to fulfill the AVL
		// guarantees -- that is,
		//
		//   |Height(A.L) - Height(A.R)| < 2
		if l, r := subtree.fh, math.Max(subtree.bh, subtree.gh)+1; math.Abs(l-r) < 2 {
			rs = append(rs, R{
				Src:    subtree.b,
				Target: subtree.f,
			})
		}
		if l, r := subtree.gh, math.Max(subtree.bh, subtree.fh)+1; math.Abs(l-r) < 2 {
			rs = append(rs, R{
				Src:    subtree.b,
				Target: subtree.g,
			})
		}
	}

	if !subtree.b.IsLeaf() {
		if l, r := math.Max(subtree.ch, subtree.eh)+1, subtree.dh; math.Abs(l-r) < 2 {
			rs = append(rs, R{
				Src:    subtree.c,
				Target: subtree.d,
			})
		}
		if l, r := math.Max(subtree.ch, subtree.dh)+1, subtree.eh; math.Abs(l-r) < 2 {
			rs = append(rs, R{
				Src:    subtree.c,
				Target: subtree.e,
			})
		}
	}

	if !subtree.b.IsLeaf() && !subtree.c.IsLeaf() {
		// By similar logic as above, simulate the potential balance
		// changes to B and C if we were to swap the leaves directly.
		//
		// For a D -> F swap, the children of B will be F and E, and the
		// children of C will be D and G.
		if l, r := math.Max(subtree.fh, subtree.eh), math.Max(subtree.dh, subtree.gh); math.Abs(l-r) < 2 {
			rs = append(rs, R{
				Src:    subtree.d,
				Target: subtree.f,
			})
		}
		if l, r := math.Max(subtree.gh, subtree.eh), math.Max(subtree.fh, subtree.dh); math.Abs(l-r) < 2 {
			rs = append(rs, R{
				Src:    subtree.d,
				Target: subtree.g,
			})
		}
	}

	return rs
}

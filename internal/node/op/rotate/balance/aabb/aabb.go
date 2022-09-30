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
	"github.com/downflux/go-bvh/internal/heuristic"
	"github.com/downflux/go-bvh/internal/node"
	"github.com/downflux/go-bvh/internal/node/op/rotate/balance/aabb/candidate"
)

type R struct {
	Src    *node.N
	Target *node.N
}

type t struct {
	b, c, d, e, f, g *node.N
}

// generate returns a set of candidate swaps. These candidates may both
// unbalance and decrease the quality of the BVH.
func generate(n *node.N) []candidate.C {
	if n.IsLeaf() {
		return nil
	}

	subtree := &t{
		b: n.Left(),
		c: n.Right(),
	}
	if !n.Left().IsLeaf() {
		subtree.d = n.Left().Left()
		subtree.e = n.Left().Right()
	}
	if !n.Right().IsLeaf() {
		subtree.f = n.Right().Left()
		subtree.g = n.Right().Right()
	}

	candidates := make([]candidate.C, 0, 8)

	if !subtree.c.IsLeaf() {
		candidates = append(candidates, candidate.C{
			B: subtree.f,
			C: candidate.P{L: subtree.b, R: subtree.g},

			Src:    subtree.b,
			Target: subtree.f,
		}, candidate.C{
			B: subtree.g,
			C: candidate.P{L: subtree.b, R: subtree.f},

			Src:    subtree.b,
			Target: subtree.g,
		})
	}
	if !subtree.b.IsLeaf() {
		candidates = append(candidates, candidate.C{
			B: candidate.P{L: subtree.c, R: subtree.e},
			C: subtree.d,

			Src:    subtree.c,
			Target: subtree.d,
		}, candidate.C{
			B: candidate.P{L: subtree.d, R: subtree.c},
			C: subtree.e,

			Src:    subtree.c,
			Target: subtree.e,
		})
	}

	if !subtree.b.IsLeaf() && !subtree.c.IsLeaf() {
		candidates = append(candidates, candidate.C{
			B: candidate.P{L: subtree.f, R: subtree.e},
			C: candidate.P{L: subtree.d, R: subtree.g},

			Src:    subtree.d,
			Target: subtree.f,
		}, candidate.C{
			B: candidate.P{L: subtree.g, R: subtree.e},
			C: candidate.P{L: subtree.f, R: subtree.d},

			Src:    subtree.d,
			Target: subtree.g,
		})
	}

	return candidates
}

func Query(n *node.N) R {
	if n == nil {
		panic("cannot query an empty node")
	}

	if n.IsLeaf() {
		return R{}
	}

	var rotation R
	h := heuristic.H(n.Left().AABB()) + heuristic.H(n.Right().AABB())

	candidates := generate(n)
	for _, c := range candidates {
		if g := heuristic.H(c.B.AABB()) + heuristic.H(c.C.AABB()); c.Balanced() && g < h {
			rotation = R{
				Src:    c.Src,
				Target: c.Target,
			}
			h = g
		}
	}

	return rotation
}

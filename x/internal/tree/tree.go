package tree

import (
	"fmt"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	bhr "github.com/downflux/go-bvh/hyperrectangle"
)

type S struct {
	N *node.N

	// B is the branch from the parent which will result in accessing this
	// node. This value is set to BranchInvalid for the root node.
	B node.Branch
}

type T struct {
	nodes *cache.C[*node.N]
	root  *node.N

	dataLookup map[id.ID]hyperrectangle.R
	leafLookup map[cache.ID][]id.ID

	aabbCache        map[cache.ID]hyperrectangle.R
	aabbCacheIsValid map[cache.ID]bool

	heightCache        map[cache.ID]int
	heightCacheIsValid map[cache.ID]bool

	size int
	k    vector.D
}

// Path returns a path to the specified node. We assume that the input node is
// valid -- that is, it exists in the tree cache and its root matches the tree
// root.
func (t *T) Path(n *node.N) []S {
	s := make([]S, 0, 16)

	for p := n; !p.IsRoot(t.nodes); p = n.Parent(t.nodes) {
		s = append(s, S{
			N: n,
			B: p.Branch(n.ID()),
		})
		n = p
	}
	s = append(s, S{
		N: t.root,
		B: node.BranchInvalid,
	})

	// Ensure root node is the first element.
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}

// Height returns the subtree height. We assume the input node is valid.
//
// Leaf nodes have a height of 0.
func (t *T) Height(n *node.N) int {
	x := n.ID()

	if t.heightCacheIsValid[x] {
		return t.heightCache[x]
	}

	t.heightCacheIsValid[x] = true

	if n.IsLeaf(t.nodes) {
		t.heightCache[x] = 0
	} else {
		h := t.Height(n.Left(t.nodes))
		if g := t.Height(n.Right(t.nodes)); g > h {
			h = g
		}
		t.heightCache[x] = 1 + g
	}

	return t.heightCache[x]
}

// AABB returns the bounding box of the subtree. We assume the input node is
// valid.
func (t *T) AABB(n *node.N) hyperrectangle.R {
	x := n.ID()

	if t.aabbCacheIsValid[x] {
		return t.aabbCache[x]
	}

	t.aabbCacheIsValid[x] = true

	if n.IsLeaf(t.nodes) {
		if len(t.leafLookup[x]) == 0 {
			panic(fmt.Sprintf("AABB is not defined for an empty leaf node %v", x))
		}

		rs := make([]hyperrectangle.R, 0, len(t.leafLookup[x]))
		for _, y := range t.leafLookup[x] {
			rs = append(rs, t.dataLookup[y])
		}

		bhr.AABBBuf(rs, t.aabbCache[x])
	} else {
		bhr.UnionBuf(
			t.AABB(n.Left(t.nodes)),
			t.AABB(n.Right(t.nodes)),
			t.aabbCache[x],
		)
	}

	return t.aabbCache[x]
}

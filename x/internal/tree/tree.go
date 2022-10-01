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

type AABBCache struct {
	cache        map[cache.ID]hyperrectangle.R
	cacheIsValid map[cache.ID]bool
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

func (t *T) Node(x cache.ID) *node.N {
	n, ok := t.nodes.Get(x)
	if !ok {
		panic(fmt.Sprintf("cannot get non-existent node %v", x))
	}

	return n
}

func (t *T) Insert(x id.ID, aabb hyperrectangle.R) {
	if _, ok := t.dataLookup[x]; ok {
		panic(fmt.Sprintf("cannot insert duplicate node %v", x))
	}

	t.dataLookup[x] = aabb

	if t.root == nil {
		m := node.New(t.nodes, node.O{
			Parent: cache.IDInvalid,
			Left:   cache.IDInvalid,
			Right:  cache.IDInvalid,
		})
		t.leafLookup[m.ID()] = append(t.leafLookup[m.ID()], x)
	}

	panic("unimplemented")
}

func (t *T) Remove(x id.ID) {
	if _, ok := t.dataLookup[x]; !ok {
		panic(fmt.Sprintf("cannot remove non-existent node %v", x))
	}

	panic("unimplemented")
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

	if n.IsLeaf() {
		t.heightCache[x] = 0
	} else {
		h := t.Height(n.Left(t.nodes))
		if g := t.Height(n.Right(t.nodes)); g > h {
			h = g
		}
		t.heightCache[x] = 1 + h
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

	if n.IsLeaf() {
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

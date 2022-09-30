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

func (t *T) Height(n *node.N) int {
	x := n.ID()

	if t.heightCacheIsValid[x] {
		return t.heightCache[x]
	}

	t.heightCacheIsValid[x] = true

	if n.IsLeaf(t.nodes) {
		t.heightCache[x] = 0
	} else {
		t.heightCache[x] = 1 + t.Height(n.Left(t.nodes)) + t.Height(n.Right(t.nodes))
	}

	return t.heightCache[x]
}

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

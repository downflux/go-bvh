package tree

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type T struct {
	nodes  *cache.C[*node.N]
	leaves map[cache.ID]map[id.ID]hyperrectangle.R
	root   *node.N

	size int
	k    vector.D
}

func (t *T) Insert(x id.ID, aabb hyperrectangle.R) {
	if t.root == nil {
		t.root = node.New(node.O{
			Nodes:  t.nodes,
			Leaves: t.leaves,

			Parent: cache.IDInvalid,
			Left:   cache.IDInvalid,
			Right:  cache.IDInvalid,
		})

		t.leaves[t.root.ID()] = map[id.ID]hyperrectangle.R{
			x: aabb,
		}

		return
	}

	panic("unimplemented")
}

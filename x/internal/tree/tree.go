package tree

import (
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type T struct {
	cache *cache.C[*node.N]
	root  *node.N

	data map[cache.ID]map[id.ID]hyperrectangle.R

	size int
	k    vector.D
}

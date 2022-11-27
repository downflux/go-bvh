package metrics

import (
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/node/util"
)

func LeafSize(n node.N) float64 {
	var leaves int
	var objects int

	util.PreOrder(n, func(n node.N) {
		if n.IsLeaf() {
			leaves += 1
			objects += len(n.Leaves())
		}
	})
	return float64(objects) / float64(leaves)
}

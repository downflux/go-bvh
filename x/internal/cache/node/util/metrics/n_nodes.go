package metrics

import (
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/node/util"
)

func NNodes(n node.N) int {
	nnodes := 0
	util.PreOrder(n, func(n node.N) { nnodes += 1 })
	return nnodes
}

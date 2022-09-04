package refit

import (
	"github.com/downflux/go-bvh/internal/node"
)

func Execute(n *node.N) *node.N {
	if n == nil {
		panic("cannot refit a nil node")
	}

	n.InvalidateAABBCache()
	return n.Root()
}

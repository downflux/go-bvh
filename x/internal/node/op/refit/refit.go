package refit

import (
	"github.com/downflux/go-bvh/x/internal/node"
)

func Execute(n *node.N) *node.N {
	if n == nil {
		panic("cannot refit a nil node")
	}

	n.InvalidateAABBCache()
	if n.IsRoot() {
		return n
	}
	return Execute(n.Parent())
}

package refit

import (
	"github.com/downflux/go-bvh/x/internal/node"
)

func Execute[T comparable](n *node.N[T]) *node.N[T] {
	n.InvalidateAABBCache()
	if n.Parent() == nil {
		return n
	}
	return Execute(n.Parent())
}

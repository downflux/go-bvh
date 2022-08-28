package rotate

import (
	"github.com/downflux/go-bvh/x/internal/node"
	"github.com/downflux/go-bvh/x/internal/node/op/rotate/rotation"
)

func Execute(n *node.N) *node.N {
	rotation.Generate(n)
	return nil
}

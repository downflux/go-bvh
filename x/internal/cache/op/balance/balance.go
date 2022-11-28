package balance

import (
	"github.com/downflux/go-bvh/x/internal/cache/node"
)

// B will update the subtree of the input node n. This node will have a valid
// cache at the end of this operation.
type B func(n node.N) node.N

package balance

import (
	"github.com/downflux/go-bvh/x/internal/cache/node"
)

// B will update the subtree of the input node n. This node will have an invalid
// cache at the end of this operation. The user must manually update the cache.
// Child nodes of n have valid caches.
type B func(n node.N) node.N

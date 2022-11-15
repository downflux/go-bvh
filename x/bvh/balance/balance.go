package balance

import (
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

func B(n node.N, data map[id.ID]hyperrectangle.R, epsilon float64) node.N {
	return rotate(avl(n, data, epsilon), data, epsilon)
}

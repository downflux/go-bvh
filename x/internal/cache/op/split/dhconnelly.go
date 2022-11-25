package split

import (
	"math"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

// DHConnelly implements the leaf node splitting function as used in
// github.com/dhconnelly/rtreego.
func DHConnelly(c *cache.C, data map[id.ID]hyperrectangle.R, n node.N, m node.N) {
	if c.LeafSize() == 1 {
		for x := range n.Leaves() {
			m.Leaves()[x] = struct{}{}
			delete(n.Leaves(), x)
			break
		}
		return
	}

	buf := hyperrectangle.New(
		vector.V(make([]float64, c.K())),
		vector.V(make([]float64, c.K())),
	).M()

	leaves := make([]id.ID, 0, len(n.Leaves()))
	for x := range n.Leaves() {
		leaves = append(leaves, x)
	}

	li, ri := seeds(data, leaves, buf)

	remaining := append(leaves[:li], leaves[li+1:ri]...)
	remaining = append(remaining, leaves[ri+1:]...)

	for len(remaining) > 0 {
		ni := next(data, remaining, buf)
	}

	return
}

func next(data map[id.ID]hyperrectangle.R, leaves []id.ID, buf hyperrectangle.M) int {
	return -1
}

func seeds(data map[id.ID]hyperrectangle.R, leaves []id.ID, buf hyperrectangle.M) (int, int) {
	var l, r int
	h := math.Inf(-1)

	for i, x := range leaves {
		for j, y := range leaves[i+1:] {
			buf.Copy(data[x])
			buf.Union(data[y])

			if g := heuristic.H(buf.R()) - (heuristic.H(data[x]) + heuristic.H(data[y])); g > h {
				h = g
				l = i
				r = j
			}
		}
	}

	return l, r
}

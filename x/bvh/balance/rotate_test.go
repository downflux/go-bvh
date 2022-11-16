package balance

import (
	"testing"

	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-geometry/epsilon"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

func TestMerge(t *testing.T) {
	const k = 2

	type w struct {
		height   int
		balanced bool
		h        float64
	}

	type config struct {
		name string
		l    node.N
		r    node.N
		want w
	}

	configs := []config{}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			buf := hyperrectangle.New(
				vector.V(make([]float64, 2)),
				vector.V(make([]float64, 2)),
			).M()
			got := &w{}

			got.height, got.balanced, got.h = merge(c.l, c.r, buf)
			if got.height != c.want.height {
				t.Errorf("height = %v, c.want = %v", got.height, c.want.height)
			}
			if got.balanced != c.want.balanced {
				t.Errorf("balanced = %v, c.want = %v", got.balanced, c.want.balanced)
			}
			if !epsilon.Within(got.h, c.want.h) {
				t.Errorf("h = %v, c.want = %v", got.h, c.want.h)
			}
		})
	}
}

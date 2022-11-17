package balance

import (
	"testing"

	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/epsilon"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

func TestCheckBF(t *testing.T) {
	const k = 2
	type w struct {
		h       float64
		optimal bool
	}

	type config struct {
		name string
		b    node.N
		f    node.N
		g    node.N
		opt  float64
		want w
	}

	configs := []config{}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			buf := hyperrectangle.New(
				vector.V(make([]float64, k)),
				vector.V(make([]float64, k)),
			).M()

			got := &w{}
			got.h, got.optimal = checkBF(c.b, c.f, c.g, c.opt, buf)
			if !epsilon.Within(got.h, c.want.h) {
				t.Errorf("h = %v, want = %v", got.h, c.want.h)
			}
			if got.optimal != c.want.optimal {
				t.Errorf("optimal = %v, want = %v", got.optimal, c.want.optimal)
			}
		})
	}
}

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

	configs := []config{
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        k,
			})

			l := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			l.AABB().Copy(*hyperrectangle.New(
				vector.V([]float64{1, 1}),
				vector.V([]float64{2, 2}),
			))
			l.SetHeight(2)

			r := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			r.AABB().Copy(*hyperrectangle.New(
				vector.V([]float64{9, 9}),
				vector.V([]float64{11, 11}),
			))
			r.SetHeight(1)

			return config{
				name: "Simple",
				l:    l,
				r:    r,
				want: w{
					height:   3,
					balanced: true,
					h: heuristic.H(*hyperrectangle.New(
						vector.V([]float64{1, 1}),
						vector.V([]float64{11, 11}),
					)),
				},
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			buf := hyperrectangle.New(
				vector.V(make([]float64, k)),
				vector.V(make([]float64, k)),
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

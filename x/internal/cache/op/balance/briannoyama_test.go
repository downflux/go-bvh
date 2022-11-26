package balance

import (
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/node/impl"
	"github.com/downflux/go-bvh/x/internal/cache/node/util/cmp"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

func TestBrianNoyama(t *testing.T) {
	type config struct {
		name string
		x    node.N
		data map[id.ID]hyperrectangle.R
		want node.N
	}

	configs := []config{
		func() config {
			data := map[id.ID]hyperrectangle.R{}

			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})

			x := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))

			want := impl.New(c, x.ID())
			want.Allocate(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid)

			return config{
				name: "NoOp/Root",
				x:    x,
				data: data,
				want: want,
			}
		}(),
		func() config {
			data := map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(
					vector.V([]float64{0, 0}),
					vector.V([]float64{1, 1}),
				),
				101: *hyperrectangle.New(
					vector.V([]float64{10, 10}),
					vector.V([]float64{11, 11}),
				),
			}

			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})

			x := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))

			l := c.GetOrDie(c.Insert(x.ID(), cid.IDInvalid, cid.IDInvalid, true))
			r := c.GetOrDie(c.Insert(x.ID(), cid.IDInvalid, cid.IDInvalid, true))

			x.SetLeft(l.ID())
			x.SetRight(r.ID())

			l.Leaves()[100] = struct{}{}
			r.Leaves()[101] = struct{}{}

			for _, n := range []node.N{r, l, x} {
				node.SetAABB(n, data, 1)
				node.SetHeight(n)
			}

			d := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})

			want := d.GetOrDie(d.Insert(cid.IDInvalid, l.ID(), r.ID(), false))
			wl := d.GetOrDie(d.Insert(x.ID(), cid.IDInvalid, cid.IDInvalid, false))
			wr := d.GetOrDie(d.Insert(x.ID(), cid.IDInvalid, cid.IDInvalid, false))

			wl.Leaves()[100] = struct{}{}
			wl.AABB().Copy(*hyperrectangle.New(
				vector.V([]float64{0, 0}),
				vector.V([]float64{1, 1}),
			))

			wr.Leaves()[101] = struct{}{}
			wr.AABB().Copy(*hyperrectangle.New(
				vector.V([]float64{10, 10}),
				vector.V([]float64{11, 11}),
			))

			want.AABB().Copy(*hyperrectangle.New(
				vector.V([]float64{0, 0}),
				vector.V([]float64{11, 11}),
			))
			want.SetHeight(1)
			want.Left().SetHeuristic(heuristic.H(want.Left().AABB().R()))

			return config{
				name: "NoOp/Child",
				x:    x.Left(),
				data: data,
				want: want.Left(),
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := BrianNoyama(c.x, c.data, 1); !cmp.Equal(got, c.want) {
				t.Errorf("B() = %v, want = %v", got, c.want)
			}
		})
	}
}

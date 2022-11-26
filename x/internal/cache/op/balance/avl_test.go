package balance

import (
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/node/impl"
	"github.com/downflux/go-bvh/x/internal/cache/node/util/cmp"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

func TestAVL(t *testing.T) {
	type config struct {
		name string
		x    node.N
		want node.N
	}

	configs := []config{
		func() config {
			data := map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(
					vector.V([]float64{1, 1}),
					vector.V([]float64{2, 2}),
				),
			}
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})
			root := c.GetOrDie(c.Insert(
				cid.IDInvalid,
				cid.IDInvalid,
				cid.IDInvalid,
				/* validate = */ true,
			))
			root.Leaves()[100] = struct{}{}
			node.SetHeight(root)
			node.SetAABB(root, data, 1)

			want := impl.New(c, root.ID())
			want.Allocate(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid)
			want.Leaves()[100] = struct{}{}
			node.SetHeight(want)
			node.SetAABB(want, data, 1)

			return config{
				name: "NoOp/Height=0",
				x:    root,
				want: want,
			}
		}(),
		//    A
		//   / \
		//  B   C
		//     / \
		//    D   E
		//       / \
		//      F   G
		//
		// to
		//
		//    A
		//   / \
		//  E   C
		//     / \
		//    D   B
		func() config {
			data := map[id.ID]hyperrectangle.R{
				// B
				100: *hyperrectangle.New(
					vector.V([]float64{1, 1}),
					vector.V([]float64{2, 2}),
				),
				// D
				101: *hyperrectangle.New(
					vector.V([]float64{2, 2}),
					vector.V([]float64{4, 4}),
				),
				// F
				102: *hyperrectangle.New(
					vector.V([]float64{8, 8}),
					vector.V([]float64{16, 16}),
				),
				// G
				103: *hyperrectangle.New(
					vector.V([]float64{100, 100}),
					vector.V([]float64{200, 200}),
				),
			}

			c := cache.New(cache.O{
				K:        2,
				LeafSize: 1,
			})
			na := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			nb := c.GetOrDie(c.Insert(na.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nc := c.GetOrDie(c.Insert(na.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nd := c.GetOrDie(c.Insert(nc.ID(), cid.IDInvalid, cid.IDInvalid, true))
			ne := c.GetOrDie(c.Insert(nc.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nf := c.GetOrDie(c.Insert(ne.ID(), cid.IDInvalid, cid.IDInvalid, true))
			ng := c.GetOrDie(c.Insert(ne.ID(), cid.IDInvalid, cid.IDInvalid, true))

			na.SetLeft(nb.ID())
			na.SetRight(nc.ID())

			nb.Leaves()[100] = struct{}{}

			nc.SetLeft(nd.ID())
			nc.SetRight(ne.ID())

			nd.Leaves()[101] = struct{}{}

			ne.SetLeft(nf.ID())
			ne.SetRight(ng.ID())

			nf.Leaves()[102] = struct{}{}

			ng.Leaves()[103] = struct{}{}

			for _, n := range []node.N{ng, nf, ne, nd, nc, nb, na} {
				node.SetAABB(n, data, 1)
				node.SetHeight(n)
			}

			//    A
			//   / \
			//  E   C
			//     / \
			//    D   B
			d := cache.New(cache.O{
				K:        2,
				LeafSize: 1,
			})
			wa := d.GetOrDie(d.Insert(cid.IDInvalid, ne.ID(), nc.ID(), false))
			wb := d.GetOrDie(d.Insert(nc.ID(), cid.IDInvalid, cid.IDInvalid, false))
			wc := d.GetOrDie(d.Insert(na.ID(), nd.ID(), nb.ID(), false))
			wd := d.GetOrDie(d.Insert(nc.ID(), cid.IDInvalid, cid.IDInvalid, false))
			we := d.GetOrDie(d.Insert(na.ID(), nf.ID(), ng.ID(), false))
			wf := d.GetOrDie(d.Insert(ne.ID(), cid.IDInvalid, cid.IDInvalid, false))
			wg := d.GetOrDie(d.Insert(ne.ID(), cid.IDInvalid, cid.IDInvalid, false))

			wb.Leaves()[100] = struct{}{}
			wd.Leaves()[101] = struct{}{}
			wf.Leaves()[102] = struct{}{}
			wg.Leaves()[103] = struct{}{}

			for _, n := range []node.N{wg, wf, wd, wb, we, wc, wa} {
				node.SetAABB(n, data, 1)
				node.SetHeight(n)
			}

			return config{
				name: "Swap/BE",
				x:    na,
				want: wa,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := AVL(c.x); !cmp.Equal(got, c.want) {
				t.Errorf("AVL() = %v, want = %v", got, c.want)
			}
		})
	}
}

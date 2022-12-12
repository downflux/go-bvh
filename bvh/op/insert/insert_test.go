package insert

import (
	"testing"

	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-bvh/internal/cache"
	"github.com/downflux/go-bvh/internal/cache/node"
	"github.com/downflux/go-bvh/internal/cache/node/util"
	"github.com/downflux/go-bvh/internal/cache/node/util/cmp"
	"github.com/downflux/go-bvh/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/internal/cache/id"
)

func TestInsert(t *testing.T) {
	type config struct {
		name      string
		c         *cache.C
		rid       cid.ID
		data      map[id.ID]hyperrectangle.R
		x         id.ID
		tolerance float64
		want      node.N
	}

	configs := []config{
		func() config {
			data := map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(vector.V{346, 0}, vector.V{347, 1}),
				101: *hyperrectangle.New(vector.V{239, 0}, vector.V{240, 1}),
				102: *hyperrectangle.New(vector.V{896, 0}, vector.V{897, 1}),
				103: *hyperrectangle.New(vector.V{826, 0}, vector.V{827, 1}),
			}

			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})

			na := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			nb := c.GetOrDie(c.Insert(na.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nc := c.GetOrDie(c.Insert(na.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nf := c.GetOrDie(c.Insert(nc.ID(), cid.IDInvalid, cid.IDInvalid, true))
			ng := c.GetOrDie(c.Insert(nc.ID(), cid.IDInvalid, cid.IDInvalid, true))

			na.SetLeft(nb.ID())
			na.SetRight(nc.ID())

			nc.SetLeft(nf.ID())
			nc.SetRight(ng.ID())

			nb.Leaves()[100] = struct{}{}
			nf.Leaves()[101] = struct{}{}
			ng.Leaves()[102] = struct{}{}

			for _, n := range []node.N{ng, nf, nc, nb, na} {
				node.SetAABB(n, data, 1)
				node.SetHeight(n)
			}

			wc := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})

			wna := wc.GetOrDie(wc.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))

			wnb := wc.GetOrDie(wc.Insert(wna.ID(), cid.IDInvalid, cid.IDInvalid, true))
			wnc := wc.GetOrDie(wc.Insert(wna.ID(), cid.IDInvalid, cid.IDInvalid, true))

			wnd := wc.GetOrDie(wc.Insert(wnb.ID(), cid.IDInvalid, cid.IDInvalid, true))
			wne := wc.GetOrDie(wc.Insert(wnb.ID(), cid.IDInvalid, cid.IDInvalid, true))

			wnf := wc.GetOrDie(wc.Insert(wnc.ID(), cid.IDInvalid, cid.IDInvalid, true))
			wng := wc.GetOrDie(wc.Insert(wnc.ID(), cid.IDInvalid, cid.IDInvalid, true))

			wna.SetLeft(wnb.ID())
			wna.SetRight(wnc.ID())

			wnb.SetLeft(wnd.ID())
			wnb.SetRight(wne.ID())

			wnc.SetLeft(wnf.ID())
			wnc.SetRight(wng.ID())

			wnd.Leaves()[103] = struct{}{}
			wne.Leaves()[102] = struct{}{}

			wnf.Leaves()[101] = struct{}{}
			wng.Leaves()[100] = struct{}{}

			for _, n := range []node.N{wng, wnf, wne, wnd, wnc, wnb, wna} {
				node.SetAABB(n, data, 1)
				node.SetHeight(n)
			}

			return config{
				name:      "Experimental",
				c:         c,
				data:      data,
				rid:       na.ID(),
				x:         103,
				tolerance: 1,
				want:      wna,
			}
		}(),
		func() config {
			data := map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}),
			}
			wc := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})
			wr := wc.GetOrDie(wc.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			wr.Leaves()[100] = struct{}{}
			wr.AABB().Copy(data[100])
			wr.SetHeuristic(heuristic.H(wr.AABB().R()))

			return config{
				name: "Trivial",
				c: cache.New(cache.O{
					LeafSize: 1,
					K:        2,
				}),
				data:      data,
				rid:       cid.IDInvalid,
				x:         100,
				tolerance: 1,
				want:      wr,
			}
		}(),
		func() config {
			data := map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}),
				101: *hyperrectangle.New(vector.V{1, 0}, vector.V{2, 1}),
			}
			wc := cache.New(cache.O{
				LeafSize: 2,
				K:        2,
			})
			wr := wc.GetOrDie(wc.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			wr.Leaves()[100] = struct{}{}
			wr.Leaves()[101] = struct{}{}
			wr.AABB().Copy(data[100])
			wr.AABB().Union(data[101])
			wr.SetHeuristic(heuristic.H(wr.AABB().R()))

			c := cache.New(cache.O{
				LeafSize: 2,
				K:        2,
			})
			root := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			root.Leaves()[100] = struct{}{}
			root.AABB().Copy(data[100])
			root.SetHeuristic(heuristic.H(root.AABB().R()))

			return config{
				name:      "Trivial/LargeLeaf",
				c:         c,
				data:      data,
				rid:       root.ID(),
				x:         101,
				tolerance: 1,
				want:      wr,
			}
		}(),
		func() config {
			data := map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}),
				101: *hyperrectangle.New(vector.V{1, 0}, vector.V{2, 1}),
			}
			wc := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})

			wna := wc.GetOrDie(wc.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			wnb := wc.GetOrDie(wc.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			wnc := wc.GetOrDie(wc.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))

			wna.SetLeft(wnb.ID())
			wna.SetRight(wnc.ID())

			wnb.SetParent(wna.ID())
			wnc.SetParent(wna.ID())

			wnb.Leaves()[100] = struct{}{}
			wnc.Leaves()[101] = struct{}{}

			for _, n := range []node.N{wnc, wnb, wna} {
				node.SetAABB(n, data, 1)
				node.SetHeight(n)
			}

			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})
			root := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			root.Leaves()[100] = struct{}{}

			node.SetAABB(root, data, 1)
			node.SetHeight(root)

			return config{
				name:      "Split",
				c:         c,
				data:      data,
				rid:       root.ID(),
				x:         101,
				tolerance: 1,
				want:      wna,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			f := cmp.F{
				ID:          false,
				LRInvariant: false,
			}

			got, mutations := Insert(c.c, c.rid, c.data, c.x, c.tolerance)

			if !f.Equal(got, c.want) {
				t.Errorf("insert() = %v, _, want = %v, _", got, c.want)
			}

			updates := map[id.ID]cid.ID{}
			for _, n := range mutations {
				for x := range n.Leaves() {
					if y, ok := updates[x]; ok {
						t.Errorf("AABB %v found in multiple nodes: %v, %v", x, n.ID(), y)
					}
					updates[x] = n.ID()
				}
			}

			if err := util.Validate(c.c, c.data, got); err != nil {
				t.Errorf("Validate() encountered an unexpected error: %v", err)
			}
		})
	}
}

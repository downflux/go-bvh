package bvh

import (
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/node/util/cmp"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
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

			got, _ := insert(c.c, c.rid, c.data, c.x, c.tolerance)

			if !f.Equal(got, c.want) {
				t.Errorf("insert() = %v, _, want = %v, _", got, c.want)
			}
		})
	}
}

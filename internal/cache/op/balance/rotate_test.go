package balance

import (
	"testing"

	"github.com/downflux/go-bvh/internal/cache"
	"github.com/downflux/go-bvh/internal/cache/node"
	"github.com/downflux/go-bvh/internal/cache/node/util/cmp"
	"github.com/downflux/go-bvh/internal/heuristic"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/internal/cache/id"
)

var (
	_ B = Rotate
	_ B = RotateNoDF
)

func TestRotate(t *testing.T) {
	type config struct {
		name string
		x    node.N
		want node.N
	}

	configs := []config{
		//    A
		//   / \
		//  B   C
		//     / \
		//    F   G
		func() config {
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

			nf.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}))
			ng.AABB().Copy(*hyperrectangle.New(vector.V{10, 0}, vector.V{11, 1}))

			nf.SetHeuristic(heuristic.H(nf.AABB().R()))
			ng.SetHeuristic(heuristic.H(ng.AABB().R()))

			nc.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{11, 1}))
			nb.AABB().Copy(*hyperrectangle.New(vector.V{9, 0}, vector.V{10, 1}))

			nc.SetHeuristic(heuristic.H(nc.AABB().R()))
			nb.SetHeuristic(heuristic.H(nb.AABB().R()))

			nc.SetHeight(1)

			na.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{11, 1}))
			na.SetHeuristic(heuristic.H(na.AABB().R()))

			na.SetHeight(2)

			//    A
			//   / \
			//  F   C
			//     / \
			//    B   G
			wc := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})

			wna := wc.GetOrDie(wc.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))

			wnb := wc.GetOrDie(wc.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			wnc := wc.GetOrDie(wc.Insert(wna.ID(), cid.IDInvalid, cid.IDInvalid, true))

			wnf := wc.GetOrDie(wc.Insert(wna.ID(), cid.IDInvalid, cid.IDInvalid, true))
			wng := wc.GetOrDie(wc.Insert(wnc.ID(), cid.IDInvalid, cid.IDInvalid, true))

			wnb.SetParent(wnc.ID())

			wna.SetLeft(wnf.ID())
			wna.SetRight(wnc.ID())

			wnc.SetLeft(wnb.ID())
			wnc.SetRight(wng.ID())

			wnf.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}))
			wng.AABB().Copy(*hyperrectangle.New(vector.V{10, 0}, vector.V{11, 1}))

			wnf.SetHeuristic(heuristic.H(wnf.AABB().R()))
			wng.SetHeuristic(heuristic.H(wng.AABB().R()))

			wnb.AABB().Copy(*hyperrectangle.New(vector.V{9, 0}, vector.V{10, 1}))
			wnc.AABB().Copy(*hyperrectangle.New(vector.V{9, 0}, vector.V{11, 1}))

			wnc.SetHeuristic(heuristic.H(wnc.AABB().R()))
			wnb.SetHeuristic(heuristic.H(wnb.AABB().R()))

			wnc.SetHeight(1)

			wna.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{11, 1}))
			wna.SetHeuristic(heuristic.H(wna.AABB().R()))

			wna.SetHeight(2)

			return config{
				name: "BF",
				x:    na,
				want: wna,
			}
		}(),
		func() config {
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

			nf.AABB().Copy(*hyperrectangle.New(vector.V{10, 0}, vector.V{11, 1}))
			ng.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}))

			nf.SetHeuristic(heuristic.H(nf.AABB().R()))
			ng.SetHeuristic(heuristic.H(ng.AABB().R()))

			nc.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{11, 1}))
			nb.AABB().Copy(*hyperrectangle.New(vector.V{9, 0}, vector.V{10, 1}))

			nc.SetHeuristic(heuristic.H(nc.AABB().R()))
			nb.SetHeuristic(heuristic.H(nb.AABB().R()))

			nc.SetHeight(1)

			na.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{11, 1}))
			na.SetHeuristic(heuristic.H(na.AABB().R()))

			na.SetHeight(2)

			//    A
			//   / \
			//  G   C
			//     / \
			//    F   B
			wc := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})

			wna := wc.GetOrDie(wc.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))

			wnb := wc.GetOrDie(wc.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			wnc := wc.GetOrDie(wc.Insert(wna.ID(), cid.IDInvalid, cid.IDInvalid, true))

			wnf := wc.GetOrDie(wc.Insert(wnc.ID(), cid.IDInvalid, cid.IDInvalid, true))
			wng := wc.GetOrDie(wc.Insert(wna.ID(), cid.IDInvalid, cid.IDInvalid, true))

			wnb.SetParent(wnc.ID())

			wna.SetLeft(wng.ID())
			wna.SetRight(wnc.ID())

			wnc.SetLeft(wnf.ID())
			wnc.SetRight(wnb.ID())

			wnf.AABB().Copy(*hyperrectangle.New(vector.V{10, 0}, vector.V{11, 1}))
			wng.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}))

			wnf.SetHeuristic(heuristic.H(wnf.AABB().R()))
			wng.SetHeuristic(heuristic.H(wng.AABB().R()))

			wnb.AABB().Copy(*hyperrectangle.New(vector.V{9, 0}, vector.V{10, 1}))
			wnc.AABB().Copy(*hyperrectangle.New(vector.V{9, 0}, vector.V{11, 1}))

			wnc.SetHeuristic(heuristic.H(wnc.AABB().R()))
			wnb.SetHeuristic(heuristic.H(wnb.AABB().R()))

			wnc.SetHeight(1)

			wna.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{11, 1}))
			wna.SetHeuristic(heuristic.H(wna.AABB().R()))

			wna.SetHeight(2)

			return config{
				name: "BG",
				x:    na,
				want: wna,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := Rotate(c.x); !cmp.Equal(got, c.want) {
				t.Errorf("Rotate() = %v, want = %v", got, c.want)
			}
		})
	}
}

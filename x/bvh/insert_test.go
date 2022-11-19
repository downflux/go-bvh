package bvh

import (
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/node/impl"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

func TestPartition(t *testing.T) {
	type config struct {
		name string
		s    node.N
		t    node.N
		axis vector.D
		data map[id.ID]hyperrectangle.R
		want node.N
	}

	configs := []config{
		func() config {
			data := map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(vector.V([]float64{0, 0}), vector.V([]float64{10, 10})),
				101: *hyperrectangle.New(vector.V([]float64{100, 0}), vector.V([]float64{110, 10})),
			}
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})
			s := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			s.Leaves()[100] = struct{}{}
			s.Leaves()[101] = struct{}{}
			t := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))

			want := impl.New(c, t.ID())
			want.Allocate(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid)
			want.Leaves()[101] = struct{}{}
			return config{
				name: "Trivial",
				s:    s,
				t:    t,
				axis: vector.AXIS_X,
				data: data,
				want: want,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			partition(c.s, c.t, c.axis, c.data)
			if !node.Equal(c.t, c.want) {
				t.Errorf("n = %v, want = %v", c.t, c.want)
			}
		})
	}
}

func TestRaw(t *testing.T) {
	const k = 2

	type w struct {
		s node.N
		t node.N
	}

	type config struct {
		name string
		c    *cache.C
		root cid.ID
		data map[id.ID]hyperrectangle.R
		x    id.ID
		want w
	}

	configs := []config{
		func() config {
			x := id.ID(100)

			c := cache.New(cache.O{
				LeafSize: 1,
				K:        k,
			})

			t := impl.New(c, 0)
			t.Allocate(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid)
			t.Leaves()[x] = struct{}{}

			return config{
				name: "Nil",
				c:    c,
				root: cid.IDInvalid,
				data: nil,
				x:    x,
				want: w{
					s: nil,
					t: t,
				},
			}
		}(),
		func() config {
			data := map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(vector.V{1, 1}, vector.V{10, 10}),
				101: *hyperrectangle.New(vector.V{100, 100}, vector.V{110, 110}),
			}

			c := cache.New(cache.O{
				LeafSize: 1,
				K:        k,
			})

			root := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			root.Leaves()[100] = struct{}{}

			s := impl.New(c, root.ID())
			s.Allocate(1, cid.IDInvalid, cid.IDInvalid)
			s.Leaves()[100] = struct{}{}

			t := impl.New(c, 2)
			t.Allocate(1, cid.IDInvalid, cid.IDInvalid)
			t.Leaves()[101] = struct{}{}

			r := impl.New(c, 1)
			r.Allocate(cid.IDInvalid, root.ID(), 2)

			return config{
				name: "Nil",
				c:    c,
				root: root.ID(),
				data: data,
				x:    101,
				want: w{
					s: s,
					t: t,
				},
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := &w{}
			got.s, got.t = raw(c.c, c.root, c.data, c.x)
			if !node.Equal(got.s, c.want.s) {
				t.Errorf("s = %v, want = %v", got.s, c.want.s)
			}
			if !node.Equal(got.t, c.want.t) {
				t.Errorf("t = %v, want = %v", got.t, c.want.t)
			}
		})
	}
}

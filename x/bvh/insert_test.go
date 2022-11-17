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

func TestExpand(t *testing.T) {
	type config struct {
		name string
		c    *cache.C
		s    node.N
		want node.N
	}

	configs := []config{
		func() config {
			c := cache.New(cache.O{
				K:        2,
				LeafSize: 1,
			})
			root := c.GetOrDie(c.Insert(
				cid.IDInvalid,
				cid.IDInvalid,
				cid.IDInvalid,
				/* validate = */ true,
			))
			n := impl.New(c, 2)
			n.Allocate(1, cid.IDInvalid, cid.IDInvalid)
			return config{
				name: "Root",
				c:    c,
				s:    root,
				want: n,
			}
		}(),
		func() config {
			c := cache.New(cache.O{
				K:        2,
				LeafSize: 1,
			})
			root := c.GetOrDie(c.Insert(
				cid.IDInvalid,
				cid.IDInvalid,
				cid.IDInvalid,
				/* validate = */ true,
			))
			s := c.GetOrDie(c.Insert(
				root.ID(),
				cid.IDInvalid,
				cid.IDInvalid,
				true,
			))
			root.SetLeft(s.ID())
			t := c.GetOrDie(c.Insert(
				root.ID(),
				cid.IDInvalid,
				cid.IDInvalid,
				true,
			))
			root.SetRight(t.ID())

			n := impl.New(c, 4)
			n.Allocate(3, cid.IDInvalid, cid.IDInvalid)
			return config{
				name: "Child",
				c:    c,
				s:    s,
				want: n,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := expand(c.c, c.s); !node.Equal(c.want, got) {
				t.Errorf("expand() = %v, want = %v", got, c.want)
			}
		})
	}
}

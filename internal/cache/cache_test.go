package cache

import (
	"fmt"
	"testing"

	"github.com/downflux/go-bvh/internal/cache/node"
	"github.com/downflux/go-bvh/internal/cache/node/impl"
	"github.com/downflux/go-bvh/perf/size"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"

	cid "github.com/downflux/go-bvh/internal/cache/id"
	ncmp "github.com/downflux/go-bvh/internal/cache/node/util/cmp"
)

func TestDelete(t *testing.T) {
	c := New(O{
		K:        1,
		LeafSize: 1,
	})
	x := c.Insert(-1, -1, -1, true)
	c.DeleteOrDie(x)
	if _, ok := c.Get(x); ok {
		t.Errorf("Get() = %v, want = %v", ok, false)
	}
}

func TestInsert(t *testing.T) {
	type config struct {
		name    string
		c       *C
		p, l, r cid.ID
		want    node.N
	}
	configs := []config{
		func() config {
			c := New(O{
				K:        1,
				LeafSize: 1,
			})
			n := impl.New(c, 0)
			n.Allocate(-1, -1, -1)
			return config{
				name: "Empty",
				c:    c,
				p:    -1,
				l:    -1,
				r:    -1,
				want: n,
			}
		}(),
		func() config {
			c := New(O{
				K:        1,
				LeafSize: 1,
			})
			c.Insert(-1, -1, -1, true)
			n := impl.New(c, 1)
			n.Allocate(-1, -1, -1)
			return config{
				name: "AfterInsert",
				c:    c,
				p:    -1,
				l:    -1,
				r:    -1,
				want: n,
			}
		}(),
		func() config {
			c := New(O{
				K:        1,
				LeafSize: 1,
			})
			c.DeleteOrDie(c.Insert(-1, -1, -1, true))
			n := impl.New(c, 0)
			n.Allocate(-1, -1, -1)
			return config{
				name: "AfterFree",
				c:    c,
				p:    -1,
				l:    -1,
				r:    -1,
				want: n,
			}
		}(),
	}
	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := c.c.GetOrDie(c.c.Insert(c.p, c.l, c.r, true))
			if !ncmp.Equal(got, c.want) {
				diff := cmp.Diff(c.want, got, cmp.AllowUnexported(
					C{}, hyperrectangle.M{}, hyperrectangle.R{},
				))
				t.Errorf("GetOrDie() mismatch(-want +got):\n%v", diff)
			}
		})
	}
}

func BenchmarkInsert(b *testing.B) {
	for _, n := range size.SizeUnit.N() {
		b.Run(fmt.Sprintf("Sequential/Batch=%v", n), func(b *testing.B) {
			b.StopTimer()
			cache := New(O{
				K:        1,
				LeafSize: 1,
			})
			b.StartTimer()

			for i := 0; i < b.N; i++ {
				for j := 0; j < n; j++ {
					cache.Insert(-1, -1, -1, false)
				}
			}
		})
		b.Run(fmt.Sprintf("Freed/Batch=%v", n), func(b *testing.B) {
			b.StopTimer()
			cache := New(O{
				K:        1,
				LeafSize: 1,
			})
			b.StartTimer()

			for i := 0; i < b.N; i++ {
				for j := 0; j < n; j++ {
					cache.Insert(-1, -1, -1, false)
				}

				b.StopTimer()
				for j := 0; j < n; j++ {
					cache.Delete(cid.ID(j))
				}
				b.StartTimer()
			}
		})
	}
}

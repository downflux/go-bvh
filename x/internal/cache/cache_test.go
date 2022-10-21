package cache

import (
	"fmt"
	"testing"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/google/go-cmp/cmp"
)

func TestDelete(t *testing.T) {
	c := New(O{
		K: 1,
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
		p, l, r ID
		want    *N
	}
	configs := []config{
		func() config {
			c := New(O{
				K: 1,
			})
			return config{
				name: "Empty",
				c:    c,
				p:    -1,
				l:    -1,
				r:    -1,
				want: &N{
					cache:       c,
					isAllocated: true,
					ids:         [4]ID{0, -1, -1, -1},
					aabbCache: hyperrectangle.New(
						vector.V([]float64{0}),
						vector.V([]float64{0}),
					).M(),
				},
			}
		}(),
		func() config {
			c := New(O{
				K: 1,
			})
			c.Insert(-1, -1, -1, true)
			return config{
				name: "AfterInsert",
				c:    c,
				p:    -1,
				l:    -1,
				r:    -1,
				want: &N{
					cache:       c,
					isAllocated: true,
					ids:         [4]ID{1, -1, -1, -1},
					aabbCache: hyperrectangle.New(
						vector.V([]float64{0}),
						vector.V([]float64{0}),
					).M(),
				},
			}
		}(),
		func() config {
			c := New(O{
				K: 1,
			})
			c.DeleteOrDie(c.Insert(-1, -1, -1, true))
			return config{
				name: "AfterFree",
				c:    c,
				p:    -1,
				l:    -1,
				r:    -1,
				want: &N{
					cache:       c,
					isAllocated: true,
					ids:         [4]ID{0, -1, -1, -1},
					aabbCache: hyperrectangle.New(
						vector.V([]float64{0}),
						vector.V([]float64{0}),
					).M(),
				},
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := c.c.GetOrDie(c.c.Insert(c.p, c.l, c.r, true))
			if !got.Within(c.want) {
				diff := cmp.Diff(c.want, got, cmp.AllowUnexported(
					N{}, C{}, hyperrectangle.M{}, hyperrectangle.R{},
				))
				t.Errorf("GetOrDie() mismatch(-want +got):\n%v", diff)
			}
		})
	}
}

func BenchmarkInsert(b *testing.B) {
	const batch = 1e4

	b.Run(fmt.Sprintf("Sequential/Batch=%v", batch), func(b *testing.B) {
		b.StopTimer()
		cache := New(O{
			K: 1,
		})
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			for j := 0; j < batch; j++ {
				cache.Insert(-1, -1, -1, false)
			}
		}
	})
	b.Run(fmt.Sprintf("Freed/Batch=%v", batch), func(b *testing.B) {
		b.StopTimer()
		cache := New(O{
			K: 1,
		})
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			for j := 0; j < batch; j++ {
				cache.Insert(-1, -1, -1, false)
			}

			b.StopTimer()
			for j := 0; j < batch; j++ {
				cache.Delete(ID(j))
			}
			b.StartTimer()
		}
	})

}

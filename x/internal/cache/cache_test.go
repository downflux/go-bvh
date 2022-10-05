package cache

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDelete(t *testing.T) {
	c := New()
	x := c.Insert(-1, -1, -1)
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
		{
			name: "Empty",
			c:    New(),
			p:    -1,
			l:    100,
			r:    101,
			want: &N{
				isAllocated: true,
				ids:         [4]ID{0, -1, 100, 101},
			},
		},
		func() config {
			c := New()
			c.Insert(-1, 100, 101)
			return config{
				name: "AfterInsert",
				c:    c,
				p:    -1,
				l:    102,
				r:    103,
				want: &N{
					isAllocated: true,
					ids:         [4]ID{1, -1, 102, 103},
				},
			}
		}(),
		func() config {
			c := New()
			c.DeleteOrDie(c.Insert(-1, 100, 101))
			return config{
				name: "AfterFree",
				c:    c,
				p:    -1,
				l:    102,
				r:    103,
				want: &N{
					isAllocated: true,
					ids:         [4]ID{0, -1, 102, 103},
				},
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := c.c.GetOrDie(c.c.Insert(c.p, c.l, c.r))
			if diff := cmp.Diff(c.want, got, cmp.AllowUnexported(N{})); diff != "" {
				t.Errorf("GetOrDie() mismatch(-want +got):\n%v", diff)
			}
		})
	}
}

func BenchmarkInsert(b *testing.B) {
	const batch = 1e4

	b.Run(fmt.Sprintf("Sequential/Batch=%v", batch), func(b *testing.B) {
		b.StopTimer()
		cache := New()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			for j := 0; j < batch; j++ {
				cache.Insert(-1, -1, -1)
			}
		}
	})
	b.Run(fmt.Sprintf("Freed/Batch=%v", batch), func(b *testing.B) {
		b.StopTimer()
		cache := New()
		b.StartTimer()

		for i := 0; i < b.N; i++ {
			for j := 0; j < batch; j++ {
				cache.Insert(-1, -1, -1)
			}

			b.StopTimer()
			for j := 0; j < batch; j++ {
				cache.Delete(ID(j))
			}
			b.StartTimer()
		}
	})

}

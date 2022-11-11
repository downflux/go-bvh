package tree

import (
	"testing"

	"github.com/downflux/go-bvh/x/internal/cache"
)

func TestExpand(t *testing.T) {
	type config struct {
		name string
		c    *cache.C
		s    *cache.N
		want *cache.N
	}

	configs := []config{
		func() config {
			c := cache.New(cache.O{
				K:        2,
				LeafSize: 1,
			})
			root := c.GetOrDie(c.Insert(
				cache.IDInvalid,
				cache.IDInvalid,
				cache.IDInvalid,
				/* validate = */ true,
			))
			return config{
				name: "Root",
				c:    c,
				s:    root,
				want: cache.NewN(c, 2, 1, cache.IDInvalid, cache.IDInvalid),
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := expand(c.c, c.s); !c.want.Within(got) {
				t.Errorf("expand() = %v, want = %v", got, c.want)
			}
		})
	}
}

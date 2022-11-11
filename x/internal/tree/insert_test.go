package tree

import (
	"testing"

	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/shared"
	"github.com/downflux/go-bvh/x/internal/cache/node"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

func TestExpand(t *testing.T) {
	type config struct {
		name string
		c    *cache.C
		s    shared.N
		want shared.N
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
			n := node.New(c, 2)
			n.Allocate(1, cid.IDInvalid, cid.IDInvalid)
			return config{
				name: "Root",
				c:    c,
				s:    root,
				want: n,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := expand(c.c, c.s); !shared.Equal(c.want, got) {
				t.Errorf("expand() = %v, want = %v", got, c.want)
			}
		})
	}
}

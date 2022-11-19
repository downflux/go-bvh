package unsafe

import (
	"testing"

	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/node/impl"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

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
			if got := Expand(c.c, c.s); !node.Equal(c.want, got) {
				t.Errorf("Expand() = %v, want = %v", got, c.want)
			}
		})
	}
}

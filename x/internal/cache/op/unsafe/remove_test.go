package unsafe

import (
	"testing"

	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/node/impl"
	"github.com/downflux/go-bvh/x/internal/cache/node/util/cmp"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

func TestRemove(t *testing.T) {
	type config struct {
		name string
		c    *cache.C
		n    node.N
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
				true,
			))

			return config{
				name: "Root",
				c:    c,
				n:    root,
				want: nil,
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
				true,
			))

			l := c.GetOrDie(c.Insert(root.ID(), cid.IDInvalid, cid.IDInvalid, true))
			r := c.GetOrDie(c.Insert(root.ID(), cid.IDInvalid, cid.IDInvalid, true))

			root.SetLeft(l.ID())
			root.SetRight(r.ID())

			want := impl.New(c, l.ID())
			want.Allocate(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid)

			return config{
				name: "Leaf",
				c:    c,
				n:    r,
				want: want,
			}
		}(),
		func() config {
			c := cache.New(cache.O{
				K:        2,
				LeafSize: 1,
			})
			nq := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))

			nt := c.GetOrDie(c.Insert(nq.ID(), cid.IDInvalid, cid.IDInvalid, true))
			np := c.GetOrDie(c.Insert(nq.ID(), cid.IDInvalid, cid.IDInvalid, true))

			nm := c.GetOrDie(c.Insert(np.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nn := c.GetOrDie(c.Insert(np.ID(), cid.IDInvalid, cid.IDInvalid, true))

			nq.SetLeft(nt.ID())
			nq.SetRight(np.ID())

			np.SetLeft(nm.ID())
			np.SetRight(nn.ID())

			want := impl.New(c, nm.ID())
			want.Allocate(nq.ID(), cid.IDInvalid, cid.IDInvalid)

			return config{
				name: "Internal",
				c:    c,
				n:    nn,
				want: want,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := Remove(c.c, c.n); !cmp.Equal(c.want, got) {
				t.Errorf("Remove() = %v, want = %v", got, c.want)
			}
		})
	}
}

package cache

import (
	"testing"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"
)

func TestIsAncestor(t *testing.T) {
	type config struct {
		name string
		c    *C
		n    ID
		m    ID
		want bool
	}

	configs := []config{
		func() config {
			c := New(O{
				K: 1,
			})
			root := c.GetOrDie(c.Insert(-1, -1, -1, true))
			n := c.Insert(root.ID(), -1, -1, true)
			m := c.Insert(root.ID(), -1, -1, true)
			root.SetLeft(n)
			root.SetRight(m)

			return config{
				name: "Sibling",
				c:    c,
				n:    n,
				m:    m,
				want: false,
			}
		}(),
		func() config {
			c := New(O{
				K: 1,
			})
			root := c.GetOrDie(c.Insert(-1, -1, -1, true))
			n := c.Insert(root.ID(), -1, -1, true)
			root.SetLeft(n)

			return config{
				name: "Parent",
				c:    c,
				n:    root.ID(),
				m:    n,
				want: true,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := IsAncestor(c.c, c.n, c.m); got != c.want {
				t.Errorf("IsAncestor() = %v, want = %v", got, c.want)
			}
		})
	}
}

func TestSwap(t *testing.T) {
	type config struct {
		name string
		c    *C
		from ID
		to   ID
		want *C
	}

	configs := []config{
		func() config {
			c := New(O{
				K: 1,
			})
			root := c.GetOrDie(c.Insert(-1, -1, -1, true))
			n := c.Insert(root.ID(), -1, -1, true)
			m := c.Insert(root.ID(), -1, -1, true)
			root.SetLeft(n)
			root.SetRight(m)

			want := &C{
				k:     1,
				freed: []ID{},
			}
			want.data = []*N{
				&N{cache: want, ids: [4]ID{0, -1, 2, 1}, isAllocated: true},
				&N{cache: want, ids: [4]ID{1, 0, -1, -1}, isAllocated: true},
				&N{cache: want, ids: [4]ID{2, 0, -1, -1}, isAllocated: true},
			}

			return config{
				name: "Sibling",
				c:    c,
				from: n,
				to:   m,
				want: want,
			}
		}(),
		func() config {
			c := New(O{
				K: 1,
			})
			root := c.GetOrDie(c.Insert(-1, -1, -1, true))

			n := c.Insert(root.ID(), -1, -1, true)
			root.SetLeft(n)

			r := c.GetOrDie(c.Insert(root.ID(), -1, -1, true))
			root.SetRight(r.ID())

			m := c.Insert(r.ID(), -1, -1, true)
			r.SetLeft(m)

			want := &C{
				k:     1,
				freed: []ID{},
			}
			want.data = []*N{
				&N{cache: want, ids: [4]ID{0, -1, m, r.ID()}, isAllocated: true},
				&N{cache: want, ids: [4]ID{n, r.ID(), -1, -1}, isAllocated: true},
				&N{cache: want, ids: [4]ID{r.ID(), root.ID(), n, -1}, isAllocated: true},
				&N{cache: want, ids: [4]ID{m, root.ID(), -1, -1}, isAllocated: true},
			}

			return config{
				name: "NonAdjacent",
				c:    c,
				from: n,
				to:   m,
				want: want,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			Swap(c.c, c.from, c.to /* validate = */, true)
			if diff := cmp.Diff(c.want, c.c, cmp.AllowUnexported(C{}, N{}, hyperrectangle.R{}, hyperrectangle.M{})); diff != "" {
				t.Errorf("Swap() mismatch(-want +got):\n%v", diff)
			}

		})
	}
}

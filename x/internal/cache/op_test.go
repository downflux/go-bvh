package cache

import (
	"testing"

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
			c := New()
			root := c.GetOrDie(c.Insert(-1, -1, -1))
			n := c.Insert(root.ID(), -1, -1)
			m := c.Insert(root.ID(), -1, -1)
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
			c := New()
			root := c.GetOrDie(c.Insert(-1, -1, -1))
			n := c.Insert(root.ID(), -1, -1)
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
			c := New()
			root := c.GetOrDie(c.Insert(-1, -1, -1))
			n := c.Insert(root.ID(), -1, -1)
			m := c.Insert(root.ID(), -1, -1)
			root.SetLeft(n)
			root.SetRight(m)

			want := &C{
				freed: []ID{},
				data: []*N{
					&N{ids: [4]ID{0, -1, 2, 1}, isAllocated: true},
					&N{ids: [4]ID{1, 0, -1, -1}, isAllocated: true},
					&N{ids: [4]ID{2, 0, -1, -1}, isAllocated: true},
				},
			}

			return config{
				name: "Sibling",
				c:    c,
				from: 1, // n,
				to:   2, // m,
				want: want,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			Swap(c.c, c.from, c.to /* validate = */, true)
			if diff := cmp.Diff(c.want, c.c, cmp.AllowUnexported(C{}, N{})); diff != "" {
				t.Errorf("Swap() mismatch(-want +got):\n%v", diff)
			}

		})
	}
}

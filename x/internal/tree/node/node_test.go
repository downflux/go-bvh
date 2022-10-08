package node

import (
	"testing"

	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNew(t *testing.T) {
	type config struct {
		name string
		c    *cache.C
		id   cache.ID
		want *N
	}

	configs := []config{
		func() config {
			c := cache.New()
			x := c.Insert(-1, -1, -1)

			return config{
				name: "SingleEntry",
				c:    c,
				id:   x,
				want: &N{
					cache:  c,
					id:     x,
					parent: cache.IDInvalid,
					branch: cache.BInvalid,
				},
			}
		}(),
		func() config {
			c := cache.New()
			root := c.GetOrDie(c.Insert(-1, -1, -1))
			root.SetLeft(c.Insert(root.ID(), -1, -1))

			return config{
				name: "Child",
				c:    c,
				id:   root.Left(),
				want: &N{
					cache:  c,
					id:     root.Left(),
					parent: root.ID(),
					branch: cache.BLeft,
				},
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			got := New(c.c, c.id)
			if diff := cmp.Diff(c.want, got, cmp.AllowUnexported(N{}), cmpopts.IgnoreFields(N{}, "cache")); diff != "" {
				t.Errorf("New() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

func TestIsRoot(t *testing.T) {
	type config struct {
		name string
		n    *N
		want bool
	}

	configs := []config{
		func() config {
			c := cache.New()
			root := c.Insert(-1, -1, -1)

			return config{
				name: "Root",
				n:    New(c, root),
				want: true,
			}

		}(),
		func() config {
			c := cache.New()
			left := c.Insert(c.GetOrDie(c.Insert(-1, -1, -1)).ID(), -1, -1)

			return config{
				name: "Child",
				n:    New(c, left),
				want: false,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := c.n.IsRoot(); got != c.want {
				t.Errorf("IsRoot() = %v, want = %v", got, c.want)
			}
		})
	}
}

func TestIsLeaf(t *testing.T) {
	type config struct {
		name string
		n    *N
		want bool
	}

	configs := []config{
		func() config {
			c := cache.New()
			root := c.Insert(-1, -1, -1)

			return config{
				name: "Leaf",
				n:    New(c, root),
				want: true,
			}
		}(),
		func() config {
			c := cache.New()
			root := c.GetOrDie(c.Insert(-1, -1, -1))
			root.SetLeft(c.Insert(root.ID(), -1, -1))
			root.SetRight(c.Insert(root.ID(), -1, -1))

			return config{
				name: "Parent",
				n:    New(c, root.ID()),
				want: false,
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := c.n.IsLeaf(); got != c.want {
				t.Errorf("IsLeaf() = %v, want = %v", got, c.want)
			}
		})
	}
}

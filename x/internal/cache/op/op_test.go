package op

import (
	"testing"

	"github.com/downflux/go-bvh/x/internal/cache"
)

func TestIsAncestor(t *testing.T) {
	type config struct {
		name string
		c    *cache.C
		n    cache.ID
		m    cache.ID
		want bool
	}

	configs := []config{
		func() config {
			c := cache.New()
			root := c.Insert(-1, -1, -1)
			n := c.Insert(root, -1, -1)
			m := c.Insert(root, -1, -1)
			c.GetOrDie(root).SetLeft(n)
			c.GetOrDie(root).SetRight(m)

			return config{
				name: "Sibling",
				c:    c,
				n:    n,
				m:    m,
				want: false,
			}
		}(),
		func() config {
			c := cache.New()
			root := c.Insert(-1, -1, -1)
			n := c.Insert(root, -1, -1)
			c.GetOrDie(root).SetLeft(n)

			return config{
				name: "Parent",
				c:    c,
				n:    root,
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

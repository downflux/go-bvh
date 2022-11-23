package candidate

import (
	"time"
	"fmt"
	"testing"
	"runtime"

	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

var (
	tests = map[string]C{
		"Bittner":     Bittner,
		"BrianNoyama": BrianNoyama,
		"Catto":       Catto,
		"Guttman":     Guttman,
	}
)

func BenchmarkC(b *testing.B) {
	type scenario struct {
		c    *cache.C
		root node.N
		aabb hyperrectangle.R
	}
	type config struct {
		name string
		c    C
		s    func() scenario
	}

	scenarios := map[string]func() scenario{
		"Trivial": func() scenario {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        2,
			})

			root := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			root.AABB().Copy(*hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}))

			return scenario{
				c:    c,
				root: root,
				aabb: *hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}),
			}
		},
	}

	configs := []config{}
	for sn, s := range scenarios {
		for label, c := range tests {
			configs = append(configs, config{
				name: fmt.Sprintf("%v/%v", sn, label),
				c:    c,
				s:    s,
			})
		}
	}

	for _, c := range configs {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s := func() scenario {
					b.StopTimer()
					runtime.MemProfileRate = 0
					defer func() { runtime.MemProfileRate = 512 * 1024 }()
					defer b.StartTimer()

					time.Sleep(time.Second)
					return c.s()
				}()
				c.c(s.c, s.root, s.aabb)
			}
		})
	}
}

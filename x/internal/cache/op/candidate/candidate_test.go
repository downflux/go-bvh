package candidate

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/cache/op/candidate/perf"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

var (
	tests = map[string]C{
		"Bittner":     Bittner,
		"BrianNoyama": BrianNoyama,
		"Catto":       Catto,
		"Guttman":     Guttman,
	}
)

// BenchmarkC will check the performance of each candidate search function.
// Because these functions will mutate a cache, and because the
// setup-to-execution time ratio is unfavorable for a high number of runs, this
// test suite should be run with e.g.
//
//	go test -bench . -benchtime=0.01s
func BenchmarkC(b *testing.B) {
	type scenario struct {
		name      string
		generator perf.G
		aabb      hyperrectangle.R
		batch     int
	}

	type config struct {
		name      string
		c         C
		generator perf.G
		aabb      hyperrectangle.R
		batch     int
	}

	scenarios := []scenario{
		{
			name:      "Trivial",
			generator: perf.Trivial,
			aabb:      *hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}),
			batch:     10000,
		},
		{
			name:      "Balanced/N=10000",
			generator: func() (*cache.C, node.N) { return perf.Balanced(10000) },
			aabb:      *hyperrectangle.New(vector.V{10, 10}, vector.V{100, 100}),
			batch:     500,
		},
	}

	configs := []config{}
	for _, s := range scenarios {
		for label, c := range tests {
			configs = append(configs, config{
				name:      fmt.Sprintf("%v/%v", s.name, label),
				c:         c,
				generator: s.generator,
				aabb:      s.aabb,
				batch:     s.batch,
			})
		}
	}

	for _, c := range configs {
		type lookup struct {
			cache *cache.C
			root  node.N
		}

		var ls []lookup

		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if i%c.batch == 0 {

					ls = func(n int) []lookup {
						b.StopTimer()
						runtime.MemProfileRate = 0
						defer func() { runtime.MemProfileRate = 512 * 1024 }()
						defer b.StartTimer()

						ls := make([]lookup, 0, n)
						for i := 0; i < n; i++ {
							l := &lookup{}
							l.cache, l.root = c.generator()
							ls = append(ls, *l)
						}
						return ls
					}(c.batch)

				}

				l := ls[i%c.batch]
				c.c(l.cache, l.root, c.aabb)
			}
		})
	}
}

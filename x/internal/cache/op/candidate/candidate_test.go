package candidate

import (
	"fmt"
	"math"
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
		"Bittner":     bittnerRO,
		"BrianNoyama": BrianNoyama,
		"Catto":       cattoRO,
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
	}

	type config struct {
		name  string
		c     C
		cache *cache.C
		root  node.N
		aabb  hyperrectangle.R
	}

	scenarios := []scenario{
		{
			name:      "Trivial",
			generator: perf.Trivial,
			aabb:      *hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}),
		},
	}
	for i := 4; i < 8; i++ {
		n := int(math.Pow(10, float64(i)))
		scenarios = append(scenarios, scenario{
			name:      fmt.Sprintf("Balanced/N=%v", n),
			generator: func() (*cache.C, node.N) { return perf.Balanced(n) },
			aabb:      *hyperrectangle.New(vector.V{10, 0}, vector.V{11, 1}),
		})
	}

	configs := []config{}
	for _, s := range scenarios {
		c, n := func() (*cache.C, node.N) {
			runtime.MemProfileRate = 0
			defer func() { runtime.MemProfileRate = 512 * 1024 }()
			return s.generator()
		}()
		for label, f := range tests {
			configs = append(configs, config{
				name:  fmt.Sprintf("%v/%v", s.name, label),
				c:     f,
				cache: c,
				root:  n,
				aabb:  s.aabb,
			})
		}
	}

	for _, c := range configs {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.c(c.cache, c.root, c.aabb)
			}
		})
	}
}

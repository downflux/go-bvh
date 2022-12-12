package split

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/perf/size"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

const (
	k = 2
)

var (
	tests = map[string]S{
		"DHConnelly":    DHConnelly,
		"GuttmanLinear": GuttmanLinear,
	}
)

func BenchmarkS(b *testing.B) {
	const j = 9

	type config struct {
		name string
		s    S
		size int
	}

	data := map[id.ID]hyperrectangle.R{}
	for i := 0.0; i < math.Pow(2, j); i++ {
		x := float64(rand.Intn(int(math.Pow(2, j))))
		y := float64(rand.Intn(int(math.Pow(2, j))))
		data[id.ID(100+i)] = *hyperrectangle.New(vector.V{x, y}, vector.V{x + 1, y + 1})
	}

	configs := []config{}
	for l, s := range tests {
		// For small leaf sizes, the iteration time is too fast, and the
		// StopTimer / StartTimer invocations take too long.
		for _, size := range size.SizeUnit.LeafSize() {
			configs = append(configs, config{
				name: fmt.Sprintf("%v/LeafSize=%v", l, size),
				s:    s,
				size: int(size),
			})
		}
	}

	for _, c := range configs {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ch, n, m := func() (*cache.C, node.N, node.N) {
					b.StopTimer()
					runtime.MemProfileRate = 0
					defer func() { runtime.MemProfileRate = 512 * 1024 }()
					defer b.StartTimer()

					c := cache.New(cache.O{
						LeafSize: c.size,
						K:        k,
					})
					n := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
					m := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))

					for x := range data {
						if len(n.Leaves()) > c.LeafSize() {
							break
						}
						n.Leaves()[x] = struct{}{}
					}
					node.SetAABB(n, data, 1)

					return c, n, m
				}()

				c.s(ch, data, n, m)
			}
		})
	}
}

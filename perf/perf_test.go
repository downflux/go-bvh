package perf

import (
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"
	"testing"

	"github.com/downflux/go-bvh/bvh"
	"github.com/downflux/go-bvh/container"
	"github.com/downflux/go-bvh/container/briannoyama"
	"github.com/downflux/go-bvh/container/bruteforce"
	"github.com/downflux/go-bvh/container/dhconnelly"
	"github.com/downflux/go-bvh/perf/size"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

var (
	suite = size.SizeUnit
)

func TestMain(m *testing.M) {
	flag.Var(&suite, "performance_test_size", "performance test size, one of (large)")
	flag.Parse()

	os.Exit(m.Run())
}

type c struct {
	name string
	t    func() container.C
	n    int
	k    vector.D
	f    float64
	load []F
}

func generate() []c {
	cs := []c{}

	for _, n := range suite.N() {
		for _, k := range suite.K() {
			load := GenerateInsertLoad(n, 0, k)
			cs = append(cs,
				c{
					name: fmt.Sprintf("bruteforce/K=%v/N=%v", k, n),
					t:    func() container.C { return bruteforce.New() },
					n:    n,
					k:    k,
					load: load,
				},
			)
			if k == 3 {
				cs = append(cs, c{
					name: fmt.Sprintf("briannoyama/K=%v/N=%v", k, n),
					t:    func() container.C { return briannoyama.New() },
					n:    n,
					k:    k,
					load: load,
				})
			}

			for _, size := range suite.LeafSize() {
				k := k
				size := size

				cs = append(cs,
					c{
						name: fmt.Sprintf("downflux/K=%v/N=%v/LeafSize=%v", k, n, size),
						t: func() container.C {
							return bvh.New(bvh.O{
								K:         k,
								LeafSize:  int(size),
								Tolerance: 1.05,
							})
						},
						k:    k,
						n:    n,
						load: load,
					},
					c{
						name: fmt.Sprintf("dhconnelly/K=%v/N=%v/LeafSize=%v", k, n, size),
						t: func() container.C {
							return dhconnelly.New(dhconnelly.O{
								K:         k,
								MinBranch: 1,
								MaxBranch: int(size),
							})
						},
						n:    n,
						k:    k,
						load: load,
					},
				)
			}
		}
	}

	return cs
}

func BenchmarkBroadPhase(b *testing.B) {
	type config struct {
		name string
		t    func() container.C
		k    vector.D
		q    hyperrectangle.R
		load []F
	}

	configs := []config{}

	for _, c := range generate() {
		for _, f := range suite.F() {
			vmin := make([]float64, c.k)
			vmax := make([]float64, c.k)
			for i := vector.D(0); i < c.k; i++ {
				vmax[i] = math.Pow(5*float64(c.n)*f, 1./float64(c.k))
			}

			// N.B.: q is a constant fractional area of the overall
			// scene. This means that on average, the fraction of
			// objects covered by this rectangle remains constant,
			// and as such, we expect this benchmark to also scale
			// linearly with N.
			q := *hyperrectangle.New(vmin, vmax)

			configs = append(configs, config{
				name: fmt.Sprintf("%v/F=%v", c.name, f),
				t:    c.t,
				k:    c.k,
				q:    q,
				load: c.load,
			})
		}
	}

	for _, c := range configs {
		b.Run(c.name, func(b *testing.B) {
			t := c.t()

			load := func() []F {
				b.StopTimer()
				runtime.MemProfileRate = 0
				defer func() { runtime.MemProfileRate = 512 * 1024 }()
				defer b.StartTimer()

				for _, f := range c.load {
					f(t)
				}

				return GenerateBroadPhaseLoad(b.N, c.q)
			}()

			for i := 0; i < b.N; i++ {
				load[i](t)
			}
		})
	}
}

func BenchmarkInsert(b *testing.B) {
	type config struct {
		name string
		t    func() container.C
		k    vector.D
		load []F
	}

	configs := []config{}

	for _, c := range generate() {
		configs = append(configs, config{
			name: c.name,
			t:    c.t,
			k:    c.k,
			load: c.load,
		})
	}

	for _, c := range configs {
		b.Run(c.name, func(b *testing.B) {
			t := c.t()
			load := func() []F {
				b.StopTimer()
				runtime.MemProfileRate = 0
				defer func() { runtime.MemProfileRate = 512 * 1024 }()
				defer b.StartTimer()

				for _, f := range c.load {
					f(t)
				}

				return GenerateInsertLoad(b.N, len(c.load), c.k)
			}()

			for i := 0; i < b.N; i++ {
				if err := load[i](t); err != nil {
					b.Errorf("Insert() = %v, want = nil", err)
				}
			}
		})
	}
}

func BenchmarkRemove(b *testing.B) {
	type config struct {
		name string
		t    func() container.C
		k    vector.D
		load []F
	}

	configs := []config{}

	for _, c := range generate() {
		configs = append(configs, config{
			name: c.name,
			t:    c.t,
			k:    c.k,
			load: c.load,
		})
	}

	for _, c := range configs {
		b.Run(c.name, func(b *testing.B) {
			t := c.t()
			load := func() []F {
				b.StopTimer()
				runtime.MemProfileRate = 0
				defer func() { runtime.MemProfileRate = 512 * 1024 }()
				defer b.StartTimer()

				for _, f := range c.load {
					f(t)
				}

				return GenerateRemoveLoad(b.N, len(c.load), t, c.k)
			}()

			for i := 0; i < b.N; i++ {
				if err := load[i](t); err != nil {
					b.Errorf("Remove() = %v, want = nil", err)
				}
			}
		})
	}
}

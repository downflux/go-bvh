package perf

import (
	"flag"
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/downflux/go-bvh/bvh"
	"github.com/downflux/go-bvh/container"
	"github.com/downflux/go-bvh/perf/generator"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

var (
	suite = SizeSmall
	logd  = flag.String("log_directory", "", "log directory")
)

func TestMain(m *testing.M) {
	flag.Var(&suite, "performance_test_size", "performance test size, one of (large)")
	flag.Parse()

	os.Exit(m.Run())
}

func BenchmarkInsert(b *testing.B) {
	type config struct {
		t    container.C
		name string
		n    int
		k    vector.D

		load generator.G
	}

	var configs []config
	for _, n := range suite.N() {
		for _, k := range suite.K() {
			configs = append(configs, func() config {
				t := generator.BY(generator.InsertRandom(nil, k, n))
				return config{
					name: fmt.Sprintf("briannoyama/K=%v/N=%v", k, n),
					t:    t,
					n:    n,
					k:    k,

					load: func(k vector.D, n int) []generator.M {
						return generator.InsertRandom(t.IDs(), k, n)
					},
				}
			}(),
			)
			for _, size := range suite.LeafSize() {
				configs = append(configs, func() config {
					t := generator.BVH(size, generator.InsertRandom(nil, k, n))
					return config{
						name: fmt.Sprintf("BVH/K=%v/N=%v/LeafSize=%v", k, n, size),
						t:    t,
						n:    n,
						k:    k,

						load: func(k vector.D, n int) []generator.M {
							return generator.InsertRandom(t.IDs(), k, n)
						},
					}
				}(),
				)
			}
		}
	}

	for _, c := range configs {
		b.Run(c.name, func(b *testing.B) {
			b.StopTimer()
			fs := c.load(c.k, b.N)
			b.StartTimer()

			for i := 0; i < b.N; i++ {
				if err := fs[i](c.t); err != nil {
					b.Errorf("Insert() = %v, want = nil", err)
				}
			}
			if bvh, ok := c.t.(*bvh.BVH); ok {
				b.StopTimer()
				m := bvh.Report()
				b.ReportMetric(m.SAH, "SAH")
				b.ReportMetric(m.LeafSize, "size")
				b.StartTimer()
			}
		})
	}
}

func BenchmarkBroadPhase(b *testing.B) {
	type config struct {
		name string
		t    container.C
		n    int
		k    vector.D
		f    float64
		size uint
	}

	var configs []config
	for _, n := range suite.N() {
		for _, k := range suite.K() {
			for _, f := range suite.F() {
				ms := generator.InsertRandom(nil, k, n)
				configs = append(configs, config{
					name: fmt.Sprintf("BruteForce/K=%v/N=%v/F=%v", k, n, f),
					t:    generator.BF(ms),
					n:    n,
					k:    k,
					f:    f,
				}, config{
					name: fmt.Sprintf("bryannoyama/K=%v/N=%v/F=%v", k, n, f),
					t:    generator.BY(ms),
					n:    n,
					k:    k,
					f:    f,
				})
				for _, size := range suite.LeafSize() {
					configs = append(configs, config{
						name: fmt.Sprintf("BVH/K=%v/N=%v/F=%v/LeafSize=%v", k, n, f, size),
						t:    generator.BVH(size, ms),
						n:    n,
						k:    k,
						f:    f,
						size: size,
					})
				}
			}
		}
	}

	for _, c := range configs {
		b.Run(c.name, func(b *testing.B) {
			b.StopTimer()
			vmin := make([]float64, c.k)
			vmax := make([]float64, c.k)
			for i := vector.D(0); i < c.k; i++ {
				vmax[i] = math.Pow(5*float64(c.n)*c.f, 1./float64(c.k))
			}
			q := *hyperrectangle.New(vmin, vmax)
			b.StartTimer()

			for i := 0; i < b.N; i++ {
				c.t.BroadPhase(q)
			}
		})
	}
}

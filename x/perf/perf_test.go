package perf

import (
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"
	"testing"

	"github.com/downflux/go-bvh/x/bvh"
	"github.com/downflux/go-bvh/x/container"
	"github.com/downflux/go-bvh/x/container/briannoyama"
	"github.com/downflux/go-bvh/x/container/bruteforce"
	"github.com/downflux/go-bvh/x/container/dhconnelly"
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

var (
	suite = SizeUnit
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
	load map[id.ID]hyperrectangle.R
}

func generate() []c {
	cs := []c{}

	for _, n := range suite.N() {
		for _, k := range suite.K() {

			load := GenerateRandomTiles(n, k)
			cs = append(cs,
				c{
					name: fmt.Sprintf("bruteforce/K=%v/N=%v", k, n),
					t:    func() container.C { return bruteforce.New() },
					n:    n,
					k:    k,
					load: load,
				},
				c{
					name: fmt.Sprintf("briannoyama/K=%v/N=%v", k, n),
					t:    func() container.C { return briannoyama.New() },
					n:    n,
					k:    k,
					load: load,
				},
			)

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
		load map[id.ID]hyperrectangle.R
		q    hyperrectangle.R
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
				load: c.load,
				q:    q,
			})
		}
	}

	for _, c := range configs {
		b.Run(c.name, func(b *testing.B) {
			t := c.t()

			func() {
				b.StopTimer()
				runtime.MemProfileRate = 0
				defer func() { runtime.MemProfileRate = 512 * 1024 }()
				defer b.StartTimer()

				Insert(t, c.load)
			}()

			for i := 0; i < b.N; i++ {
				t.BroadPhase(c.q)
			}
		})
	}
}

func BenchmarkInsert(b *testing.B) {
	type config struct {
		name string
		t    func() container.C
		k    vector.D
		load map[id.ID]hyperrectangle.R
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
			type obj struct {
				id   id.ID
				aabb hyperrectangle.R
			}

			t := c.t()
			data := func() []obj {
				b.StopTimer()
				runtime.MemProfileRate = 0
				defer func() { runtime.MemProfileRate = 512 * 1024 }()
				defer b.StartTimer()

				Insert(t, c.load)

				data := make([]obj, 0, b.N)
				offset := id.ID(len(c.load))
				for x, aabb := range GenerateRandomTiles(b.N, c.k) {
					data = append(data, obj{id: x + offset, aabb: aabb})
				}
				return data
			}()

			for i := 0; i < b.N; i++ {
				if err := t.Insert(data[i].id, data[i].aabb); err != nil {
					b.Errorf("Insert() = %v, want = nil", err)
				}
			}

			b.StopTimer()
			defer b.StartTimer()

			if u, ok := t.(*bvh.T); ok {
				b.ReportMetric(u.SAH(), "SAH")
				b.ReportMetric(u.LeafSize(), "LeafSize")
				b.ReportMetric(float64(u.H()), "H")
			} else if u, ok := t.(*briannoyama.BVH); ok {
				b.ReportMetric(u.SAH(), "SAH")
			}
		})
	}
}

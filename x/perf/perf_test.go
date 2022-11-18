package perf

import (
	"flag"
	"fmt"
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
	suite = SizeSmall
)

func TestMain(m *testing.M) {
	flag.Var(&suite, "performance_test_size", "performance test size, one of (large)")
	flag.Parse()

	os.Exit(m.Run())
}

func BenchmarkInsert(b *testing.B) {
	type config struct {
		name string
		t    func() container.C
		k    vector.D
		load map[id.ID]hyperrectangle.R
	}

	var configs []config

	for _, n := range suite.N() {
		for _, k := range suite.K() {
			load := GenerateRandomTiles(n, k)
			configs = append(configs,
				config{
					name: fmt.Sprintf("bruteforce/K=%v/N=%v", k, n),
					t:    func() container.C { return bruteforce.New() },
					k:    k,
					load: load,
				},
				config{
					name: fmt.Sprintf("briannoyama/K=%v/N=%v", k, n),
					t:    func() container.C { return briannoyama.New() },
					k:    k,
					load: load,
				},
			)

			for _, size := range suite.LeafSize() {
				configs = append(configs,
					config{
						name: fmt.Sprintf("downflux/K=%v/N=%v/LeafSize=%v", k, n, size),
						t: func() container.C {
							return bvh.New(bvh.O{
								K:         k,
								LeafSize:  int(size),
								Tolerance: 1.05,
							})
						},
						k:    k,
						load: load,
					},
				)
			}

			for _, size := range suite.LeafSize() {
				configs = append(configs,
					config{
						name: fmt.Sprintf("dhconnelly/K=%v/N=%v/LeafSize=%v", k, n, size),
						t: func() container.C {
							return dhconnelly.New(dhconnelly.O{
								K:         k,
								MinBranch: 1,
								MaxBranch: int(size),
							})
						},
						k:    k,
						load: load,
					},
				)
			}
		}
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
			if u, ok := t.(*bvh.T); ok {
				b.ReportMetric(u.SAH(), "SAH")
				b.ReportMetric(float64(u.H()), "H")
			}
			defer b.StartTimer()
		})
	}
}

package perf

import (
	"fmt"
	"math"
	"testing"

	"github.com/downflux/go-bvh/bvh"
)

func BenchmarkNew(b *testing.B) {
	type config struct {
		name string
		n    int
		k    int
		size int
	}

	var configs []config
	for _, n := range []int{1e3, 1e4, 1e5, 1e6} {
		for _, k := range []int{16} {
			for _, size := range []int{1, 4, 16, 64} {
				configs = append(configs, config{
					name: fmt.Sprintf("K=%v/N=%v/LeafSize=%v", k, n, size),
					n:    n,
					k:    k,
					size: size,
				})
			}
		}
	}

	for _, c := range configs {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				New(O{
					Insert: 1,
					K:      c.k,
					N:      c.n,
					Size:   c.size,
				}).Apply(0, 500)
			}
		})
	}
}

func BenchmarkBroadPhase(b *testing.B) {
	type config struct {
		name string
		bvh  *bvh.BVH
		k    int
		f    float64
		size int
	}

	var configs []config
	for _, n := range []int{1e3, 1e4, 1e5, 1e6} {
		for _, k := range []int{16} {
			for _, f := range []float64{0.05} {
				for _, size := range []int{1, 4, 16, 64} {
					l := New(O{
						Insert: 1,
						K:      k,
						N:      n,
						Size:   size,
					})
					configs = append(configs, config{
						name: fmt.Sprintf("K=%v/N=%v/F=%v/LeafSize=%v", k, n, f, size),
						bvh:  l.Apply(0, 500),
						k:    k,
						f:    f,
					})
				}
			}
		}
	}

	for _, c := range configs {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				bvh.BroadPhase(c.bvh, rr(0, 500*math.Pow(c.f, 1./float64(c.k)), c.k))
			}
		})
	}
}

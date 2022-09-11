package perf

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"testing"

	"github.com/downflux/go-bvh/perf/generator"
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

func l(d string, fn string) *log.Logger {
	if d == "" {
		return nil
	}
	f, err := os.OpenFile(path.Join(d, fmt.Sprintf("%v.log", fn)), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("could not create logger: %v", err))
	}
	return log.New(f, "", log.Lshortfile)
}

func BenchmarkInsert(b *testing.B) {
	type config struct {
		name string
		n    int
		k    vector.D
		size uint
	}

	var configs []config
	for _, n := range suite.N() {
		for _, k := range suite.K() {
			for _, size := range suite.LeafSize() {
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
		t := generator.BVH(c.size, generator.Generate(generator.O{
			Insert: 1,
			K:      c.k,
		}, c.n))
		b.Run(fmt.Sprintf("Real/%v", c.name), func(b *testing.B) {
			fs := generator.Generate(generator.O{
				IDs:    t.IDs(),
				Insert: 1,
				K:      c.k,
			}, b.N)

			for i := 0; i < b.N; i++ {
				if err := fs[i](t); err != nil {
					b.Errorf("Insert() = %v, want = nil", err)
				}
			}
		})
	}
}

func BenchmarkBroadPhase(b *testing.B) {
	type config struct {
		name string
		n    int
		k    vector.D
		f    float64
		size uint
	}

	var configs []config
	for _, n := range suite.N() {
		for _, k := range suite.K() {
			for _, f := range suite.F() {
				for _, size := range suite.LeafSize() {
					configs = append(configs, config{
						name: fmt.Sprintf("K=%v/N=%v/F=%v/LeafSize=%v", k, n, f, size),
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
		ms := generator.Generate(generator.O{
			Insert: 1,
			K:      c.k,
		}, c.n)
		t := generator.BVH(c.size, ms)
		bf := generator.BF(ms)
		q := RR(0, 500*math.Pow(c.f, 1./float64(c.k)), c.k)

		if c.size == 1 {
			b.Run(fmt.Sprintf("BruteForce/%v", c.name), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					bf.BroadPhase(q)
				}
			})
		}
		b.Run(fmt.Sprintf("Real/%v", c.name), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				t.BroadPhase(q)
			}
		})
	}
}

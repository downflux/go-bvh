package perf

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"testing"

	"github.com/downflux/go-bvh/bvh"
)

var (
	suite = SizeLarge
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

func BenchmarkNew(b *testing.B) {
	type config struct {
		name string
		n    int
		k    int
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
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				New(O{
					Insert: 1,
					K:      c.k,
					N:      c.n,
					Size:   c.size,
					Logger: l(*logd, fmt.Sprintf("new-%v-%v-%v", c.k, c.n, c.size)),
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
		size uint
	}

	var configs []config
	for _, n := range suite.N() {
		for _, k := range suite.K() {
			for _, f := range suite.F() {
				for _, size := range suite.LeafSize() {
					l := New(O{
						Insert: 1,
						K:      k,
						N:      n,
						Size:   size,
						Logger: l(*logd, fmt.Sprintf("broadphase-%v-%v-%v-%v", k, n, f, size)),
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

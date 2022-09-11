package perf

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/downflux/go-bvh/bvh"
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type PerfTestSize int

const (
	SizeUnknown PerfTestSize = iota
	SizeSmall
	SizeLarge
)

func (s *PerfTestSize) String() string {
	return map[PerfTestSize]string{
		SizeLarge: "large",
	}[*s]
}

func (s *PerfTestSize) Set(v string) error {
	size, ok := map[string]PerfTestSize{
		"large": SizeLarge,
		"small": SizeSmall,
	}[v]
	if !ok {
		return fmt.Errorf("invalid test size value: %v", v)
	}
	*s = size
	return nil
}

func (s PerfTestSize) N() []int {
	return map[PerfTestSize][]int{
		SizeLarge: []int{1e3, 1e4, 1e5},
		SizeSmall: []int{1e3},
	}[s]
}

func (s PerfTestSize) F() []float64 {
	return map[PerfTestSize][]float64{
		SizeLarge: []float64{0.05},
		SizeSmall: []float64{0.05},
	}[s]
}

func (s PerfTestSize) LeafSize() []uint {
	return map[PerfTestSize][]uint{
		SizeLarge: []uint{1, 16, 256, 1024},
		SizeSmall: []uint{1, 4},
	}[s]
}

func (s PerfTestSize) K() []int {
	return map[PerfTestSize][]int{
		SizeLarge: []int{2, 3},
		SizeSmall: []int{2},
	}[s]
}

func rn(min, max float64) float64 { return rand.Float64()*(max-min) + min }
func rv(min, max float64, k int) vector.V {
	vs := []float64{}
	for i := 0; i < k; i++ {
		vs = append(vs, rn(min, max))
	}
	return vector.V(vs)
}
func rr(min, max float64, k int) hyperrectangle.R {
	a := rv(min, max, k)
	b := rv(min, max, k)

	vmin := make([]float64, k)
	vmax := make([]float64, k)

	for i := 0; i < k; i++ {
		vmin[i] = math.Min(a[i], b[i])
		vmax[i] = math.Max(a[i], b[i])
	}
	return *hyperrectangle.New(vmin, vmax)
}

type O struct {
	// Insert is the insert node weight.
	Insert float64
	// Remove is the remove node weight.
	Remove float64

	// N is the number of opts to call.
	N int

	K int

	Size   uint
	Logger *log.Logger
}

type L struct {
	insert float64
	remove float64
	n      int
	k      int

	bvh *bvh.BVH
	ids map[id.ID]bool
}

func New(o O) *L {
	return &L{
		insert: o.Insert / (o.Insert + o.Remove),
		remove: o.Remove / (o.Insert + o.Remove),
		n:      o.N,
		k:      o.K,
		bvh: bvh.New(bvh.O{
			Size:   o.Size,
			Logger: o.Logger,
		}),
		ids: map[id.ID]bool{},
	}
}

func (l *L) IDs() []id.ID {
	ids := make([]id.ID, 0, len(l.ids))
	for i := range l.ids {
		ids = append(ids, i)
	}
	return ids
}

func (l *L) BVH() *bvh.BVH { return l.bvh }
func (l *L) Apply() *bvh.BVH {
	for _, f := range l.Generate() {
		f()
	}
	return l.BVH()
}

func (l *L) Generate() []func() {
	tiles := make([]int, 0, l.n)
	for i := 0; i < l.n; i++ {
		tiles = append(tiles, i)
	}
	rand.Shuffle(len(tiles), func(i, j int) { tiles[i], tiles[j] = tiles[j], tiles[i] })

	fs := make([]func(), 0, l.n)

	for i := 0; i < l.n; i++ {
		j := id.ID(i + 1)
		l.ids[j] = true

		vmin := make([]float64, l.k)
		vmax := make([]float64, l.k)
		for k := 0; k < l.k; k++ {
			vmax[k] = 1
		}
		vmin[0] = float64(tiles[i])
		vmax[0] = float64(tiles[i]) + 1

		fs = append(fs, func() { l.bvh.Insert(j, *hyperrectangle.New(vmin, vmax)) })
	}
	return fs
}

package generator

import (
	"math"
	"math/rand"
	"runtime"

	"github.com/downflux/go-bvh/bvh"
	"github.com/downflux/go-bvh/container"
	"github.com/downflux/go-bvh/container/briannoyama"
	"github.com/downflux/go-bvh/container/bruteforce"
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type M func(c container.C) error
type G func(k vector.D, n int) []M

type O struct {
	// IDs is current list of objects in the BVH -- this is used for
	// generating valid remove calls.
	IDs    []id.ID
	Insert float64
	Remove float64
	K      vector.D
}

func allocateID(ids map[id.ID]bool) id.ID {
	var x id.ID
	for ; ids[x]; x = id.ID(rand.Uint64()) {
	}
	return x
}

func BY(ms []M) *briannoyama.BVH {
	runtime.MemProfileRate = 0
	defer func() { runtime.MemProfileRate = 512 * 1024 }()

	t := briannoyama.New()
	for _, f := range ms {
		f(t)
	}
	return t
}

func BVH(size uint, ms []M) *bvh.BVH {
	// Generating large number of points in tests will mess with data
	// collection figures. We should ignore these allocs.
	runtime.MemProfileRate = 0
	defer func() { runtime.MemProfileRate = 512 * 1024 }()

	t := bvh.New(bvh.O{Size: size})
	for _, f := range ms {
		f(t)
	}
	return t
}

func BF(ms []M) bruteforce.L {
	runtime.MemProfileRate = 0
	defer func() { runtime.MemProfileRate = 512 * 1024 }()

	l := bruteforce.New()
	for _, f := range ms {
		f(l)
	}
	return l
}

// grid generates a complete K-dimensional grid.
func grid(min, max int, k vector.D) []vector.V {
	if k == 1 {
		var vs []vector.V
		for i := min; i < max; i++ {
			vs = append(vs, []float64{float64(i)})
		}
		return vs
	}

	var vs []vector.V
	for _, j := range grid(min, max, k-1) {
		for i := min; i < max; i++ {
			v := []float64{float64(i)}
			v = append(v, j...)
			vs = append(vs, v)
		}
	}
	return vs
}

// InsertRandom generates a K-dimensional grid and returns a 20% dense grid
// of inserts.
//
// We generate this grid by first generating a complete grid of unit tiles, and
// then shuffling them in-place. We use the first N values from this list.
func InsertRandom(ids []id.ID, k vector.D, n int) []M {
	runtime.MemProfileRate = 0
	defer func() { runtime.MemProfileRate = 512 * 1024 }()

	for i := 0; i < n; i++ {
	}
	xs := map[id.ID]bool{}
	for _, x := range ids {
		xs[x] = true
	}

	tiles := grid(0, int(math.Pow(5*float64(n), 1.0/float64(k))), k)
	rand.Shuffle(len(tiles), func(i, j int) { tiles[i], tiles[j] = tiles[j], tiles[i] })

	fs := make([]M, 0, n)
	// Only use the first n tiles.
	for i := 0; i < n; i++ {
		x := allocateID(xs)
		xs[x] = true

		vmin := tiles[i]
		vmax := make([]float64, k)
		for i := vector.D(0); i < k; i++ {
			vmax[i] = vmin[i] + 1
		}
		fs = append(fs,
			func(c container.C) error {
				return c.Insert(x, *hyperrectangle.New(vmin, vmax))
			},
		)
	}
	return fs

}

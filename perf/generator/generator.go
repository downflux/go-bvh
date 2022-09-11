package generator

import (
	"math/rand"
	"runtime"

	"github.com/downflux/go-bvh/bvh"
	"github.com/downflux/go-bvh/container"
	"github.com/downflux/go-bvh/container/bruteforce"
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type M func(c container.C) error

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

func Generate(o O, n int) []M {
	runtime.MemProfileRate = 0
	defer func() { runtime.MemProfileRate = 512 * 1024 }()

	ids := map[id.ID]bool{}
	for _, x := range o.IDs {
		ids[x] = true
	}

	tiles := make([]int, 0, n)
	for i := 0; i < n; i++ {
		tiles = append(tiles, i)
	}
	rand.Shuffle(len(tiles), func(i, j int) { tiles[i], tiles[j] = tiles[j], tiles[i] })

	fs := make([]M, 0, n)
	for i := 0; i < n; i++ {
		x := allocateID(ids)
		ids[x] = true

		vmin := make([]float64, o.K)
		vmax := make([]float64, o.K)
		for k := vector.D(0); k < o.K; k++ {
			vmax[k] = 1
		}
		vmin[0] = float64(tiles[i])
		vmax[0] = float64(tiles[i]) + 1

		fs = append(fs,
			func(c container.C) error {
				return c.Insert(x, *hyperrectangle.New(vmin, vmax))
			},
		)
	}
	return fs
}

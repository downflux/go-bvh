package generator

import (
	"log"
	"math/rand"

	"github.com/downflux/go-bvh/bvh"
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type O struct {
	// Insert is the insert node weight.
	Insert float64
	// Remove is the remove node weight.
	Remove float64

	// N is the number of opts to call.
	N int

	K vector.D

	Size   uint
	Logger *log.Logger
}

type L struct {
	insert float64
	remove float64
	n      int
	k      vector.D

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
		for k := vector.D(0); k < l.k; k++ {
			vmax[k] = 1
		}
		vmin[0] = float64(tiles[i])
		vmax[0] = float64(tiles[i]) + 1

		fs = append(fs, func() { l.bvh.Insert(j, *hyperrectangle.New(vmin, vmax)) })
	}
	return fs
}

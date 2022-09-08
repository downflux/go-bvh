package perf

import (
	"math"
	"math/rand"

	"github.com/downflux/go-bvh/bvh"
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

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

	K    int
	Size uint
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
		bvh:    bvh.New(o.Size),
		ids:    map[id.ID]bool{},
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
func (l *L) Apply(min, max float64) *bvh.BVH {
	for _, f := range l.Generate(min, max) {
		f()
	}
	return l.BVH()
}

func (l *L) Generate(min, max float64) []func() {
	fs := make([]func(), 0, l.n)
	for i := 0; i < l.n; i++ {
		if rand.Float64() <= l.insert {
			var j id.ID
			for j = id.ID(rand.Uint64()); l.ids[j]; j = id.ID(rand.Uint64()) {
			}
			l.ids[j] = true
			fs = append(fs, func() { l.bvh.Insert(j, rr(min, max, l.k)) })
		} else {
			if len(l.ids) > 0 {
				j := l.IDs()[rand.Intn(len(l.ids))]
				l.ids[j] = false
				fs = append(fs, func() { l.bvh.Remove(j) })
			}
		}
	}
	return fs
}

package hyperrectangle

import (
	"math"

	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

type BoundOpt[T point.P] struct {
	Data []T
	K    vector.D
	Size int
}

func Bound[T point.P](o BoundOpt[T]) hyperrectangle.R {
	min := make([]float64, o.K)
	max := make([]float64, o.K)

	for i := vector.D(0); i < o.K; i++ {
		min[i] = math.Inf(0)
		max[i] = math.Inf(-1)
	}

	b := *hyperrectangle.New(vector.V(min), vector.V(max))

	if len(o.Data) < o.Size {
		for _, p := range o.Data {
			b = Union(b, p.Bound())
		}
	} else {
		l := make(chan hyperrectangle.R)
		r := make(chan hyperrectangle.R)
		go func(ch chan<- hyperrectangle.R) {
			ch <- Bound[T](BoundOpt[T]{
				Data: o.Data[0 : len(o.Data)/2],
				K:    o.K,
				Size: o.Size,
			})
			close(ch)
		}(l)
		go func(ch chan<- hyperrectangle.R) {
			ch <- Bound[T](BoundOpt[T]{
				Data: o.Data[len(o.Data)/2 : len(o.Data)-1],
				K:    o.K,
				Size: o.Size,
			})
			close(ch)
		}(r)
		b = Union(<-l, <-r)
	}

	return b
}

func Collision(r hyperrectangle.R, s hyperrectangle.R) bool {
	if r.Min().Dimension() != s.Min().Dimension() {
		panic("mismatching vector dimensions")
	}

	for i := vector.D(0); i < r.Min().Dimension(); i++ {
		l := math.Max(r.Min().X(i), s.Min().X(i))
		u := math.Min(r.Max().X(i), s.Max().X(i))
		if l < u {
			return false
		}
	}
	return true
}

func Union(r hyperrectangle.R, s hyperrectangle.R) hyperrectangle.R {
	if r.Min().Dimension() != s.Min().Dimension() {
		panic("mismatching vector dimensions")
	}

	min := make([]float64, r.Min().Dimension())
	max := make([]float64, r.Min().Dimension())

	for i := vector.D(0); i < r.Min().Dimension(); i++ {
		min[i] = math.Min(r.Min().X(i), s.Min().X(i))
		max[i] = math.Max(r.Max().X(i), s.Max().X(i))
	}

	return *hyperrectangle.New(min, max)
}

func Heuristic(r hyperrectangle.R) float64 {
	h := 1.0
	d := r.D()
	for i := vector.D(0); i < d.Dimension(); i++ {
		h *= d.X(i)
	}
	return h
}

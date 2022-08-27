package hyperrectangle

import (
	"math"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

const size = 128

func AABB(rs []hyperrectangle.R) hyperrectangle.R {
	if len(rs) == 0 {
		return hyperrectangle.R{}
	}
	if len(rs) == 1 {
		return rs[0]
	}
	if len(rs) == 2 {
		return Union(rs[0], rs[1])
	}

	var b hyperrectangle.R
	if len(rs) <= size {
		b = rs[0]
		for _, r := range rs[1:] {
			b = Union(b, r)
		}
	} else {
		l := make(chan hyperrectangle.R)
		r := make(chan hyperrectangle.R)
		go func(ch chan<- hyperrectangle.R) {
			ch <- AABB(rs[:len(rs)/2])
			close(ch)
		}(l)
		go func(ch chan<- hyperrectangle.R) {
			ch <- AABB(rs[len(rs)/2:])
			close(ch)
		}(r)
		b = Union(<-l, <-r)
	}

	return b
}

// Contains checks if the input rectangle r fully encloses s.
func Contains(r hyperrectangle.R, s hyperrectangle.R) bool {
	if r.Min().Dimension() != s.Min().Dimension() {
		panic("mismatching vector dimensions")
	}

	for i := vector.D(0); i < r.Min().Dimension(); i++ {
		if s.Min().X(i) <= r.Min().X(i) || s.Max().X(i) >= r.Max().X(i) {
			return false
		}
	}
	return true
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

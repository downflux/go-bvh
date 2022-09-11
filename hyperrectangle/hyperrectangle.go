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
		b = *hyperrectangle.New(
			make([]float64, rs[0].Min().Dimension()),
			make([]float64, rs[0].Min().Dimension()),
		)
		AABBBuf(rs, b)
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

func AABBBuf(rs []hyperrectangle.R, buf hyperrectangle.R) {
	if len(rs) == 0 {
		return
	}
	copy(buf.Min(), rs[0].Min())
	copy(buf.Max(), rs[0].Max())

	for _, r := range rs {
		UnionBuf(r, buf, buf)
	}
}

// Contains checks if the input rectangle r fully encloses s.
//
// We are treating r as a closed interval.
func Contains(r hyperrectangle.R, s hyperrectangle.R) bool {
	if r.Min().Dimension() != s.Min().Dimension() {
		panic("mismatching vector dimensions")
	}

	for i := vector.D(0); i < r.Min().Dimension(); i++ {
		if s.Min().X(i) < r.Min().X(i) || s.Max().X(i) > r.Max().X(i) {
			return false
		}
	}
	return true
}

func Disjoint(r hyperrectangle.R, s hyperrectangle.R) bool {
	if r.Min().Dimension() != s.Min().Dimension() {
		panic("mismatching vector dimensions")
	}

	for i := vector.D(0); i < r.Min().Dimension(); i++ {
		if (r.Min().X(i) < s.Min().X(i) && r.Max().X(i) < s.Min().X(i)) || (s.Min().X(i) < r.Min().X(i) && s.Max().X(i) < r.Min().X(i)) {
			return true
		}
	}
	return false
}

func UnionBuf(r hyperrectangle.R, s hyperrectangle.R, buf hyperrectangle.R) {
	if r.Min().Dimension() != s.Min().Dimension() || r.Min().Dimension() != buf.Min().Dimension() {
		panic("mismatching vector dimensions")
	}

	for i := vector.D(0); i < r.Min().Dimension(); i++ {
		buf.Min()[i] = math.Min(r.Min().X(i), s.Min().X(i))
		buf.Max()[i] = math.Max(r.Max().X(i), s.Max().X(i))
	}
}

func Union(r hyperrectangle.R, s hyperrectangle.R) hyperrectangle.R {
	if r.Min().Dimension() != s.Min().Dimension() {
		panic("mismatching vector dimensions")
	}

	buf := *hyperrectangle.New(
		make([]float64, r.Min().Dimension()),
		make([]float64, r.Min().Dimension()),
	)

	UnionBuf(r, s, buf)
	return buf
}

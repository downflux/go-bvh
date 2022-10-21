package hyperrectangle

import (
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

const size = 4

func AABB(rs []hyperrectangle.R) hyperrectangle.R {
	if len(rs) == 0 {
		return hyperrectangle.R{}
	}
	if len(rs) == 1 {
		return rs[0]
	}
	if len(rs) == 2 {
		return hyperrectangle.Union(rs[0], rs[1])
	}

	var b hyperrectangle.R
	if len(rs) <= size {
		b = *hyperrectangle.New(
			make([]float64, rs[0].Min().Dimension()),
			make([]float64, rs[0].Min().Dimension()),
		)
		AABBBuf(rs, b.M())
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
		b = hyperrectangle.Union(<-l, <-r)
	}

	return b
}

func AABBBuf(rs []hyperrectangle.R, buf hyperrectangle.M) {
	if len(rs) == 0 {
		return
	}
	buf.Copy(rs[0])

	for _, r := range rs[1:] {
		buf.Union(r)
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

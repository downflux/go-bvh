package util

import (
	"math"
	"math/rand"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

func RN(min, max float64) float64 { return rand.Float64()*(max-min) + min }
func RV(min, max float64, k vector.D) vector.V {
	vs := []float64{}
	for i := vector.D(0); i < k; i++ {
		vs = append(vs, RN(min, max))
	}
	return vector.V(vs)
}
func RR(min, max float64, k vector.D) hyperrectangle.R {
	a := RV(min, max, k)
	b := RV(min, max, k)

	vmin := make([]float64, k)
	vmax := make([]float64, k)

	for i := vector.D(0); i < k; i++ {
		vmin[i] = math.Min(a[i], b[i])
		vmax[i] = math.Max(a[i], b[i])
	}
	return *hyperrectangle.New(vmin, vmax)
}

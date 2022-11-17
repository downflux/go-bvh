package perf

import (
	"math"
	"math/rand"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

func rn(min float64, max float64) float64 { return rand.Float64() * (max - min) }

func GenerateAABB(k vector.D, min float64, max float64) hyperrectangle.R {
	aabb := hyperrectangle.New(
		vector.V(make([]float64, k)),
		vector.V(make([]float64, k)),
	).M()
	for i := vector.D(0); i < k; i++ {
		a, b := rn(min, max), rn(min, max)
		aabb.Min().SetX(i, math.Min(a, b))
		aabb.Max().SetX(i, math.Max(a, b))
	}
	return aabb.R()
}

func GenerateObjects(n int, k vector.D, min float64, max float64) map[id.ID]hyperrectangle.R {
	data := make(map[id.ID]hyperrectangle.R, n)
	for i := 0; i < n; i++ {
		data[id.ID(i)] = GenerateAABB(k, min, max)
	}
	return data
}

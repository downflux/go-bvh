package perf

import (
	"math"
	"math/rand"
	"runtime"

	"github.com/downflux/go-bvh/x/container"
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

func GenerateRandomTiles(n int, k vector.D) map[id.ID]hyperrectangle.R {
	runtime.MemProfileRate = 0
	defer func() { runtime.MemProfileRate = 512 * 1024 }()

	tiles := grid(0, int(math.Ceil(math.Pow(5*float64(n), 1.0/float64(k)))), k)
	rand.Shuffle(len(tiles), func(i, j int) { tiles[i], tiles[j] = tiles[j], tiles[i] })

	aabbs := make(map[id.ID]hyperrectangle.R, n)
	for i := 0; i < n; i++ {
		max := vector.M(make([]float64, k))
		for j := vector.D(0); j < k; j++ {
			max.SetX(j, tiles[i].X(j)+1)
		}
		aabbs[id.ID(i)] = *hyperrectangle.New(tiles[i], max.V())
	}
	return aabbs
}

func GenerateRandomBoxes(n int, k vector.D, min float64, max float64) map[id.ID]hyperrectangle.R {
	runtime.MemProfileRate = 0
	defer func() { runtime.MemProfileRate = 512 * 1024 }()

	data := make(map[id.ID]hyperrectangle.R, n)
	for i := 0; i < n; i++ {
		data[id.ID(i)] = GenerateAABB(k, min, max)
	}
	return data
}

type F func(c container.C) error

func GenerateInsertLoad(n int, offset int, k vector.D) []F {
	ops := make([]F, 0, n)
	for x, aabb := range GenerateRandomTiles(n, k) {
		x, aabb := x, aabb
		ops = append(ops, func(c container.C) error {
			return c.Insert(x+id.ID(offset), aabb)
		})
	}
	return ops
}

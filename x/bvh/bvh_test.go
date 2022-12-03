package bvh

import (
	"fmt"
	"testing"

	"github.com/downflux/go-bvh/x/container/bruteforce"
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/perf"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBroadPhaseConformance(t *testing.T) {
	const k = 3
	type config struct {
		size      int
		tolerance float64
		data      map[id.ID]hyperrectangle.R
		q         hyperrectangle.R
	}

	configs := []config{
		{
			size:      4,
			tolerance: 1.05,
			data:      map[id.ID]hyperrectangle.R{},
			q: *hyperrectangle.New(
				vector.V([]float64{10, 10, 10}),
				vector.V([]float64{90, 90, 90}),
			),
		},
		{
			size:      1,
			tolerance: 1.05,
			data: map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(
					vector.V([]float64{0, 0, 0}),
					vector.V([]float64{1, 1, 1}),
				),
				101: *hyperrectangle.New(
					vector.V([]float64{10, 10, 10}),
					vector.V([]float64{11, 11, 11}),
				),
			},
			q: *hyperrectangle.New(
				vector.V([]float64{1, 1, 1}),
				vector.V([]float64{10, 10, 10}),
			),
		},
		{
			size:      4,
			tolerance: 1.05,
			data: map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(
					vector.V([]float64{0, 0, 0}),
					vector.V([]float64{1, 1, 1}),
				),
				101: *hyperrectangle.New(
					vector.V([]float64{10, 10, 10}),
					vector.V([]float64{11, 11, 11}),
				),
				102: *hyperrectangle.New(
					vector.V([]float64{9, 9, 9}),
					vector.V([]float64{11, 11, 11}),
				),
				103: *hyperrectangle.New(
					vector.V([]float64{30, 30, 30}),
					vector.V([]float64{40, 40, 40}),
				),
				104: *hyperrectangle.New(
					vector.V([]float64{100, 100, 100}),
					vector.V([]float64{101, 101, 101}),
				),
				105: *hyperrectangle.New(
					vector.V([]float64{90.01, 90.01, 90.01}),
					vector.V([]float64{95, 95, 95}),
				),
				106: *hyperrectangle.New(
					vector.V([]float64{0, 0, 0}),
					vector.V([]float64{100, 100, 100}),
				),
			},
			q: *hyperrectangle.New(
				vector.V([]float64{10, 10, 10}),
				vector.V([]float64{90, 90, 90}),
			),
		},
		{
			size:      1,
			tolerance: 1.05,
			data:      perf.GenerateRandomBoxes(10000, k, 100, 200),
			q:         perf.GenerateAABB(k, 50, 70),
		},
		{
			size:      4,
			tolerance: 1.05,
			data:      perf.GenerateRandomBoxes(1000, k, 100, 200),
			q:         perf.GenerateAABB(k, 50, 70),
		},
	}

	for _, c := range configs {
		t.Run(fmt.Sprintf("BVH/K=%v/LeafSize=%v/Tolerance=%v/N=%v", k, c.size, c.tolerance, len(c.data)), func(t *testing.T) {
			tbf := bruteforce.New()
			tbvh := New(O{
				K:         k,
				LeafSize:  c.size,
				Tolerance: c.tolerance,
			})

			for x, h := range c.data {
				tbf.Insert(x, h)
				tbvh.Insert(x, h)
			}

			want := tbf.BroadPhase(c.q)
			got := tbvh.BroadPhase(c.q)

			if diff := cmp.Diff(
				want, got,
				cmpopts.SortSlices(func(a, b id.ID) bool { return a < b }),
			); diff != "" {
				t.Errorf("BroadPhase() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

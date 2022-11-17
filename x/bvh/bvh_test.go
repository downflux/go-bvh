package bvh

import (
	"fmt"
	"testing"

	"github.com/downflux/go-bvh/x/container/bruteforce"
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/perf"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBroadPhaseConformance(t *testing.T) {
	type config struct {
		k         vector.D
		size      int
		tolerance float64
		data      map[id.ID]hyperrectangle.R
		q         hyperrectangle.R
	}

	configs := []config{
		{
			k:         2,
			size:      4,
			tolerance: 1.05,
			data:      map[id.ID]hyperrectangle.R{},
			q: *hyperrectangle.New(
				vector.V([]float64{10, 10}),
				vector.V([]float64{90, 90}),
			),
		},
		{
			k:         2,
			size:      1,
			tolerance: 1.05,
			data: map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(
					vector.V([]float64{0, 0}),
					vector.V([]float64{1, 1}),
				),
				101: *hyperrectangle.New(
					vector.V([]float64{10, 10}),
					vector.V([]float64{11, 11}),
				),
			},
			q: *hyperrectangle.New(
				vector.V([]float64{1, 1}),
				vector.V([]float64{10, 10}),
			),
		},
		{
			k:         2,
			size:      4,
			tolerance: 1.05,
			data: map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(
					vector.V([]float64{0, 0}),
					vector.V([]float64{1, 1}),
				),
				101: *hyperrectangle.New(
					vector.V([]float64{10, 10}),
					vector.V([]float64{11, 11}),
				),
				102: *hyperrectangle.New(
					vector.V([]float64{9, 9}),
					vector.V([]float64{11, 11}),
				),
				103: *hyperrectangle.New(
					vector.V([]float64{30, 30}),
					vector.V([]float64{40, 40}),
				),
				104: *hyperrectangle.New(
					vector.V([]float64{100, 100}),
					vector.V([]float64{101, 101}),
				),
				105: *hyperrectangle.New(
					vector.V([]float64{90.01, 90.01}),
					vector.V([]float64{95, 95}),
				),
				106: *hyperrectangle.New(
					vector.V([]float64{0, 0}),
					vector.V([]float64{100, 100}),
				),
			},
			q: *hyperrectangle.New(
				vector.V([]float64{10, 10}),
				vector.V([]float64{90, 90}),
			),
		},
		{
			k:         2,
			size:      1,
			tolerance: 1.05,
			data:      perf.GenerateObjects(10000, 2, 100, 200),
			q:         perf.GenerateAABB(2, 50, 70),
		},
		{
			k:         2,
			size:      4,
			tolerance: 1.05,
			data:      perf.GenerateObjects(1000, 2, 100, 200),
			q:         perf.GenerateAABB(2, 50, 70),
		},
	}

	for _, c := range configs {
		name := fmt.Sprintf("K=%v/LeafSize=%v/Tolerance=%v/N=%v", c.k, c.size, c.tolerance, len(c.data))
		t.Run(name, func(t *testing.T) {
			tbf := bruteforce.New()
			tbvh := New(O{
				K:         c.k,
				LeafSize:  c.size,
				Tolerance: c.tolerance,
			})

			for x, h := range c.data {
				tbf.Insert(x, h)
				tbvh.Insert(x, h)
			}

			q := *hyperrectangle.New(
				vector.V([]float64{10, 10}),
				vector.V([]float64{90, 90}),
			)

			want := tbf.BroadPhase(q)
			got := tbvh.BroadPhase(q)

			if diff := cmp.Diff(
				want, got,
				cmpopts.SortSlices(func(a, b id.ID) bool { return a < b }),
			); diff != "" {
				t.Errorf("BroadPhase() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}

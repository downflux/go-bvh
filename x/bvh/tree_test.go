package bvh

import (
	"fmt"
	"testing"

	"github.com/downflux/go-bvh/x/container/bruteforce"
	"github.com/downflux/go-bvh/x/id"
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
			size:      4,
			tolerance: 1.05,
			data: map[id.ID]hyperrectangle.R{
				0: *hyperrectangle.New(
					vector.V([]float64{0, 0}),
					vector.V([]float64{1, 1}),
				),
				1: *hyperrectangle.New(
					vector.V([]float64{10, 10}),
					vector.V([]float64{11, 11}),
				),
				2: *hyperrectangle.New(
					vector.V([]float64{9, 9}),
					vector.V([]float64{11, 11}),
				),
				3: *hyperrectangle.New(
					vector.V([]float64{30, 30}),
					vector.V([]float64{40, 40}),
				),
				4: *hyperrectangle.New(
					vector.V([]float64{100, 100}),
					vector.V([]float64{101, 101}),
				),
				5: *hyperrectangle.New(
					vector.V([]float64{0, 0}),
					vector.V([]float64{100, 100}),
				),
			},
			q: *hyperrectangle.New(
				vector.V([]float64{10, 10}),
				vector.V([]float64{90, 90}),
			),
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

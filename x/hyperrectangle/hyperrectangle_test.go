package hyperrectangle

import (
	"testing"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"
)

func TestAABB(t *testing.T) {
	type config struct {
		name string
		rs   []hyperrectangle.R
		want hyperrectangle.R
	}

	configs := []config{
		{
			name: "Single",
			rs: []hyperrectangle.R{
				*hyperrectangle.New(
					[]float64{0, 0},
					[]float64{100, 100},
				),
			},
			want: *hyperrectangle.New(
				[]float64{0, 0},
				[]float64{100, 100},
			),
		},
		{
			name: "Double",
			rs: []hyperrectangle.R{
				*hyperrectangle.New(
					[]float64{0, 0},
					[]float64{50, 50},
				),
				*hyperrectangle.New(
					[]float64{99, 99},
					[]float64{100, 100},
				),
			},
			want: *hyperrectangle.New(
				[]float64{0, 0},
				[]float64{100, 100},
			),
		},
		func() config {
			n := float64(2*size + 1)
			var rs []hyperrectangle.R
			for i := 0.0; i <= n; i++ {
				rs = append(rs, *hyperrectangle.New(
					[]float64{0},
					[]float64{i},
				))
			}
			return config{
				name: "Concurrent",
				rs:   rs,
				want: *hyperrectangle.New(
					[]float64{0},
					[]float64{n},
				),
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if diff := cmp.Diff(
				c.want,
				AABB(c.rs),
				cmp.AllowUnexported(hyperrectangle.R{})); diff != "" {
				t.Errorf("AABB() mismatch (-want +got):\n%v", diff)
			}
		})
	}
}
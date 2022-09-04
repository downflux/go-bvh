package hyperrectangle

import (
	"testing"

	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/google/go-cmp/cmp"
)

func TestV(t *testing.T) {
	type config struct {
		name string
		r    hyperrectangle.R
		want float64
	}

	configs := []config{
		{
			name: "1D",
			r:    *hyperrectangle.New([]float64{10}, []float64{100}),
			want: 90,
		},
		{
			name: "2D",
			r:    *hyperrectangle.New([]float64{10, 10}, []float64{100, 100}),
			want: 8100,
		},
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			if got := V(c.r); got != c.want {
				t.Errorf("V() = %v, want = %v", got, c.want)
			}
		})
	}
}

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

func TestDisjoint(t *testing.T) {
	type config struct {
		name string
		r    hyperrectangle.R
		s    hyperrectangle.R
		want bool
	}

	configs := []config{
		{
			name: "Simple/Disjoint",
			r: *hyperrectangle.New(
				[]float64{0},
				[]float64{10},
			),
			s: *hyperrectangle.New(
				[]float64{11},
				[]float64{20},
			),
			want: true,
		},
		{
			name: "Simple/Overlap",
			r: *hyperrectangle.New(
				[]float64{0},
				[]float64{10},
			),
			s: *hyperrectangle.New(
				[]float64{9},
				[]float64{20},
			),
			want: false,
		},
		{
			name: "Simple/Disjoint/Commutative",
			r: *hyperrectangle.New(
				[]float64{11},
				[]float64{20},
			),
			s: *hyperrectangle.New(
				[]float64{0},
				[]float64{10},
			),
			want: true,
		},
		{
			name: "2D/Disjoint",
			r: *hyperrectangle.New(
				[]float64{0, 0},
				[]float64{10, 10},
			),
			s: *hyperrectangle.New(
				[]float64{5, 11},
				[]float64{20, 20},
			),
			want: true,
		},
		{
			name: "2D/Overlap",
			r: *hyperrectangle.New(
				[]float64{0, 0},
				[]float64{10, 10},
			),
			s: *hyperrectangle.New(
				[]float64{5, 5},
				[]float64{20, 20},
			),
			want: false,
		},
	}

	for _, c := range configs {
		if got := Disjoint(c.r, c.s); got != c.want {
			t.Errorf("Disjoint() = %v, want = %v", got, c.want)
		}
	}
}

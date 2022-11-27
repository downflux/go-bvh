package perf

import (
	"fmt"

	"github.com/downflux/go-geometry/nd/vector"
)

type PerfTestSize int

const (
	SizeUnknown PerfTestSize = iota
	SizeUnit
	SizeSmall
	SizeLarge
)

func (s *PerfTestSize) String() string {
	return map[PerfTestSize]string{
		SizeLarge: "large",
		SizeSmall: "small",
		SizeUnit:  "unit",
	}[*s]
}

func (s *PerfTestSize) Set(v string) error {
	size, ok := map[string]PerfTestSize{
		"large": SizeLarge,
		"small": SizeSmall,
		"unit":  SizeUnit,
	}[v]
	if !ok {
		return fmt.Errorf("invalid test size value: %v", v)
	}
	*s = size
	return nil
}

func (s PerfTestSize) N() []int {
	return map[PerfTestSize][]int{
		SizeLarge: []int{1e3, 1e4, 1e5},
		SizeSmall: []int{1e3, 1e4},
		SizeUnit:  []int{1e3},
	}[s]
}

func (s PerfTestSize) F() []float64 {
	return map[PerfTestSize][]float64{
		SizeLarge: []float64{0.05},
		SizeSmall: []float64{0.05},
		SizeUnit:  []float64{0.05},
	}[s]
}

func (s PerfTestSize) LeafSize() []uint {
	return map[PerfTestSize][]uint{
		SizeLarge: []uint{2, 4, 8, 16, 32, 64},
		SizeSmall: []uint{2, 4, 8, 16},
		SizeUnit:  []uint{1, 2},
	}[s]
}

func (s PerfTestSize) K() []vector.D {
	return map[PerfTestSize][]vector.D{
		SizeLarge: []vector.D{3},
		SizeSmall: []vector.D{3},
		SizeUnit:  []vector.D{3},
	}[s]
}

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
		SizeSmall: []int{1e3, 1e4, 1e5},
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
		SizeLarge: []uint{1, 16, 256, 1024},
		SizeSmall: []uint{1, 4, 16, 64},
		SizeUnit:  []uint{1, 2},
	}[s]
}

func (s PerfTestSize) K() []vector.D {
	return map[PerfTestSize][]vector.D{
		SizeLarge: []vector.D{2, 3, 10},
		SizeSmall: []vector.D{2},
		SizeUnit:  []vector.D{2},
	}[s]
}

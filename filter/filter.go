package filter

import (
	"github.com/downflux/go-bvh/point"
)

type F[T point.P] func(p T) bool

package filter

import (
	"github.com/downflux/go-bvh/point"
)

type F[T point.RO] func(p T) bool

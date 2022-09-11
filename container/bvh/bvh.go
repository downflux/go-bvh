package bvh

import (
	"github.com/downflux/go-bvh/bvh"
	"github.com/downflux/go-bvh/container"
)

var (
	_ container.C = &bvh.BVH{}
)

package bvh

import (
	"github.com/downflux/go-bvh/x/bvh"
	"github.com/downflux/go-bvh/x/container"
)

var (
	_ container.C = &bvh.T{}
)

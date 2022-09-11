package bvh

import (
	"github.com/downflux/go-bvh/bvh"
	"github.com/downflux/go-bvh/container"
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

var (
	_ container.C = &BVH{}
)

type BVH bvh.BVH

type O struct {
	Size uint
}

func New(o O) *BVH { return (*BVH)(bvh.New(bvh.O{Size: o.Size})) }

func (t *BVH) Insert(x id.ID, aabb hyperrectangle.R) error { return (*bvh.BVH)(t).Insert(x, aabb) }
func (t *BVH) Remove(x id.ID) error                        { return (*bvh.BVH)(t).Remove(x) }
func (t *BVH) Update(x id.ID, q hyperrectangle.R, aabb hyperrectangle.R) error {
	return (*bvh.BVH)(t).Update(x, q, aabb)
}
func (t *BVH) BroadPhase(q hyperrectangle.R) []id.ID { return bvh.BroadPhase((*bvh.BVH)(t), q) }

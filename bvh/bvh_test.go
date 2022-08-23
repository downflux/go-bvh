package bvh

import (
	"testing"

	"github.com/downflux/go-bvh/internal/node/util"
	"github.com/downflux/go-bvh/point"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

type n struct {
	id point.ID
}

func (n n) Bound() hyperrectangle.R { return hyperrectangle.R{} }
func (n n) ID() point.ID            { return n.id }

func Equal(a, b BVH[*n]) bool {
	if !util.Equal(a.allocation, a.root, b.allocation, b.root) {
		return false
	}

	for pid, t := range a.data {
		if t != b.data[pid] {
			return false
		}
	}

	for pid, aid := range a.lookup {
		if !util.Equal(a.allocation, aid, b.allocation, b.lookup[pid]) {
			return false
		}
	}
	return true
}

func TestRemove(t *testing.T) {

}

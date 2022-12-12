package container

import (
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

type C interface {
	IDs() []id.ID
	Insert(x id.ID, aabb hyperrectangle.R) error
	BroadPhase(q hyperrectangle.R) []id.ID
	Remove(x id.ID) error
	Update(x id.ID, aabb hyperrectangle.R) error
}

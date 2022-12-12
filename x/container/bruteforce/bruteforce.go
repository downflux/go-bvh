package bruteforce

import (
	"fmt"

	"github.com/downflux/go-bvh/x/container"
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
)

var (
	_ container.C = L{}
)

type L map[id.ID]hyperrectangle.R

func New() L { return L(map[id.ID]hyperrectangle.R{}) }

func (l L) IDs() []id.ID {
	ids := make([]id.ID, len(l))
	for x := range l {
		ids = append(ids, x)
	}
	return ids
}

func (l L) Insert(x id.ID, aabb hyperrectangle.R) error {
	if _, ok := l[x]; ok {
		return fmt.Errorf("cannot insert a node with duplicate ID %v", x)
	}
	l[x] = aabb
	return nil
}

func (l L) Remove(x id.ID) error {
	if _, ok := l[x]; !ok {
		return fmt.Errorf("cannot remove a non-existent object with ID %v", x)
	}
	delete(l, x)
	return nil
}

func (l L) Update(x id.ID, aabb hyperrectangle.R) error {
	if err := l.Remove(x); err != nil {
		return fmt.Errorf("cannot update object: %v", err)
	}
	if err := l.Insert(x, aabb); err != nil {
		return fmt.Errorf("cannot update object: %v", err)
	}
	return nil
}

func (l L) BroadPhase(q hyperrectangle.R) []id.ID {
	ids := make([]id.ID, 0, 128)
	for x, aabb := range l {
		if !hyperrectangle.Disjoint(q, aabb) {
			ids = append(ids, x)
		}
	}
	return ids
}

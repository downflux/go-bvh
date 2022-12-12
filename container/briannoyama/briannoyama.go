package briannoyama

import (
	"fmt"

	"github.com/briannoyama/bvh/rect"
	"github.com/downflux/go-bvh/container"
	"github.com/downflux/go-bvh/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

var (
	_ container.C = &BVH{}
)

type BVH struct {
	lookup  map[id.ID]*rect.Orthotope
	reverse map[*rect.Orthotope]id.ID

	bvh *rect.BVol
}

func New() *BVH {
	return &BVH{
		bvh:     &rect.BVol{},
		lookup:  map[id.ID]*rect.Orthotope{},
		reverse: map[*rect.Orthotope]id.ID{},
	}
}

func orth(aabb hyperrectangle.R) *rect.Orthotope {
	k := aabb.Min().Dimension()
	if k != 3 {
		panic(fmt.Sprintf("unsupported vector length %v", k))
	}

	vmin := [3]int32{}
	vdelta := [3]int32{}
	for i := vector.D(0); i < k; i++ {
		vmin[i] = int32(aabb.Min().X(i))
		vdelta[i] = int32(aabb.Max().X(i) - aabb.Min().X(i))
	}
	return &rect.Orthotope{
		Point: vmin, Delta: vdelta,
	}
}

func (bvh *BVH) Insert(x id.ID, aabb hyperrectangle.R) error {
	if bvh.lookup[x] != nil {
		return fmt.Errorf("cannot insert duplicate ID %v", x)
	}

	bvh.lookup[x] = orth(aabb)
	bvh.reverse[bvh.lookup[x]] = x
	bvh.bvh.Add(bvh.lookup[x])
	return nil
}

func (bvh *BVH) IDs() []id.ID {
	xs := make([]id.ID, 0, len(bvh.lookup))
	for x := range bvh.lookup {
		xs = append(xs, x)
	}
	return xs
}

func (bvh *BVH) SAH() float64 { return bvh.bvh.SAH() }

func (bvh *BVH) Remove(x id.ID) error {
	if bvh.lookup[x] == nil {
		return fmt.Errorf("cannot remove non-existent ID %v", x)
	}

	bvh.bvh.Remove(bvh.lookup[x])
	delete(bvh.reverse, bvh.lookup[x])
	delete(bvh.lookup, x)
	return nil
}

func (bvh *BVH) Update(x id.ID, aabb hyperrectangle.R) error {
	if err := bvh.Remove(x); err != nil {
		return fmt.Errorf("cannot update ID %v: %v", x, err)
	}
	if err := bvh.Insert(x, aabb); err != nil {
		return fmt.Errorf("cannot update ID %v: %v", x, err)
	}
	return nil
}

func (bvh *BVH) BroadPhase(q hyperrectangle.R) []id.ID {
	var xs []id.ID
	iter := bvh.bvh.Iterator()
	p := orth(q)
	for r := iter.Query(p); r != nil; r = iter.Query(p) {
		xs = append(xs, bvh.reverse[r])
	}
	return xs
}

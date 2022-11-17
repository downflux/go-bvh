package dhconnelly

import (
	"fmt"

	"github.com/dhconnelly/rtreego"
	"github.com/downflux/go-bvh/x/container"
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

var (
	_ container.C = &BVH{}
)

type r rtreego.Rect

func (rect *r) Bounds() *rtreego.Rect { return (*rtreego.Rect)(rect) }

type BVH struct {
	lookup  map[id.ID]*r
	reverse map[*r]id.ID

	bvh *rtreego.Rtree
}

func rect(aabb hyperrectangle.R) *rtreego.Rect {
	r, err := rtreego.NewRectFromPoints([]float64(aabb.Min()), []float64(aabb.Max()))
	if err != nil {
		panic(fmt.Sprintf("cannot construct rectangle: %v", err))
	}

	return r
}

type O struct {
	K         vector.D
	MinBranch int
	MaxBranch int
}

func New(o O) *BVH {
	return &BVH{
		bvh:     rtreego.NewTree(int(o.K), o.MinBranch, o.MaxBranch),
		lookup:  map[id.ID]*r{},
		reverse: map[*r]id.ID{},
	}
}

func (bvh *BVH) IDs() []id.ID {
	xs := make([]id.ID, 0, len(bvh.lookup))
	for x := range bvh.lookup {
		xs = append(xs, x)
	}
	return xs
}

func (bvh *BVH) Insert(x id.ID, aabb hyperrectangle.R) error {
	if bvh.lookup[x] != nil {
		return fmt.Errorf("cannot insert duplicate ID %v", x)
	}

	bvh.lookup[x] = (*r)(rect(aabb))
	bvh.reverse[bvh.lookup[x]] = x
	bvh.bvh.Insert(bvh.lookup[x])
	return nil
}

func (bvh *BVH) Remove(x id.ID) error {
	if bvh.lookup[x] == nil {
		return fmt.Errorf("cannot remove non-existent ID %v", x)
	}

	bvh.bvh.Delete(bvh.lookup[x])
	delete(bvh.reverse, bvh.lookup[x])
	delete(bvh.lookup, x)
	return nil
}

func (bvh *BVH) Update(x id.ID, q hyperrectangle.R, aabb hyperrectangle.R) error {
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
	for _, aabb := range bvh.bvh.SearchIntersect(rect(q)) {
		xs = append(xs, bvh.reverse[aabb.(*r)])
	}
	return xs
}

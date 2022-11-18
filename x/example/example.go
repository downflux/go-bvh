package main

import (
	"fmt"

	"github.com/downflux/go-bvh/x/bvh"
	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"
)

func main() {
	t := bvh.New(bvh.O{
		LeafSize:  1,
		K:         vector.D(2),
		Tolerance: 1,
	})

	data := map[id.ID]hyperrectangle.R{
		100: *hyperrectangle.New(vector.V{0, 0}, vector.V{1, 1}),
		101: *hyperrectangle.New(vector.V{90, 90}, vector.V{100, 100}),
		102: *hyperrectangle.New(vector.V{9, 9}, vector.V{11, 11}),
		103: *hyperrectangle.New(vector.V{50.1, 50.1}, vector.V{60, 60}),
		104: *hyperrectangle.New(vector.V{0, 0}, vector.V{100, 100}),
	}

	for _, x := range []id.ID{100, 101, 102, 103, 104} {
		fmt.Printf("Adding ID = %v, AABB = %v\n", x, data[x])
		t.Insert(x, data[x])
	}

	q := *hyperrectangle.New(vector.V{10, 10}, vector.V{50, 50})
	fmt.Printf("Querying q = %v\n", q)

	for _, x := range t.BroadPhase(q) {
		fmt.Printf("ID = %v intersects with the query AABB\n", x)
	}
}

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
		101: *hyperrectangle.New(vector.V{10, 10}, vector.V{11, 11}),
		102: *hyperrectangle.New(vector.V{9, 9}, vector.V{11, 11}),
		103: *hyperrectangle.New(vector.V{30, 30}, vector.V{40, 40}),
		104: *hyperrectangle.New(vector.V{100, 100}, vector.V{101, 101}),
		105: *hyperrectangle.New(vector.V{90.01, 90.01}, vector.V{95, 95}),
		106: *hyperrectangle.New(vector.V{0, 0}, vector.V{100, 100}),
	}

	for x, aabb := range data {
		fmt.Printf("Adding ID = %v, AABB = %v\n", x, aabb)
		t.Insert(x, aabb)
	}

	q := *hyperrectangle.New(vector.V{10, 10}, vector.V{50, 50})
	fmt.Printf("Querying q = %v\n", q)

	for _, x := range t.BroadPhase(q) {
		fmt.Printf("ID = %v intersects with the query AABB\n", x)
	}
}

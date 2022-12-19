# go-bvh
Golang AABB-backed BVH implementation

## Example

```golang
t := bvh.New(bvh.O{
	LeafSize: 4,
	K: vector.D(2),
	Tolerance: 1.05,
})

t.Insert(100, *hyperrectangle.New(vector.V{0, 0}, vector.V{10, 10}))

for _, x := range t.BroadPhase(*hyperrectangle.New(vector.V{9, 0}, vector.V{11, 1})) {
	fmt.Printf("ID %v overlaps query", x)
}

t.Update(100, *hyperrectangle.New(vector.V{0, 0}, vector.V{5, 10}))

// No objects remain in the query rectangle.
for _, x := range t.BroadPhase(*hyperrectangle.New(vector.V{9, 0}, vector.V{11, 1})) {
	fmt.Printf("ID %v overlaps query", x)
}
```

See [example.go](example/example.go) for a full implementation.

## Testing

```
go test -race -coverprofile=coverage.out github.com/downflux/go-bvh/...
go tool cover -html=coverage.out
```

## Sources

1. Gold, Nash. "Is BVH faster than the octree/kd-tree for raytracing the objects on a GPU?" https://computergraphics.stackexchange.com/q/10098. Aug 2011.
1. Catto, Erin. "Dynamic Bounding Volume Hierarchies." https://box2d.org/files/ErinCatto_DynamicBVH_Full.pdf. 2019.
1. Randall, James. "Introductory Guide to AABB Tree Collision Detection." https://www.azurefromthetrenches.com/introductory-guide-to-aabb-tree-collision-detection/. 2017.
1. @erincatto. "Box2D." https://github.com/erincatto/box2d. 2021.
1. @ChubbyBubba91. "Implementing 2D unit collision for top down RTS." https://www.reddit.com/r/gamedev/comments/9osh44. 2018.
1. @briannoyama. "Online Bounding Volume Hierarchy." https://github.com/briannoyama/bvh. 2020.
1. Hedges, Lester. "AABB.cc." https://github.com/lohedges/aabbcc. 2021.
1. MacDonald, J.D., and Kellogg Booth. "Heuristics for Tay Tracing Using Space Subdivision." https://graphicsinterface.org/wp-content/uploads/gi1989-22.pdf. 1990.
1. Aila, et al. "On Quality Metrics of Bounding Volume Hierarchies." https://users.aalto.fi/~ailat1/publications/aila2013hpg_paper.pdf. 2013.
1. Bittner, et al. "Fast Insertion-Based Optimization of Bounding Volume Hierarchies."  https://doi.org/10.1111/cgf.12000. 2013.
1. Kopta, et al. "Fast, Effective BVH Updates for Animated Scenes." https://hwrt.cs.utah.edu/papers/hwrt_rotations.pdf. 2012.
1. Guttman, Antonin. "R-trees: A Dynamic Index Structure for Spacial Searching." http://www-db.deis.unibo.it/courses/SI-LS/papers/Gut84.pdf. 1984.
1. @dhconnelly. "rtreego." https://github.com/dhconnelly/rtreego. 2019.

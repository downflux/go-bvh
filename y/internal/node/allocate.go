package node

func New(k vector.D, n int, tolerance float64, x id.ID) *N {
	if tolerance < 1 {
		panic(fmt.Sprintf("cannot set expansion tolerance to be less than resultant AABB"))
	}

	aabb := *hyperrectangle.New(
		vector.V(make([]float64, k)),
		vector.V(make([]float64, k)),
	)

	return &N{
		id: x,

		k:         k,
		n:         n,
		tolerance: tolerance,

		aabbCache: aabb,

		children: make([]*N, 0, n),
	}
}

func Free(n *N) {
	// AABB will set by subsequent update calls. Skip zeroing the bounding
	// box.

	for x := range n.children {
		delete(n.children, x)
	}

	n.isLeaf = false
}

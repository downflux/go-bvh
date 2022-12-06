package node

// Update makes the node up-to-date. This function assumes the child nodes are
// also up-to-date.
func Update(n *N, data map[id.ID]hyperrectangle.R) {
	if n.IsLeaf() {
		n.heightCache = 0
	} else {
		h := 0
		for _, n := range n.children {
			if g := n.Height(); g > h {
				h = g
			}
		}

		n.heightCache = h + 1
	}

	aabb := n.AABB().M()
	var init bool
	if !n.IsLeaf() {
		for _, m := range n.children {
			if !init {
				init = true
				aabb.Copy(m.AABB())
			} else {
				aabb.Union(m.AABB())
			}
		}
	} else {
		aabb.Copy(data[n.leaf])

		tmin, tmax := aabb.Min(), aabb.Max()
		f := math.Pow(n.tolerance, 1/float64(k))

		for i := vector.D(0); i < k; i++ {
			delta := tmax[i] - tmin[i]
			offset := delta * (n.tolerance - 1) / 2
			tmin[i] = tmin[i] - offset
			tmax[i] = tmax[i] + offset
		}
	}
	heuristic.H(aabb.R())
}

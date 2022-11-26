package balance

import (
	"testing"

	"github.com/downflux/go-bvh/x/id"
	"github.com/downflux/go-bvh/x/internal/cache"
	"github.com/downflux/go-bvh/x/internal/cache/node"
	"github.com/downflux/go-bvh/x/internal/heuristic"
	"github.com/downflux/go-geometry/epsilon"
	"github.com/downflux/go-geometry/nd/hyperrectangle"
	"github.com/downflux/go-geometry/nd/vector"

	cid "github.com/downflux/go-bvh/x/internal/cache/id"
)

func TestCheckDF(t *testing.T) {
	const k = 2
	type w struct {
		h       float64
		optimal bool
	}

	type config struct {
		name string
		d    node.N
		e    node.N
		f    node.N
		g    node.N
		opt  float64
		want w
	}

	configs := []config{
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        k,
			})

			nd := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			nd.AABB().Copy(*hyperrectangle.New(vector.V([]float64{1, 1}), vector.V([]float64{2, 2})))

			ne := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			ne.AABB().Copy(*hyperrectangle.New(vector.V([]float64{99, 99}), vector.V([]float64{100, 100})))

			nf := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			nf.AABB().Copy(*hyperrectangle.New(vector.V([]float64{97, 97}), vector.V([]float64{98, 98})))

			ng := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			ng.AABB().Copy(*hyperrectangle.New(vector.V([]float64{3, 3}), vector.V([]float64{4, 4})))

			nb := c.GetOrDie(c.Insert(cid.IDInvalid, nd.ID(), ne.ID(), true))
			nc := c.GetOrDie(c.Insert(cid.IDInvalid, nf.ID(), ng.ID(), true))
			na := c.GetOrDie(c.Insert(cid.IDInvalid, nb.ID(), nc.ID(), true))

			for _, n := range []node.N{nb, nc, na} {
				node.SetAABB(n, nil, 1)
				node.SetHeight(n)
			}

			return config{
				name: "Swap",
				d:    nd,
				e:    ne,
				f:    nf,
				g:    ng,
				opt:  heuristic.H(nb.AABB().R()) + heuristic.H(nc.AABB().R()),
				want: w{
					h:       24,
					optimal: true,
				},
			}
		}(),
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        k,
			})

			nd := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			nd.AABB().Copy(*hyperrectangle.New(vector.V([]float64{97, 97}), vector.V([]float64{98, 98})))

			ne := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			ne.AABB().Copy(*hyperrectangle.New(vector.V([]float64{99, 99}), vector.V([]float64{100, 100})))

			nf := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			nf.AABB().Copy(*hyperrectangle.New(vector.V([]float64{1, 1}), vector.V([]float64{2, 2})))

			ng := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			ng.AABB().Copy(*hyperrectangle.New(vector.V([]float64{3, 3}), vector.V([]float64{4, 4})))

			nb := c.GetOrDie(c.Insert(cid.IDInvalid, nd.ID(), ne.ID(), true))
			nc := c.GetOrDie(c.Insert(cid.IDInvalid, nf.ID(), ng.ID(), true))
			na := c.GetOrDie(c.Insert(cid.IDInvalid, nb.ID(), nc.ID(), true))

			for _, n := range []node.N{nb, nc, na} {
				node.SetAABB(n, nil, 1)
				node.SetHeight(n)
			}

			return config{
				name: "NoSwap/AABB",
				d:    nd,
				e:    ne,
				f:    nf,
				g:    ng,
				opt:  heuristic.H(nb.AABB().R()) + heuristic.H(nc.AABB().R()),
				want: w{
					h:       24,
					optimal: false,
				},
			}
		}(),
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        k,
			})

			nd := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			nd.AABB().Copy(*hyperrectangle.New(vector.V([]float64{1, 1}), vector.V([]float64{2, 2})))
			nd.SetHeight(2)

			ne := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			ne.AABB().Copy(*hyperrectangle.New(vector.V([]float64{99, 99}), vector.V([]float64{100, 100})))

			nf := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			nf.AABB().Copy(*hyperrectangle.New(vector.V([]float64{97, 97}), vector.V([]float64{98, 98})))

			ng := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			ng.AABB().Copy(*hyperrectangle.New(vector.V([]float64{3, 3}), vector.V([]float64{4, 4})))
			ng.SetHeight(2)

			nb := c.GetOrDie(c.Insert(cid.IDInvalid, nd.ID(), ne.ID(), true))
			nc := c.GetOrDie(c.Insert(cid.IDInvalid, nf.ID(), ng.ID(), true))
			na := c.GetOrDie(c.Insert(cid.IDInvalid, nb.ID(), nc.ID(), true))

			for _, n := range []node.N{nb, nc, na} {
				node.SetAABB(n, nil, 1)
				node.SetHeight(n)
			}

			return config{
				name: "NoSwap/Balance",
				d:    nd,
				e:    ne,
				f:    nf,
				g:    ng,
				opt:  heuristic.H(nb.AABB().R()) + heuristic.H(nc.AABB().R()),
				want: w{
					h:       776,
					optimal: false,
				},
			}
		}(),
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        k,
			})

			nd := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			nd.AABB().Copy(*hyperrectangle.New(vector.V([]float64{1, 1}), vector.V([]float64{2, 2})))
			nd.SetHeight(2)

			ne := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			ne.AABB().Copy(*hyperrectangle.New(vector.V([]float64{99, 99}), vector.V([]float64{100, 100})))
			ne.SetHeight(2)

			nf := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			nf.AABB().Copy(*hyperrectangle.New(vector.V([]float64{97, 97}), vector.V([]float64{98, 98})))

			ng := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			ng.AABB().Copy(*hyperrectangle.New(vector.V([]float64{3, 3}), vector.V([]float64{4, 4})))

			nb := c.GetOrDie(c.Insert(cid.IDInvalid, nd.ID(), ne.ID(), true))
			nc := c.GetOrDie(c.Insert(cid.IDInvalid, nf.ID(), ng.ID(), true))
			na := c.GetOrDie(c.Insert(cid.IDInvalid, nb.ID(), nc.ID(), true))

			for _, n := range []node.N{nb, nc, na} {
				node.SetAABB(n, nil, 1)
				node.SetHeight(n)
			}

			return config{
				name: "NoSwap/Balance/Child",
				d:    nd,
				e:    ne,
				f:    nf,
				g:    ng,
				opt:  heuristic.H(nb.AABB().R()) + heuristic.H(nc.AABB().R()),
				want: w{
					h:       776,
					optimal: false,
				},
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			buf := hyperrectangle.New(
				vector.V(make([]float64, k)),
				vector.V(make([]float64, k)),
			).M()
			got := &w{}
			got.h, got.optimal = checkDF(c.d, c.e, c.f, c.g, c.opt, buf)
			if got.h != c.want.h {
				t.Errorf("h = %v, want = %v", got.h, c.want.h)
			}
			if got.optimal != c.want.optimal {
				t.Errorf("optimal = %v, want = %v", got.optimal, c.want.optimal)
			}
		})
	}
}

func TestCheckBF(t *testing.T) {
	const k = 2
	type w struct {
		h       float64
		optimal bool
	}

	type config struct {
		name string
		b    node.N
		f    node.N
		g    node.N
		opt  float64
		want w
	}

	configs := []config{
		func() config {
			data := map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(vector.V([]float64{1, 1}), vector.V([]float64{2, 2})),
				101: *hyperrectangle.New(vector.V([]float64{99, 99}), vector.V([]float64{100, 100})),
				102: *hyperrectangle.New(vector.V([]float64{3, 3}), vector.V([]float64{4, 4})),
			}
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        k,
			})

			na := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))

			nb := c.GetOrDie(c.Insert(na.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nb.Leaves()[100] = struct{}{}
			na.SetLeft(nb.ID())

			nc := c.GetOrDie(c.Insert(na.ID(), cid.IDInvalid, cid.IDInvalid, true))
			na.SetRight(nc.ID())

			nf := c.GetOrDie(c.Insert(nc.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nf.Leaves()[101] = struct{}{}
			nc.SetLeft(nf.ID())

			ng := c.GetOrDie(c.Insert(nc.ID(), cid.IDInvalid, cid.IDInvalid, true))
			ng.Leaves()[102] = struct{}{}
			nc.SetRight(ng.ID())

			for _, n := range []node.N{ng, nf, nc, nb, na} {
				node.SetAABB(n, data, 1)
				node.SetHeight(n)
			}

			return config{
				name: "Swap",
				b:    nb,
				f:    nf,
				g:    ng,
				opt:  heuristic.H(nb.AABB().R()) + heuristic.H(nc.AABB().R()),
				want: w{
					h:       16,
					optimal: true,
				},
			}
		}(),
		func() config {
			data := map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(vector.V([]float64{99, 99}), vector.V([]float64{100, 100})),
				101: *hyperrectangle.New(vector.V([]float64{1, 1}), vector.V([]float64{2, 2})),
				102: *hyperrectangle.New(vector.V([]float64{3, 3}), vector.V([]float64{4, 4})),
			}
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        k,
			})

			na := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))

			nb := c.GetOrDie(c.Insert(na.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nb.Leaves()[100] = struct{}{}
			na.SetLeft(nb.ID())

			nc := c.GetOrDie(c.Insert(na.ID(), cid.IDInvalid, cid.IDInvalid, true))
			na.SetRight(nc.ID())

			nf := c.GetOrDie(c.Insert(nc.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nf.Leaves()[101] = struct{}{}
			nc.SetLeft(nf.ID())

			ng := c.GetOrDie(c.Insert(nc.ID(), cid.IDInvalid, cid.IDInvalid, true))
			ng.Leaves()[102] = struct{}{}
			nc.SetRight(ng.ID())

			for _, n := range []node.N{ng, nf, nc, nb, na} {
				node.SetAABB(n, data, 1)
				node.SetHeight(n)
			}

			return config{
				name: "NoSwap/AABB",
				b:    nb,
				f:    nf,
				g:    ng,
				opt:  heuristic.H(nb.AABB().R()) + heuristic.H(nc.AABB().R()),
				want: w{
					h:       16,
					optimal: false,
				},
			}
		}(),
		func() config {
			data := map[id.ID]hyperrectangle.R{
				100: *hyperrectangle.New(vector.V([]float64{1, 1}), vector.V([]float64{2, 2})),
				101: *hyperrectangle.New(vector.V([]float64{99, 99}), vector.V([]float64{100, 100})),
				102: *hyperrectangle.New(vector.V([]float64{3, 3}), vector.V([]float64{4, 4})),
			}
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        k,
			})

			na := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))

			nb := c.GetOrDie(c.Insert(na.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nb.Leaves()[100] = struct{}{}
			na.SetLeft(nb.ID())

			nc := c.GetOrDie(c.Insert(na.ID(), cid.IDInvalid, cid.IDInvalid, true))
			na.SetRight(nc.ID())

			nf := c.GetOrDie(c.Insert(nc.ID(), cid.IDInvalid, cid.IDInvalid, true))
			nf.AABB().Copy(data[101])
			nc.SetLeft(nf.ID())
			nf.SetHeight(2) // Manually set height for a pseudo-leaf.

			ng := c.GetOrDie(c.Insert(nc.ID(), cid.IDInvalid, cid.IDInvalid, true))
			ng.AABB().Copy(data[102])
			nc.SetRight(ng.ID())
			ng.SetHeight(2)

			for _, n := range []node.N{nc, nb, na} {
				node.SetAABB(n, data, 1)
				node.SetHeight(n)
			}

			return config{
				name: "NoSwap/Height",
				b:    nb,
				f:    nf,
				g:    ng,
				opt:  heuristic.H(nb.AABB().R()) + heuristic.H(nc.AABB().R()),
				want: w{
					h:       392,
					optimal: false,
				},
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			buf := hyperrectangle.New(
				vector.V(make([]float64, k)),
				vector.V(make([]float64, k)),
			).M()

			got := &w{}
			got.h, got.optimal = checkBF(c.b, c.f, c.g, c.opt, buf)
			if !epsilon.Within(got.h, c.want.h) {
				t.Errorf("h = %v, want = %v", got.h, c.want.h)
			}
			if got.optimal != c.want.optimal {
				t.Errorf("optimal = %v, want = %v", got.optimal, c.want.optimal)
			}
		})
	}
}

func TestMerge(t *testing.T) {
	const k = 2
	type w struct {
		height   int
		balanced bool
		h        float64
	}

	type config struct {
		name string
		l    node.N
		r    node.N
		want w
	}

	configs := []config{
		func() config {
			c := cache.New(cache.O{
				LeafSize: 1,
				K:        k,
			})

			l := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			l.AABB().Copy(*hyperrectangle.New(
				vector.V([]float64{1, 1}),
				vector.V([]float64{2, 2}),
			))
			l.SetHeight(2)

			r := c.GetOrDie(c.Insert(cid.IDInvalid, cid.IDInvalid, cid.IDInvalid, true))
			r.AABB().Copy(*hyperrectangle.New(
				vector.V([]float64{9, 9}),
				vector.V([]float64{11, 11}),
			))
			r.SetHeight(1)

			return config{
				name: "Simple",
				l:    l,
				r:    r,
				want: w{
					height:   3,
					balanced: true,
					h: heuristic.H(*hyperrectangle.New(
						vector.V([]float64{1, 1}),
						vector.V([]float64{11, 11}),
					)),
				},
			}
		}(),
	}

	for _, c := range configs {
		t.Run(c.name, func(t *testing.T) {
			buf := hyperrectangle.New(
				vector.V(make([]float64, k)),
				vector.V(make([]float64, k)),
			).M()
			got := &w{}

			got.height, got.balanced, got.h = merge(c.l, c.r, buf)
			if got.height != c.want.height {
				t.Errorf("height = %v, c.want = %v", got.height, c.want.height)
			}
			if got.balanced != c.want.balanced {
				t.Errorf("balanced = %v, c.want = %v", got.balanced, c.want.balanced)
			}
			if !epsilon.Within(got.h, c.want.h) {
				t.Errorf("h = %v, c.want = %v", got.h, c.want.h)
			}
		})
	}
}

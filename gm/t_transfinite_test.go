// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gm

import (
	"math"
	"testing"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/fun"
	"github.com/cpmech/gosl/io"
	"github.com/cpmech/gosl/la"
	"github.com/cpmech/gosl/plt"
	"github.com/cpmech/gosl/utl"
)

func TestTransfinite01(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Transfinite01")

	π := math.Pi
	e0 := []float64{1, 0}
	e1 := []float64{0, 1}

	// mapping
	trf := NewTransfinite(2, []fun.Vs{
		func(x la.Vector, r float64) {
			for i := 0; i < 2; i++ {
				x[i] = (2 + r) * e0[i]
			}
		},
		func(x la.Vector, s float64) {
			θ := π * (s + 1) / 4.0
			for i := 0; i < 2; i++ {
				x[i] = 3*math.Cos(θ)*e0[i] + 3*math.Sin(θ)*e1[i]
			}
		},
		func(x la.Vector, r float64) {
			for i := 0; i < 2; i++ {
				x[i] = (2 + r) * e1[i]
			}
		},
		func(x la.Vector, s float64) {
			θ := π * (s + 1) / 4.0
			for i := 0; i < 2; i++ {
				x[i] = math.Cos(θ)*e0[i] + math.Sin(θ)*e1[i]
			}
		},
	}, []fun.Vs{
		func(dxdr la.Vector, r float64) {
			for i := 0; i < 2; i++ {
				dxdr[i] = e0[i]
			}
		},
		func(dxds la.Vector, s float64) {
			θ := π * (s + 1) / 4.0
			dθds := π / 4.0
			for i := 0; i < 2; i++ {
				dxds[i] = (-3*math.Sin(θ)*e0[i] + 3*math.Cos(θ)*e1[i]) * dθds
			}
		},
		func(dxdr la.Vector, r float64) {
			for i := 0; i < 2; i++ {
				dxdr[i] = e1[i]
			}
		},
		func(dxds la.Vector, s float64) {
			θ := π * (s + 1) / 4.0
			dθds := π / 4.0
			for i := 0; i < 2; i++ {
				dxds[i] = (-math.Sin(θ)*e0[i] + math.Cos(θ)*e1[i]) * dθds
			}
		},
	})

	// check corners
	chk.Array(tst, "C0", 1e-17, trf.C[0], []float64{1, 0})
	chk.Array(tst, "C1", 1e-17, trf.C[1], []float64{3, 0})
	chk.Array(tst, "C2", 1e-17, trf.C[2], []float64{0, 3})
	chk.Array(tst, "C3", 1e-17, trf.C[3], []float64{0, 1})

	// auxiliary
	a := 1.0 / math.Sqrt(2)
	b := 2.0 / math.Sqrt(2)
	c := 3.0 / math.Sqrt(2)
	x := la.NewVector(2)

	// check points
	trf.Point(x, []float64{-1, -1})
	chk.Array(tst, "x(-1,-1)", 1e-17, x, []float64{1, 0})

	trf.Point(x, []float64{0, -1})
	chk.Array(tst, "x( 0,-1)", 1e-17, x, []float64{2, 0})

	trf.Point(x, []float64{+1, -1})
	chk.Array(tst, "x(+1,-1)", 1e-17, x, []float64{3, 0})

	trf.Point(x, []float64{-1, 0})
	chk.Array(tst, "x(-1, 0)", 1e-15, x, []float64{a, a})

	trf.Point(x, []float64{0, 0})
	chk.Array(tst, "x( 0, 0)", 1e-15, x, []float64{b, b})

	trf.Point(x, []float64{+1, 0})
	chk.Array(tst, "x(+1, 0)", 1e-15, x, []float64{c, c})

	trf.Point(x, []float64{-1, +1})
	chk.Array(tst, "x(-1,+1)", 1e-15, x, []float64{0, 1})

	trf.Point(x, []float64{0, +1})
	chk.Array(tst, "x( 0,+1)", 1e-15, x, []float64{0, 2})

	trf.Point(x, []float64{+1, +1})
	chk.Array(tst, "x(+1,+1)", 1e-15, x, []float64{0, 3})

	// check derivs
	dxdu := la.NewMatrix(2, 2)
	u := la.NewVector(2)
	rvals := utl.LinSpace(-1, 1, 3)
	svals := utl.LinSpace(-1, 1, 3)
	verb := false
	for _, s := range svals {
		for _, r := range rvals {
			u[0] = r
			u[1] = s
			trf.Derivs(dxdu, x, u)
			chk.DerivVecVec(tst, "dx/dr", 1e-9, dxdu.GetDeep2(), u, 1e-3, verb, func(xx, rr []float64) {
				trf.Point(xx, rr)
			})
		}
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		trf.Draw([]int{21, 21}, &plt.A{C: plt.C(2, 9)}, &plt.A{C: plt.C(3, 9), Lw: 2})
		plt.Arc(0, 0, 1, 0, 90, &plt.A{C: plt.C(5, 9), NoClip: true, Z: 10})
		plt.Arc(0, 0, 3, 0, 90, &plt.A{C: plt.C(5, 9), NoClip: true, Z: 10})
		for _, s := range svals {
			for _, r := range rvals {
				u[0] = r
				u[1] = s
				trf.Derivs(dxdu, x, u)
				DrawArrow2dM(x, dxdu, 0, true, 0.3, &plt.A{C: plt.C(0, 0), Scale: 7, Z: 10})
				DrawArrow2dM(x, dxdu, 1, true, 0.3, &plt.A{C: plt.C(1, 0), Scale: 7, Z: 10})
			}
		}
		plt.HideAllBorders()
		plt.Equal()
		plt.Save("/tmp/gosl/gm", "transfinite01")
	}
}

func TestTransfinite02(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Transfinite02")

	π := math.Pi
	e0 := []float64{1, 0}
	e1 := []float64{0, 1}

	trf := NewTransfinite(2, []fun.Vs{
		func(x la.Vector, r float64) {
			for i := 0; i < len(x); i++ {
				x[i] = (2 + r) * e0[i]
			}
		},
		func(x la.Vector, s float64) {
			for i := 0; i < len(x); i++ {
				x[i] = 1.5*(1-s)*e0[i] + 1.5*(1+s)*e1[i]
			}
		},
		func(x la.Vector, r float64) {
			for i := 0; i < len(x); i++ {
				x[i] = (2 + r) * e1[i]
			}
		},
		func(x la.Vector, s float64) {
			for i := 0; i < len(x); i++ {
				θ := π * (s + 1) / 4.0
				x[i] = math.Cos(θ)*e0[i] + math.Sin(θ)*e1[i]
			}
		},
	}, []fun.Vs{
		func(dxdr la.Vector, r float64) {
			for i := 0; i < 2; i++ {
				dxdr[i] = e0[i]
			}
		},
		func(dxds la.Vector, s float64) {
			θ := π * (s + 1) / 4.0
			dθds := π / 4.0
			for i := 0; i < 2; i++ {
				dxds[i] = (-3*math.Sin(θ)*e0[i] + 3*math.Cos(θ)*e1[i]) * dθds
			}
		},
		func(dxdr la.Vector, r float64) {
			for i := 0; i < 2; i++ {
				dxdr[i] = e1[i]
			}
		},
		func(dxds la.Vector, s float64) {
			θ := π * (s + 1) / 4.0
			dθds := π / 4.0
			for i := 0; i < 2; i++ {
				dxds[i] = (-math.Sin(θ)*e0[i] + math.Cos(θ)*e1[i]) * dθds
			}
		},
	})

	chk.Array(tst, "C0", 1e-17, trf.C[0], []float64{1, 0})
	chk.Array(tst, "C1", 1e-17, trf.C[1], []float64{3, 0})
	chk.Array(tst, "C2", 1e-17, trf.C[2], []float64{0, 3})
	chk.Array(tst, "C3", 1e-17, trf.C[3], []float64{0, 1})

	a := 1.0 / math.Sqrt(2)
	c := 1.5
	b := (a + c) / 2.0
	x := la.NewVector(2)

	trf.Point(x, []float64{-1, -1})
	chk.Array(tst, "x(-1,-1)", 1e-17, x, []float64{1, 0})

	trf.Point(x, []float64{0, -1})
	chk.Array(tst, "x( 0,-1)", 1e-17, x, []float64{2, 0})

	trf.Point(x, []float64{+1, -1})
	chk.Array(tst, "x(+1,-1)", 1e-17, x, []float64{3, 0})

	trf.Point(x, []float64{-1, 0})
	chk.Array(tst, "x(-1, 0)", 1e-15, x, []float64{a, a})

	trf.Point(x, []float64{0, 0})
	chk.Array(tst, "x( 0, 0)", 1e-15, x, []float64{b, b})

	trf.Point(x, []float64{+1, 0})
	chk.Array(tst, "x(+1, 0)", 1e-15, x, []float64{c, c})

	trf.Point(x, []float64{-1, +1})
	chk.Array(tst, "x(-1,+1)", 1e-15, x, []float64{0, 1})

	trf.Point(x, []float64{0, +1})
	chk.Array(tst, "x( 0,+1)", 1e-15, x, []float64{0, 2})

	trf.Point(x, []float64{+1, +1})
	chk.Array(tst, "x(+1,+1)", 1e-15, x, []float64{0, 3})

	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		trf.Draw([]int{21, 21}, nil, nil)
		plt.Arc(0, 0, 1, 0, 90, &plt.A{C: plt.C(2, 0), NoClip: true, Z: 10})
		plt.HideAllBorders()
		plt.Equal()
		plt.Save("/tmp/gosl/gm", "transfinite02")
	}
}

func TestTransfinite03(tst *testing.T) {

	//verbose()
	chk.PrintTitle("Transfinite03")

	// boundary functions
	curve0 := FactoryNurbs.Curve2dExample1()
	e0 := []float64{1, 0}
	e1 := []float64{0, 1}
	knot := []float64{0}
	trf := NewTransfinite(2, []fun.Vs{

		// B0
		func(x la.Vector, r float64) {
			knot[0] = (1 + r) / 2.0
			for i := 0; i < len(x); i++ {
				curve0.Point(x, knot, 2)
			}
		},

		// B1
		func(x la.Vector, s float64) {
			x[0] = 3
			x[1] = 1.5 * (1 + s) * e1[1]
		},

		// B2
		func(x la.Vector, r float64) {
			x[0] = 1.5 * (1 + r) * e0[0]
			x[1] = 3
		},

		// B3
		func(x la.Vector, s float64) {
			x[0] = 0
			x[1] = 1.5 * (1 + s) * e1[1]
		},
	}, []fun.Vs{

		// dB0/dr
		func(dxdr la.Vector, r float64) {
			knot[0] = (1 + r) / 2.0
			dCdu := la.NewMatrix(2, curve0.Gnd())
			C := la.NewVector(2)
			curve0.PointDeriv(dCdu, C, knot, 2)
			for i := 0; i < 2; i++ {
				dxdr[i] = dCdu.Get(i, 0) * 0.5
			}
		},

		// dB1/ds
		func(dxds la.Vector, s float64) {
			dxds[0] = 0
			dxds[1] = 1.5 * e1[1]
		},

		// dB2/dr
		func(dxdr la.Vector, r float64) {
			dxdr[0] = 1.5 * e0[0]
			dxdr[1] = 0
		},

		// dB3/ds
		func(dxds la.Vector, s float64) {
			dxds[0] = 0
			dxds[1] = 1.5 * e1[1]
		},
	})

	// auxiliary
	xtmp := la.NewVector(2)
	dxdr := la.NewVector(2)
	dxds := la.NewVector(2)
	rvals := utl.LinSpace(-1, 1, 5)
	svals := utl.LinSpace(-1, 1, 5)

	// check: dB0/dr
	//verb := chk.Verbose
	verb := false
	for _, r := range rvals {
		trf.Bd[0](dxdr, r)
		for i := 0; i < 2; i++ {
			chk.DerivScaSca(tst, io.Sf("dB0_%d/dr", i), 1e-10, dxdr[i], r, 1e-3, verb, func(s float64) float64 {
				trf.B[0](xtmp, s)
				return xtmp[i]
			})
		}
	}

	// check: dB1/ds
	io.Pl()
	for _, s := range svals {
		trf.Bd[1](dxds, s)
		for i := 0; i < 2; i++ {
			chk.DerivScaSca(tst, io.Sf("dB1_%d/ds", i), 1e-12, dxds[i], s, 1e-3, verb, func(s float64) float64 {
				trf.B[1](xtmp, s)
				return xtmp[i]
			})
		}
	}

	// check: dB2/dr
	io.Pl()
	for _, r := range rvals {
		trf.Bd[2](dxdr, r)
		for i := 0; i < 2; i++ {
			chk.DerivScaSca(tst, io.Sf("dB2_%d/dr", i), 1e-12, dxdr[i], r, 1e-3, verb, func(s float64) float64 {
				trf.B[2](xtmp, s)
				return xtmp[i]
			})
		}
	}

	// check: dB3/ds
	io.Pl()
	for _, s := range svals {
		trf.Bd[3](dxds, s)
		for i := 0; i < 2; i++ {
			chk.DerivScaSca(tst, io.Sf("dB3_%d/ds", i), 1e-12, dxds[i], s, 1e-3, verb, func(s float64) float64 {
				trf.B[3](xtmp, s)
				return xtmp[i]
			})
		}
	}

	// check derivs
	dxdu := la.NewMatrix(2, 2)
	x := la.NewVector(2)
	u := la.NewVector(2)
	for _, s := range svals {
		for _, r := range rvals {
			u[0] = r
			u[1] = s
			trf.Derivs(dxdu, x, u)
			chk.DerivVecVec(tst, "dx/dr", 1e-9, dxdu.GetDeep2(), u, 1e-3, verb, func(xx, rr []float64) {
				trf.Point(xx, rr)
			})
		}
	}

	// plot
	if chk.Verbose {
		plt.Reset(true, &plt.A{WidthPt: 400, Dpi: 150})
		curve0.DrawElems(2, 41, false, &plt.A{C: plt.C(2, 0), Z: 10}, nil)
		trf.Draw([]int{21, 21}, &plt.A{C: plt.C(2, 9)}, &plt.A{C: plt.C(3, 9), Lw: 2})
		for _, s := range svals {
			for _, r := range rvals {
				u[0] = r
				u[1] = s
				trf.Derivs(dxdu, x, u)
				DrawArrow2dM(x, dxdu, 0, true, 0.3, &plt.A{C: plt.C(0, 0), Scale: 7, Z: 10})
				DrawArrow2dM(x, dxdu, 1, true, 0.3, &plt.A{C: plt.C(1, 0), Scale: 7, Z: 10})
			}
		}
		plt.HideAllBorders()
		plt.Equal()
		plt.Save("/tmp/gosl/gm", "transfinite03")
	}
}

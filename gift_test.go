package gift

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
)

type testFilter struct {
	z int
}

func (p *testFilter) Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle) {
	dstBounds = image.Rect(0, 0, srcBounds.Dx()+p.z, srcBounds.Dy()+p.z*2)
	return
}

func (p *testFilter) Draw(dst draw.Image, src image.Image, options *Options) {
	dst.Set(dst.Bounds().Min.X, dst.Bounds().Min.Y, color.Gray{123})
	return
}

func TestGIFT(t *testing.T) {
	g := New()
	if g.Parallelization() != defaultOptions.Parallelization {
		t.Error("unexpected parallelization property")
	}
	g.SetParallelization(true)
	if g.Parallelization() != true {
		t.Error("unexpected parallelization property")
	}
	g.SetParallelization(false)
	if g.Parallelization() != false {
		t.Error("unexpected parallelization property")
	}

	g = New(
		&testFilter{1},
		&testFilter{2},
		&testFilter{3},
	)
	if len(g.Filters) != 3 {
		t.Error("unexpected filters count")
	}

	g.Add(
		&testFilter{4},
		&testFilter{5},
		&testFilter{6},
	)
	if len(g.Filters) != 6 {
		t.Error("unexpected filters count")
	}
	b := g.Bounds(image.Rect(0, 0, 1, 2))
	if !b.Eq(image.Rect(0, 0, 22, 44)) {
		t.Error("unexpected gift bounds")
	}

	g.Empty()
	if len(g.Filters) != 0 {
		t.Error("unexpected filters count")
	}
	b = g.Bounds(image.Rect(0, 0, 1, 2))
	if !b.Eq(image.Rect(0, 0, 1, 2)) {
		t.Error("unexpected gift bounds")
	}

	g = &GIFT{}
	src := image.NewGray(image.Rect(-1, -1, 1, 1))
	src.Pix = []uint8{
		1, 2,
		3, 4,
	}
	dst := image.NewGray(g.Bounds(src.Bounds()))
	g.Draw(dst, src)
	if !dst.Bounds().Size().Eq(src.Bounds().Size()) {
		t.Error("unexpected dst bounds")
	}
	for i := range dst.Pix {
		if dst.Pix[i] != src.Pix[i] {
			t.Error("unexpected dst pix")
		}
	}

	g.Add(&testFilter{1})
	g.Add(&testFilter{2})
	dst = image.NewGray(g.Bounds(src.Bounds()))
	g.Draw(dst, src)
	if dst.Bounds().Dx() != src.Bounds().Dx()+3 || dst.Bounds().Dy() != src.Bounds().Dy()+6 {
		t.Error("unexpected dst bounds")
	}
	if dst.Pix[0] != 123 {
		t.Error("unexpected dst pix")
	}
}

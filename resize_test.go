package gift

import (
	"image"
	"testing"
)

func TestResize(t *testing.T) {
	var img0, img1 *image.Gray

	// Testing various sizes and parallelization settings
	w, h := 10, 20
	img0 = image.NewGray(image.Rect(0, 0, w, h))
	sz := []struct{ w0, h0, w1, h1 int }{
		{w, h, w, h},
		{w * 2, h, w * 2, h},
		{w, h * 2, w, h * 2},
		{w * 2, h * 2, w * 2, h * 2},
		{w / 2, h, w / 2, h},
		{w, 0, w, h},
		{0, h, w, h},
		{w * 2, 0, w * 2, h * 2},
		{0, h / 2, w / 2, h / 2},
		{0, 0, 0, 0},
		{1, -1, 0, 0},
		{-1, 1, 0, 0},
	}
	rfilters := []Resampling{
		NearestNeighborResampling,
		BoxResampling,
		LinearResampling,
		CubicResampling,
		LanczosResampling,
	}
	for _, prlz := range []bool{true, false} {
		for _, z := range sz {
			for _, f := range rfilters {
				g := New(Resize(z.w0, z.h0, f))
				g.SetParallelization(prlz)
				img1 := image.NewGray(g.Bounds(img0.Bounds()))
				g.Draw(img1, img0)
				w2, h2 := img1.Bounds().Dx(), img1.Bounds().Dy()
				if w2 != z.w1 || h2 != z.h1 {
					t.Errorf("resize %s %dx%d: expected %dx%d got %dx%d", f, z.w0, z.h0, z.w1, z.h1, w2, h2)
				}
			}
		}
	}

	// Nearest filter resize
	img0 = image.NewGray(image.Rect(-1, -1, 4, 1))
	img0.Pix = []uint8{
		1, 2, 3, 4, 5,
		6, 7, 8, 0, 1,
	}
	img1_exp := image.NewGray(image.Rect(0, 0, 2, 2))
	img1_exp.Pix = []uint8{
		2, 4,
		7, 0,
	}
	f := Resize(2, 2, NearestNeighborResampling)
	img1 = image.NewGray(f.Bounds(img0.Bounds()))
	f.Draw(img1, img0, nil)
	if img1.Bounds().Size() != img1_exp.Bounds().Size() {
		t.Errorf("expected %v got %v", img1_exp.Bounds().Size(), img1.Bounds().Size())
	}
	if !comparePix(img1_exp.Pix, img1.Pix) {
		t.Errorf("expected %v got %v", img1_exp.Pix, img1.Pix)
	}

	// Box Filter resize
	img0 = image.NewGray(image.Rect(-1, -1, 3, 1))
	img0.Pix = []uint8{
		1, 2, 2, 1,
		4, 5, 8, 9,
	}
	img1_exp = image.NewGray(image.Rect(0, 0, 2, 1))
	img1_exp.Pix = []uint8{
		3, 5,
	}
	f = Resize(2, 1, BoxResampling)
	img1 = image.NewGray(f.Bounds(img0.Bounds()))
	f.Draw(img1, img0, nil)
	if img1.Bounds().Size() != img1_exp.Bounds().Size() {
		t.Errorf("expected %v got %v", img1_exp.Bounds().Size(), img1.Bounds().Size())
	}
	if !comparePix(img1_exp.Pix, img1.Pix) {
		t.Errorf("expected %v got %v", img1_exp.Pix, img1.Pix)
	}

	// Empty image should remain empty and not panic
	img0 = &image.Gray{}
	f = Resize(100, 100, BoxResampling)
	img1 = image.NewGray(f.Bounds(img0.Bounds()))
	f.Draw(img1, img0, nil)
	if img1.Bounds().Dx() != 0 || img1.Bounds().Dy() != 0 {
		t.Errorf("empty image resized is not empty: %dx%d", img1.Bounds().Dx(), img1.Bounds().Dy())
	}

	// Testing kernel values outside the window
	for _, f := range rfilters {
		if f.Kernel(f.Support()+0.000001) != 0 {
			t.Errorf("filter %s value outside support != 0", f)
		}
	}

	// Testing spline and sinc edge cases
	if sinc(0) != 1 {
		t.Errorf("sinc(0) != 1")
	}
	if bcspline(-2.0, 0.0, 0.5) != 0 {
		t.Errorf("bcspline(2.0, ...) != 0")
	}

	if (resamplingStruct{name: "test"}).String() != "test" {
		t.Error("resamplingStruct String() fail")
	}
}

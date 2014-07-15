package gift

import (
	"image"
	"testing"
)

func comparePix(p1, p2 []uint8) bool {
	if len(p1) != len(p2) {
		return false
	}
	for i := 0; i < len(p1); i++ {
		if p1[i] != p2[i] {
			return false
		}
	}
	return true
}

func TestRotate90(t *testing.T) {
	img0 := image.NewGray(image.Rect(-1, -1, 3, 1))
	img0.Pix = []uint8{
		1, 2, 3, 4,
		5, 6, 7, 8,
	}
	img1_exp := image.NewGray(image.Rect(0, 0, 2, 4))
	img1_exp.Pix = []uint8{
		4, 8,
		3, 7,
		2, 6,
		1, 5,
	}

	f := Rotate90()
	img1 := image.NewGray(f.Bounds(img0.Bounds()))
	f.Draw(img1, img0, nil)

	if img1.Bounds().Size() != img1_exp.Bounds().Size() {
		t.Errorf("expected %v got %v", img1_exp.Bounds().Size(), img1.Bounds().Size())
	}
	if !comparePix(img1_exp.Pix, img1.Pix) {
		t.Errorf("expected %v got %v", img1_exp.Pix, img1.Pix)
	}
}

func TestRotate180(t *testing.T) {
	img0 := image.NewGray(image.Rect(-1, -1, 3, 1))
	img0.Pix = []uint8{
		1, 2, 3, 4,
		5, 6, 7, 8,
	}
	img1_exp := image.NewGray(image.Rect(0, 0, 4, 2))
	img1_exp.Pix = []uint8{
		8, 7, 6, 5,
		4, 3, 2, 1,
	}

	f := Rotate180()
	img1 := image.NewGray(f.Bounds(img0.Bounds()))
	f.Draw(img1, img0, nil)

	if img1.Bounds().Size() != img1_exp.Bounds().Size() {
		t.Errorf("expected %v got %v", img1_exp.Bounds().Size(), img1.Bounds().Size())
	}
	if !comparePix(img1_exp.Pix, img1.Pix) {
		t.Errorf("expected %v got %v", img1_exp.Pix, img1.Pix)
	}
}

func TestRotate270(t *testing.T) {
	img0 := image.NewGray(image.Rect(-1, -1, 3, 1))
	img0.Pix = []uint8{
		1, 2, 3, 4,
		5, 6, 7, 8,
	}
	img1_exp := image.NewGray(image.Rect(0, 0, 2, 4))
	img1_exp.Pix = []uint8{
		5, 1,
		6, 2,
		7, 3,
		8, 4,
	}

	f := Rotate270()
	img1 := image.NewGray(f.Bounds(img0.Bounds()))
	f.Draw(img1, img0, nil)

	if img1.Bounds().Size() != img1_exp.Bounds().Size() {
		t.Errorf("expected %v got %v", img1_exp.Bounds().Size(), img1.Bounds().Size())
	}
	if !comparePix(img1_exp.Pix, img1.Pix) {
		t.Errorf("expected %v got %v", img1_exp.Pix, img1.Pix)
	}
}

func TestFlipHorizontal(t *testing.T) {
	img0 := image.NewGray(image.Rect(-1, -1, 3, 1))
	img0.Pix = []uint8{
		1, 2, 3, 4,
		5, 6, 7, 8,
	}
	img1_exp := image.NewGray(image.Rect(0, 0, 4, 2))
	img1_exp.Pix = []uint8{
		4, 3, 2, 1,
		8, 7, 6, 5,
	}

	f := FlipHorizontal()
	img1 := image.NewGray(f.Bounds(img0.Bounds()))
	f.Draw(img1, img0, nil)

	if img1.Bounds().Size() != img1_exp.Bounds().Size() {
		t.Errorf("expected %v got %v", img1_exp.Bounds().Size(), img1.Bounds().Size())
	}
	if !comparePix(img1_exp.Pix, img1.Pix) {
		t.Errorf("expected %v got %v", img1_exp.Pix, img1.Pix)
	}
}

func TestFlipVertical(t *testing.T) {
	img0 := image.NewGray(image.Rect(-1, -1, 3, 1))
	img0.Pix = []uint8{
		1, 2, 3, 4,
		5, 6, 7, 8,
	}
	img1_exp := image.NewGray(image.Rect(0, 0, 4, 2))
	img1_exp.Pix = []uint8{
		5, 6, 7, 8,
		1, 2, 3, 4,
	}

	f := FlipVertical()
	img1 := image.NewGray(f.Bounds(img0.Bounds()))
	f.Draw(img1, img0, nil)

	if img1.Bounds().Size() != img1_exp.Bounds().Size() {
		t.Errorf("expected %v got %v", img1_exp.Bounds().Size(), img1.Bounds().Size())
	}
	if !comparePix(img1_exp.Pix, img1.Pix) {
		t.Errorf("expected %v got %v", img1_exp.Pix, img1.Pix)
	}
}

func TestTranspose(t *testing.T) {
	img0 := image.NewGray(image.Rect(-1, -1, 3, 1))
	img0.Pix = []uint8{
		1, 2, 3, 4,
		5, 6, 7, 8,
	}
	img1_exp := image.NewGray(image.Rect(0, 0, 2, 4))
	img1_exp.Pix = []uint8{
		1, 5,
		2, 6,
		3, 7,
		4, 8,
	}

	f := Transpose()
	img1 := image.NewGray(f.Bounds(img0.Bounds()))
	f.Draw(img1, img0, nil)

	if img1.Bounds().Size() != img1_exp.Bounds().Size() {
		t.Errorf("expected %v got %v", img1_exp.Bounds().Size(), img1.Bounds().Size())
	}
	if !comparePix(img1_exp.Pix, img1.Pix) {
		t.Errorf("expected %v got %v", img1_exp.Pix, img1.Pix)
	}
}

func TestTransverse(t *testing.T) {
	img0 := image.NewGray(image.Rect(-1, -1, 3, 1))
	img0.Pix = []uint8{
		1, 2, 3, 4,
		5, 6, 7, 8,
	}
	img1_exp := image.NewGray(image.Rect(0, 0, 2, 4))
	img1_exp.Pix = []uint8{
		8, 4,
		7, 3,
		6, 2,
		5, 1,
	}

	f := Transverse()
	img1 := image.NewGray(f.Bounds(img0.Bounds()))
	f.Draw(img1, img0, nil)

	if img1.Bounds().Size() != img1_exp.Bounds().Size() {
		t.Errorf("expected %v got %v", img1_exp.Bounds().Size(), img1.Bounds().Size())
	}
	if !comparePix(img1_exp.Pix, img1.Pix) {
		t.Errorf("expected %v got %v", img1_exp.Pix, img1.Pix)
	}
}

func TestCrop(t *testing.T) {
	testData := []struct {
		desc           string
		r              image.Rectangle
		srcb, dstb     image.Rectangle
		srcPix, dstPix []uint8
	}{
		{
			"crop (0, 0, 0, 0)",
			image.Rect(0, 0, 0, 0),
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 0, 0),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{},
		},
		{
			"crop (1, 1, -1, -1)",
			image.Rectangle{image.Pt(1, 1), image.Pt(-1, -1)},
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 0, 0),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{},
		},
		{
			"crop (-1, 0, 3, 2)",
			image.Rect(-1, 0, 3, 2),
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 4, 2),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0x60, 0xB0, 0xA0, 0xB0,
				0x00, 0x80, 0x00, 0x80,
			},
		},
		{
			"crop (-100, -100, 2, 2)",
			image.Rect(-100, -100, 2, 2),
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 3, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0,
				0x00, 0x80, 0x00,
			},
		},
		{
			"crop (-100, -100, 100, 100)",
			image.Rect(-100, -100, 100, 100),
			image.Rect(-1, -1, 4, 2),
			image.Rect(0, 0, 5, 3),
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
			[]uint8{
				0x00, 0x40, 0x00, 0x40, 0x00,
				0x60, 0xB0, 0xA0, 0xB0, 0x60,
				0x00, 0x80, 0x00, 0x80, 0x00,
			},
		},
	}

	for _, d := range testData {
		src := image.NewGray(d.srcb)
		src.Pix = d.srcPix

		f := Crop(d.r)
		dst := image.NewGray(f.Bounds(src.Bounds()))
		f.Draw(dst, src, nil)

		if !checkBoundsAndPix(dst.Bounds(), d.dstb, dst.Pix, d.dstPix) {
			t.Errorf("test [%s] failed: %#v, %#v", d.desc, dst.Bounds(), dst.Pix)
		}
	}

}

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

func TestFlipVertically(t *testing.T) {
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

	f := FlipVertically()
	img1 := image.NewGray(f.Bounds(img0.Bounds()))
	f.Draw(img1, img0, nil)

	if img1.Bounds().Size() != img1_exp.Bounds().Size() {
		t.Errorf("expected %v got %v", img1_exp.Bounds().Size(), img1.Bounds().Size())
	}
	if !comparePix(img1_exp.Pix, img1.Pix) {
		t.Errorf("expected %v got %v", img1_exp.Pix, img1.Pix)
	}
}

func TestFlipHorizontally(t *testing.T) {
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

	f := FlipHorizontally()
	img1 := image.NewGray(f.Bounds(img0.Bounds()))
	f.Draw(img1, img0, nil)

	if img1.Bounds().Size() != img1_exp.Bounds().Size() {
		t.Errorf("expected %v got %v", img1_exp.Bounds().Size(), img1.Bounds().Size())
	}
	if !comparePix(img1_exp.Pix, img1.Pix) {
		t.Errorf("expected %v got %v", img1_exp.Pix, img1.Pix)
	}
}

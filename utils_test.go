package gift

import (
	"image"
	"image/color"
	"runtime"
	"testing"
)

func testParallelizeN(enabled bool, n, procs int) bool {
	data := make([]bool, n)
	runtime.GOMAXPROCS(procs)
	parallelize(enabled, 0, n, func(start, end int) {
		for i := start; i < end; i++ {
			data[i] = true
		}
	})
	for i := 0; i < n; i++ {
		if !data[i] {
			return false
		}
	}
	return true
}

func TestParallelize(t *testing.T) {
	for _, e := range []bool{true, false} {
		for _, n := range []int{1, 10, 100, 1000} {
			for _, p := range []int{1, 2, 4, 8, 16, 100} {
				if !testParallelizeN(e, n, p) {
					t.Errorf("failed testParallelizeN(%v, %d, %d)", e, n, p)
				}
			}
		}
	}
}

func TestTempImageCopy(t *testing.T) {
	tmp1 := createTempImage(image.Rect(-1, -2, 1, 2))
	if !tmp1.Bounds().Eq(image.Rect(-1, -2, 1, 2)) {
		t.Error("unexpected temp image bounds")
	}
	tmp2 := createTempImage(image.Rect(-3, -4, 3, 4))
	if !tmp2.Bounds().Eq(image.Rect(-3, -4, 3, 4)) {
		t.Error("unexpected temp image bounds")
	}
	copyimage(tmp1, tmp2, nil)
}

func TestQSort(t *testing.T) {
	testData := []struct {
		a, b []float32
	}{
		{
			[]float32{},
			[]float32{},
		},
		{
			[]float32{0.1},
			[]float32{0.1},
		},
		{
			[]float32{0.4, 0.2, 0.5, -0.5, 0.3, 0.0, 0.1},
			[]float32{-0.5, 0.0, 0.1, 0.2, 0.3, 0.4, 0.5},
		},
		{
			[]float32{-10, 10, -20, 20, -30, 30},
			[]float32{-30, -20, -10, 10, 20, 30},
		},
	}

	for _, d := range testData {
		qsortf32(d.a)
		for i := range d.a {
			if d.a[i] != d.b[i] {
				t.Errorf("qsort failed: %#v", d.a)
			}
		}
	}
}

func TestDisk(t *testing.T) {
	testData := []struct {
		ksize int
		k     []float32
	}{
		{
			-5,
			[]float32{},
		},
		{
			0,
			[]float32{},
		},
		{
			1,
			[]float32{1},
		},
		{
			2,
			[]float32{1},
		},
		{
			3,
			[]float32{
				0, 1, 0,
				1, 1, 1,
				0, 1, 0,
			},
		},
		{
			4,
			[]float32{
				0, 1, 0,
				1, 1, 1,
				0, 1, 0,
			},
		},
		{
			5,
			[]float32{
				0, 0, 1, 0, 0,
				0, 1, 1, 1, 0,
				1, 1, 1, 1, 1,
				0, 1, 1, 1, 0,
				0, 0, 1, 0, 0,
			},
		},
		{
			6,
			[]float32{
				0, 0, 1, 0, 0,
				0, 1, 1, 1, 0,
				1, 1, 1, 1, 1,
				0, 1, 1, 1, 0,
				0, 0, 1, 0, 0,
			},
		},
		{
			7,
			[]float32{
				0, 0, 0, 1, 0, 0, 0,
				0, 1, 1, 1, 1, 1, 0,
				0, 1, 1, 1, 1, 1, 0,
				1, 1, 1, 1, 1, 1, 1,
				0, 1, 1, 1, 1, 1, 0,
				0, 1, 1, 1, 1, 1, 0,
				0, 0, 0, 1, 0, 0, 0,
			},
		},
	}

	for _, d := range testData {
		disk := genDisk(d.ksize)
		for i := range disk {
			if disk[i] != d.k[i] {
				t.Errorf("gen disk failed: %d %#v", d.ksize, disk)
			}
		}
	}
}

func TestIsOpaque(t *testing.T) {
	type opqt struct {
		img    image.Image
		opaque bool
	}
	var testData []opqt

	testData = append(testData, opqt{image.NewNRGBA(image.Rect(0, 0, 1, 1)), false})
	testData = append(testData, opqt{image.NewNRGBA64(image.Rect(0, 0, 1, 1)), false})
	testData = append(testData, opqt{image.NewRGBA(image.Rect(0, 0, 1, 1)), false})
	testData = append(testData, opqt{image.NewRGBA64(image.Rect(0, 0, 1, 1)), false})
	testData = append(testData, opqt{image.NewGray(image.Rect(0, 0, 1, 1)), true})
	testData = append(testData, opqt{image.NewGray16(image.Rect(0, 0, 1, 1)), true})
	testData = append(testData, opqt{image.NewYCbCr(image.Rect(0, 0, 1, 1), image.YCbCrSubsampleRatio444), true})
	testData = append(testData, opqt{image.NewAlpha(image.Rect(0, 0, 1, 1)), false})

	img1 := image.NewNRGBA(image.Rect(0, 0, 1, 1))
	img1.Set(0, 0, color.NRGBA{0x00, 0x00, 0x00, 0xff})
	testData = append(testData, opqt{img1, true})
	img2 := image.NewNRGBA64(image.Rect(0, 0, 1, 1))
	img2.Set(0, 0, color.NRGBA{0x00, 0x00, 0x00, 0xff})
	testData = append(testData, opqt{img2, true})
	img3 := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img3.Set(0, 0, color.NRGBA{0x00, 0x00, 0x00, 0xff})
	testData = append(testData, opqt{img3, true})
	img4 := image.NewRGBA64(image.Rect(0, 0, 1, 1))
	img4.Set(0, 0, color.NRGBA{0x00, 0x00, 0x00, 0xff})
	testData = append(testData, opqt{img4, true})
	imgp1 := image.NewPaletted(image.Rect(0, 0, 1, 1), []color.Color{color.NRGBA{0x00, 0x00, 0x00, 0xff}})
	imgp1.SetColorIndex(0, 0, 0)
	testData = append(testData, opqt{imgp1, true})
	imgp2 := image.NewPaletted(image.Rect(0, 0, 1, 1), []color.Color{color.NRGBA{0x00, 0x00, 0x00, 0xfe}})
	imgp2.SetColorIndex(0, 0, 0)
	testData = append(testData, opqt{imgp2, false})

	for _, d := range testData {
		isop := isOpaque(d.img)
		if isop != d.opaque {
			t.Errorf("isOpaque failed %#v, %v", d.img, isop)
		}
	}
}

func checkBoundsAndPix(b1, b2 image.Rectangle, pix1, pix2 []uint8) bool {
	if !b1.Eq(b2) {
		return false
	}
	if len(pix1) != len(pix2) {
		return false
	}
	for i := 0; i < len(pix1); i++ {
		if pix1[i] != pix2[i] {
			return false
		}
	}
	return true
}

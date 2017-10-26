package gift

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"testing"
)

func TestNewPixelGetter(t *testing.T) {
	var img image.Image
	var pg *pixelGetter
	img = image.NewNRGBA(image.Rect(0, 0, 1, 1))
	pg = newPixelGetter(img)
	if pg.it != itNRGBA || pg.nrgba == nil || !img.Bounds().Eq(pg.bounds) {
		t.Error("newPixelGetter NRGBA")
	}
	img = image.NewNRGBA64(image.Rect(0, 0, 1, 1))
	pg = newPixelGetter(img)
	if pg.it != itNRGBA64 || pg.nrgba64 == nil || !img.Bounds().Eq(pg.bounds) {
		t.Error("newPixelGetter NRGBA64")
	}
	img = image.NewRGBA(image.Rect(0, 0, 1, 1))
	pg = newPixelGetter(img)
	if pg.it != itRGBA || pg.rgba == nil || !img.Bounds().Eq(pg.bounds) {
		t.Error("newPixelGetter RGBA")
	}
	img = image.NewRGBA64(image.Rect(0, 0, 1, 1))
	pg = newPixelGetter(img)
	if pg.it != itRGBA64 || pg.rgba64 == nil || !img.Bounds().Eq(pg.bounds) {
		t.Error("newPixelGetter RGBA64")
	}
	img = image.NewGray(image.Rect(0, 0, 1, 1))
	pg = newPixelGetter(img)
	if pg.it != itGray || pg.gray == nil || !img.Bounds().Eq(pg.bounds) {
		t.Error("newPixelGetter Gray")
	}
	img = image.NewGray16(image.Rect(0, 0, 1, 1))
	pg = newPixelGetter(img)
	if pg.it != itGray16 || pg.gray16 == nil || !img.Bounds().Eq(pg.bounds) {
		t.Error("newPixelGetter Gray16")
	}
	img = image.NewYCbCr(image.Rect(0, 0, 1, 1), image.YCbCrSubsampleRatio422)
	pg = newPixelGetter(img)
	if pg.it != itYCbCr || pg.ycbcr == nil || !img.Bounds().Eq(pg.bounds) {
		t.Error("newPixelGetter YCbCr")
	}
	img = image.NewUniform(color.NRGBA64{0, 0, 0, 0})
	pg = newPixelGetter(img)
	if pg.it != itGeneric || pg.image == nil || !img.Bounds().Eq(pg.bounds) {
		t.Error("newPixelGetter Generic(Uniform)")
	}
	img = image.NewAlpha(image.Rect(0, 0, 1, 1))
	pg = newPixelGetter(img)
	if pg.it != itGeneric || pg.image == nil || !img.Bounds().Eq(pg.bounds) {
		t.Error("newPixelGetter Generic(Alpha)")
	}
}

func comparePixels(px1, px2 pixel, dif float64) bool {
	if math.Abs(float64(px1.r)-float64(px2.r)) > dif {
		return false
	}
	if math.Abs(float64(px1.g)-float64(px2.g)) > dif {
		return false
	}
	if math.Abs(float64(px1.b)-float64(px2.b)) > dif {
		return false
	}
	if math.Abs(float64(px1.a)-float64(px2.a)) > dif {
		return false
	}
	return true

}

func compareColorsNRGBA(c1, c2 color.NRGBA, dif int) bool {
	if math.Abs(float64(c1.R)-float64(c2.R)) > float64(dif) {
		return false
	}
	if math.Abs(float64(c1.G)-float64(c2.G)) > float64(dif) {
		return false
	}
	if math.Abs(float64(c1.B)-float64(c2.B)) > float64(dif) {
		return false
	}
	if math.Abs(float64(c1.A)-float64(c2.A)) > float64(dif) {
		return false
	}
	return true
}

func TestGetPixel(t *testing.T) {
	var pg *pixelGetter

	// RGBA, NRGBA, RGBA64, NRGBA64

	palette := []color.Color{
		color.NRGBA{0, 0, 0, 0},
		color.NRGBA{255, 255, 255, 255},
		color.NRGBA{50, 100, 150, 255},
		color.NRGBA{150, 100, 50, 200},
	}

	images1 := []draw.Image{
		image.NewRGBA(image.Rect(-1, -2, 3, 4)),
		image.NewRGBA64(image.Rect(-1, -2, 3, 4)),
		image.NewNRGBA(image.Rect(-1, -2, 3, 4)),
		image.NewNRGBA64(image.Rect(-1, -2, 3, 4)),
		image.NewPaletted(image.Rect(-1, -2, 3, 4), palette),
	}

	colors1 := []struct {
		c  color.NRGBA
		px pixel
	}{
		{color.NRGBA{0, 0, 0, 0}, pixel{0, 0, 0, 0}},
		{color.NRGBA{255, 255, 255, 255}, pixel{1, 1, 1, 1}},
		{color.NRGBA{50, 100, 150, 255}, pixel{0.196, 0.392, 0.588, 1}},
		{color.NRGBA{150, 100, 50, 200}, pixel{0.588, 0.392, 0.196, 0.784}},
	}

	for _, img := range images1 {
		pg = newPixelGetter(img)
		for _, k := range colors1 {
			for _, x := range []int{-1, 0, 2} {
				for _, y := range []int{-2, 0, 3} {
					img.Set(x, y, k.c)
					px := pg.getPixel(x, y)
					if !comparePixels(k.px, px, 0.005) {
						t.Errorf("getPixel %T %v %dx%d %v %v", img, k.c, x, y, k.px, px)
					}
				}
			}
		}
	}

	// Uniform (Generic)

	for _, k := range colors1 {
		img := image.NewUniform(k.c)
		pg = newPixelGetter(img)
		for _, x := range []int{-1, 0, 2} {
			for _, y := range []int{-2, 0, 3} {
				px := pg.getPixel(x, y)
				if !comparePixels(k.px, px, 0.005) {
					t.Errorf("getPixel %T %v %dx%d %v %v", img, k.c, x, y, k.px, px)
				}
			}
		}
	}

	// YCbCr

	colors2 := []struct {
		c  color.NRGBA
		px pixel
	}{
		{color.NRGBA{0, 0, 0, 255}, pixel{0, 0, 0, 1}},
		{color.NRGBA{255, 255, 255, 255}, pixel{1, 1, 1, 1}},
		{color.NRGBA{50, 100, 150, 255}, pixel{0.196, 0.392, 0.588, 1}},
		{color.NRGBA{150, 100, 50, 255}, pixel{0.588, 0.392, 0.196, 1}},
	}

	for _, k := range colors2 {
		for _, sr := range []image.YCbCrSubsampleRatio{
			image.YCbCrSubsampleRatio444,
			image.YCbCrSubsampleRatio422,
			image.YCbCrSubsampleRatio420,
			image.YCbCrSubsampleRatio440,
			image.YCbCrSubsampleRatio410,
			image.YCbCrSubsampleRatio411,
		} {
			img := image.NewYCbCr(image.Rect(-1, -2, 3, 4), sr)
			pg = newPixelGetter(img)
			for _, x := range []int{-1, 0, 2} {
				for _, y := range []int{-2, 0, 3} {
					iy := img.YOffset(x, y)
					ic := img.COffset(x, y)
					img.Y[iy], img.Cb[ic], img.Cr[ic] = color.RGBToYCbCr(k.c.R, k.c.G, k.c.B)
					px := pg.getPixel(x, y)
					if !comparePixels(k.px, px, 0.005) {
						t.Errorf("getPixel %T %v %dx%d %v %v", img, k.c, x, y, k.px, px)
					}
				}
			}
		}
	}

	// Gray, Gray16

	images2 := []draw.Image{
		image.NewGray(image.Rect(-1, -2, 3, 4)),
		image.NewGray16(image.Rect(-1, -2, 3, 4)),
	}

	colors3 := []struct {
		c  color.NRGBA
		px pixel
	}{
		{color.NRGBA{0, 0, 0, 0}, pixel{0, 0, 0, 1}},
		{color.NRGBA{255, 255, 255, 255}, pixel{1, 1, 1, 1}},
		{color.NRGBA{50, 100, 150, 255}, pixel{0.356, 0.356, 0.356, 1}},
		{color.NRGBA{150, 100, 50, 200}, pixel{0.337, 0.337, 0.337, 1}},
	}

	for _, img := range images2 {
		pg = newPixelGetter(img)
		for _, k := range colors3 {
			for _, x := range []int{-1, 0, 2} {
				for _, y := range []int{-2, 0, 3} {
					img.Set(x, y, k.c)
					px := pg.getPixel(x, y)
					if !comparePixels(k.px, px, 0.005) {
						t.Errorf("getPixel %T %v %dx%d %v %v", img, k.c, x, y, k.px, px)
					}
				}
			}
		}
	}
}

func comparePixelSlices(s1, s2 []pixel, dif float64) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 1; i < len(s1); i++ {
		if !comparePixels(s1[i], s2[i], dif) {
			return false
		}
	}
	return true
}

func TestGetPixelRow(t *testing.T) {
	colors := []color.NRGBA{
		{0, 0, 0, 0},
		{255, 255, 255, 255},
		{50, 100, 150, 255},
		{150, 100, 50, 200},
	}
	pixels := []pixel{
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{0.196, 0.392, 0.588, 1},
		{0.588, 0.392, 0.196, 0.784},
	}

	img := image.NewNRGBA(image.Rect(-1, -2, 3, 2))
	pg := newPixelGetter(img)
	var row []pixel
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			img.Set(x, y, colors[x-img.Bounds().Min.X])
		}
		pg.getPixelRow(y, &row)
		if !comparePixelSlices(row, pixels, 0.005) {
			t.Errorf("getPixelRow y=%d %v %v", y, row, pixels)
		}
	}
}

func TestGetPixelColumn(t *testing.T) {
	colors := []color.NRGBA{
		{0, 0, 0, 0},
		{255, 255, 255, 255},
		{50, 100, 150, 255},
		{150, 100, 50, 200},
	}
	pixels := []pixel{
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{0.196, 0.392, 0.588, 1},
		{0.588, 0.392, 0.196, 0.784},
	}

	img := image.NewNRGBA(image.Rect(-1, -2, 3, 2))
	pg := newPixelGetter(img)
	var column []pixel
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			img.Set(x, y, colors[y-img.Bounds().Min.Y])
		}
		pg.getPixelColumn(x, &column)
		if !comparePixelSlices(column, pixels, 0.005) {
			t.Errorf("getPixelColumn x=%d %v %v", x, column, pixels)
		}
	}
}

func TestF32u8(t *testing.T) {
	testData := []struct {
		x float32
		y uint8
	}{
		{-1, 0},
		{0, 0},
		{100, 100},
		{255, 255},
		{256, 255},
	}
	for _, p := range testData {
		v := f32u8(p.x)
		if v != p.y {
			t.Errorf("f32u8(%f) != %d: %d", p.x, p.y, v)
		}
	}
}

func TestF32u16(t *testing.T) {
	testData := []struct {
		x float32
		y uint16
	}{
		{-1, 0},
		{0, 0},
		{1, 1},
		{10000, 10000},
		{65535, 65535},
		{65536, 65535},
	}
	for _, p := range testData {
		v := f32u16(p.x)
		if v != p.y {
			t.Errorf("f32u16(%f) != %d: %d", p.x, p.y, v)
		}
	}
}

func TestClampi32(t *testing.T) {
	testData := []struct {
		x int32
		y int32
	}{
		{-1, 0},
		{0, 0},
		{1, 1},
		{99, 99},
		{100, 100},
		{101, 100},
	}
	for _, p := range testData {
		v := clampi32(p.x, 0, 100)
		if v != p.y {
			t.Errorf("clampi32(%d) != %d: %d", p.x, p.y, v)
		}
	}
}

func TestNewPixelSetter(t *testing.T) {
	var img draw.Image
	var pg *pixelSetter
	img = image.NewNRGBA(image.Rect(0, 0, 1, 1))
	pg = newPixelSetter(img)
	if pg.it != itNRGBA || pg.nrgba == nil || !img.Bounds().Eq(pg.bounds) {
		t.Error("newPixelSetter NRGBA")
	}
	img = image.NewNRGBA64(image.Rect(0, 0, 1, 1))
	pg = newPixelSetter(img)
	if pg.it != itNRGBA64 || pg.nrgba64 == nil || !img.Bounds().Eq(pg.bounds) {
		t.Error("newPixelSetter NRGBA64")
	}
	img = image.NewRGBA(image.Rect(0, 0, 1, 1))
	pg = newPixelSetter(img)
	if pg.it != itRGBA || pg.rgba == nil || !img.Bounds().Eq(pg.bounds) {
		t.Error("newPixelSetter RGBA")
	}
	img = image.NewRGBA64(image.Rect(0, 0, 1, 1))
	pg = newPixelSetter(img)
	if pg.it != itRGBA64 || pg.rgba64 == nil || !img.Bounds().Eq(pg.bounds) {
		t.Error("newPixelSetter RGBA64")
	}
	img = image.NewGray(image.Rect(0, 0, 1, 1))
	pg = newPixelSetter(img)
	if pg.it != itGray || pg.gray == nil || !img.Bounds().Eq(pg.bounds) {
		t.Error("newPixelSetter Gray")
	}
	img = image.NewGray16(image.Rect(0, 0, 1, 1))
	pg = newPixelSetter(img)
	if pg.it != itGray16 || pg.gray16 == nil || !img.Bounds().Eq(pg.bounds) {
		t.Error("newPixelSetter Gray16")
	}
	img = image.NewPaletted(image.Rect(0, 0, 1, 1), color.Palette{})
	pg = newPixelSetter(img)
	if pg.it != itPaletted || pg.paletted == nil || !img.Bounds().Eq(pg.bounds) {
		t.Error("newPixelSetter Paletted")
	}
	img = image.NewAlpha(image.Rect(0, 0, 1, 1))
	pg = newPixelSetter(img)
	if pg.it != itGeneric || pg.image == nil || !img.Bounds().Eq(pg.bounds) {
		t.Error("newPixelSetter Generic(Alpha)")
	}
}

func TestSetPixel(t *testing.T) {
	var ps *pixelSetter

	// RGBA, NRGBA, RGBA64, NRGBA64

	images1 := []draw.Image{
		image.NewRGBA(image.Rect(-1, -2, 3, 4)),
		image.NewRGBA64(image.Rect(-1, -2, 3, 4)),
		image.NewNRGBA(image.Rect(-1, -2, 3, 4)),
		image.NewNRGBA64(image.Rect(-1, -2, 3, 4)),
	}

	colors1 := []struct {
		c  color.NRGBA
		px pixel
	}{
		{color.NRGBA{0, 0, 0, 0}, pixel{0, 0, 0, 0}},
		{color.NRGBA{0, 0, 0, 255}, pixel{0, 0, 0, 1}},
		{color.NRGBA{255, 255, 255, 255}, pixel{1, 1, 1, 1}},
		{color.NRGBA{50, 100, 150, 255}, pixel{0.196, 0.392, 0.588, 1}},
		{color.NRGBA{150, 100, 50, 200}, pixel{0.588, 0.392, 0.196, 0.784}},
	}

	for _, img := range images1 {
		ps = newPixelSetter(img)
		for _, k := range colors1 {
			for _, x := range []int{-1, 0, 2} {
				for _, y := range []int{-2, 0, 3} {
					ps.setPixel(x, y, k.px)
					c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
					if !compareColorsNRGBA(c, k.c, 1) {
						t.Errorf("setPixel %T %v %dx%d %v %v", img, k.px, x, y, k.c, c)
					}
				}
			}
		}
	}

	// Gray, Gray16

	images2 := []draw.Image{
		image.NewGray(image.Rect(-1, -2, 3, 4)),
		image.NewGray16(image.Rect(-1, -2, 3, 4)),
	}

	colors2 := []struct {
		c  color.NRGBA
		px pixel
	}{
		{color.NRGBA{0, 0, 0, 255}, pixel{0, 0, 0, 1}},
		{color.NRGBA{255, 255, 255, 255}, pixel{1, 1, 1, 1}},
		{color.NRGBA{110, 110, 110, 255}, pixel{0.2, 0.5, 0.7, 1}},
		{color.NRGBA{55, 55, 55, 255}, pixel{0.2, 0.5, 0.7, 0.5}},
	}

	for _, img := range images2 {
		ps = newPixelSetter(img)
		for _, k := range colors2 {
			for _, x := range []int{-1, 0, 2} {
				for _, y := range []int{-2, 0, 3} {
					ps.setPixel(x, y, k.px)
					c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
					if !compareColorsNRGBA(c, k.c, 1) {
						t.Errorf("setPixel %T %v %dx%d %v %v", img, k.px, x, y, k.c, c)
					}
				}
			}
		}
	}

	// Generic(Alpha)

	colors3 := []struct {
		c  color.NRGBA
		px pixel
	}{
		{color.NRGBA{255, 255, 255, 255}, pixel{0, 0, 0, 1}},
		{color.NRGBA{255, 255, 255, 127}, pixel{0.2, 0.5, 0.7, 0.5}},
		{color.NRGBA{255, 255, 255, 63}, pixel{0.1, 0.2, 0.3, 0.25}},
	}

	img := image.NewAlpha(image.Rect(-1, -2, 3, 4))
	ps = newPixelSetter(img)
	for _, k := range colors3 {
		for _, x := range []int{-1, 0, 2} {
			for _, y := range []int{-2, 0, 3} {
				ps.setPixel(x, y, k.px)
				c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
				if !compareColorsNRGBA(c, k.c, 1) {
					t.Errorf("setPixel %T %v %dx%d %v %v", img, k.px, x, y, k.c, c)
				}
			}
		}
	}

	// Paletted

	images4 := []draw.Image{
		image.NewPaletted(
			image.Rect(-1, -2, 3, 4),
			color.Palette{
				color.NRGBA{0, 0, 0, 0},
				color.NRGBA{0, 0, 0, 255},
				color.NRGBA{255, 255, 255, 255},
				color.NRGBA{50, 100, 150, 255},
				color.NRGBA{150, 100, 50, 200},
				color.NRGBA{1, 255, 255, 255},
				color.NRGBA{2, 255, 255, 255},
				color.NRGBA{3, 255, 255, 255},
			},
		),
	}

	colors4 := []struct {
		c  color.NRGBA
		px pixel
	}{
		{color.NRGBA{0, 0, 0, 0}, pixel{0, 0, 0, 0}},
		{color.NRGBA{0, 0, 0, 255}, pixel{0, 0, 0, 1}},
		{color.NRGBA{255, 255, 255, 255}, pixel{1, 1, 1, 1}},
		{color.NRGBA{50, 100, 150, 255}, pixel{0.196, 0.392, 0.588, 1}},
		{color.NRGBA{150, 100, 50, 200}, pixel{0.588, 0.392, 0.196, 0.784}},
		{color.NRGBA{0, 0, 0, 0}, pixel{0.1, 0.01, 0.001, 0.1}},
		{color.NRGBA{0, 0, 0, 255}, pixel{0, 0, 0, 0.9}},
		{color.NRGBA{255, 255, 255, 255}, pixel{1, 0.9, 1, 0.9}},
		{color.NRGBA{1, 255, 255, 255}, pixel{0.001 / 255, 1, 1, 1}},
		{color.NRGBA{1, 255, 255, 255}, pixel{1.001 / 255, 1, 1, 1}},
		{color.NRGBA{2, 255, 255, 255}, pixel{2.001 / 255, 1, 1, 1}},
		{color.NRGBA{3, 255, 255, 255}, pixel{3.001 / 255, 1, 1, 1}},
		{color.NRGBA{3, 255, 255, 255}, pixel{4.001 / 255, 1, 1, 1}},
	}

	for _, img := range images4 {
		ps = newPixelSetter(img)
		for _, k := range colors4 {
			for _, x := range []int{-1, 0, 2} {
				for _, y := range []int{-2, 0, 3} {
					ps.setPixel(x, y, k.px)
					c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
					if !compareColorsNRGBA(c, k.c, 0) {
						t.Errorf("setPixel %T %v %dx%d %v %v", img, k.px, x, y, k.c, c)
					}
				}
			}
		}
	}

}

func TestSetPixelRow(t *testing.T) {
	colors := []color.NRGBA{
		{0, 0, 0, 0},
		{255, 255, 255, 255},
		{50, 100, 150, 255},
		{150, 100, 50, 200},
	}
	pixels := []pixel{
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{0.196, 0.392, 0.588, 1},
		{0.588, 0.392, 0.196, 0.784},
	}

	img := image.NewNRGBA(image.Rect(-1, -2, 3, 2))
	ps := newPixelSetter(img)
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		ps.setPixelRow(y, pixels)
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			c := img.At(x, y).(color.NRGBA)
			wantedColor := colors[x-img.Bounds().Min.X]
			if !compareColorsNRGBA(wantedColor, c, 1) {
				t.Errorf("setPixelRow y=%d x=%d %v %v", y, x, wantedColor, c)
			}
		}
	}
}

func TestSetPixelColumn(t *testing.T) {
	colors := []color.NRGBA{
		{0, 0, 0, 0},
		{255, 255, 255, 255},
		{50, 100, 150, 255},
		{150, 100, 50, 200},
	}
	pixels := []pixel{
		{0, 0, 0, 0},
		{1, 1, 1, 1},
		{0.196, 0.392, 0.588, 1},
		{0.588, 0.392, 0.196, 0.784},
	}

	img := image.NewNRGBA(image.Rect(-1, -2, 3, 2))
	ps := newPixelSetter(img)
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		ps.setPixelColumn(x, pixels)
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			c := img.At(x, y).(color.NRGBA)
			wantedColor := colors[y-img.Bounds().Min.Y]
			if !compareColorsNRGBA(wantedColor, c, 1) {
				t.Errorf("setPixelColumn x=%d y=%d %v %v", x, y, wantedColor, c)
			}
		}
	}
}

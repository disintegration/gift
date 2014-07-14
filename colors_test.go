package gift

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"testing"
)

func TestLut(t *testing.T) {
	fn := func(v float32) float32 {
		return v
	}
	for _, size := range []int{10, 100, 1000} {
		lut := prepareLut(size, fn)
		l := len(lut)
		if l != size {
			t.Errorf("LUT bad size: expected %v got %v", size, l)
		}
		if lut[0] != 0 {
			t.Errorf("LUT bad start value: expected 0 got %v", lut[0])
		}
		if lut[l-1] != 1 {
			t.Errorf("LUT bad end value: expected 1 got %v", lut[l-1])
		}
	}
	lut := prepareLut(10000, fn)
	for _, u := range []float32{0.0, 0.0001, 0.5555, 0.9999, 1.0} {
		v := getFromLut(lut, u)
		if math.Abs(float64(v-u)) > 0.0001 {
			t.Errorf("LUT bad value: expected %v got %v", u, v)
		}
	}
}

func TestInvertColors(t *testing.T) {
	src := image.NewGray(image.Rect(0, 0, 256, 1))
	for i := 0; i <= 255; i++ {
		src.Pix[i] = uint8(i)
	}
	g := New(InvertColors())
	dst := image.NewGray(g.Bounds(src.Bounds()))
	g.Draw(dst, src)

	for i := 0; i <= 255; i++ {
		if dst.Pix[i] != 255-src.Pix[i] {
			t.Errorf("InvertColors: index %d: expected %d got %d", i, 255-src.Pix[i], dst.Pix[i])
		}
	}
}

func TestColorspaceSRGBToLinear(t *testing.T) {
	vals := []float32{
		0.00000,
		0.01002,
		0.03310,
		0.07324,
		0.13287,
		0.21404,
		0.31855,
		0.44799,
		0.60383,
		0.78741,
		1.00000,
	}

	imgs := []draw.Image{
		image.NewGray(image.Rect(0, 0, 11, 11)),
		image.NewGray(image.Rect(0, 0, 111, 111)),
		image.NewGray16(image.Rect(0, 0, 11, 11)),
		image.NewGray16(image.Rect(0, 0, 1111, 1111)),
	}
	for _, img := range imgs {
		for i := 0; i <= 10; i++ {
			img.Set(i, 0, color.Gray{uint8(255 * float32(i) / 10.0)})
		}
		img2 := image.NewGray(img.Bounds())
		New(ColorspaceSRGBToLinear()).Draw(img2, img)
		if !img2.Bounds().Size().Eq(img.Bounds().Size()) {
			t.Errorf("ColorspaceSRGBToLinear bad result size: expected %v got %v", img.Bounds().Size(), img2.Bounds().Size())
		}
		for i := 0; i <= 10; i++ {
			expected := uint8(vals[i]*255.0 + 0.5)
			c := img2.At(i, 0).(color.Gray)
			if math.Abs(float64(c.Y)-float64(expected)) > 1 {
				t.Errorf("ColorspaceSRGBToLinear bad color value at index %v expected %v got %v", i, expected, c.Y)
			}
		}
	}
}

func TestColorspaceLinearToSRGB(t *testing.T) {
	vals := []float32{
		0.00000,
		0.34919,
		0.48453,
		0.58383,
		0.66519,
		0.73536,
		0.79774,
		0.85431,
		0.90633,
		0.95469,
		1.00000,
	}

	imgs := []draw.Image{
		image.NewGray(image.Rect(0, 0, 11, 11)),
		image.NewGray(image.Rect(0, 0, 111, 111)),
		image.NewGray16(image.Rect(0, 0, 11, 11)),
		image.NewGray16(image.Rect(0, 0, 1111, 1111)),
	}
	for _, img := range imgs {
		for i := 0; i <= 10; i++ {
			img.Set(i, 0, color.Gray{uint8(255 * float32(i) / 10.0)})
		}
		img2 := image.NewGray(img.Bounds())
		New(ColorspaceLinearToSRGB()).Draw(img2, img)
		if !img2.Bounds().Size().Eq(img.Bounds().Size()) {
			t.Errorf("ColorspaceLinearRGBToSRGB bad result size: expected %v got %v", img.Bounds().Size(), img2.Bounds().Size())
		}
		for i := 0; i <= 10; i++ {
			expected := uint8(vals[i]*255.0 + 0.5)
			c := img2.At(i, 0).(color.Gray)
			if math.Abs(float64(c.Y)-float64(expected)) > 1 {
				t.Errorf("ColorspaceLinearRGBToSRGB bad color value at index %v expected %v got %v", i, expected, c.Y)
			}
		}
	}

}

func TestAdjustGamma(t *testing.T) {
	src := image.NewGray(image.Rect(0, 0, 256, 1))
	dst := image.NewGray(image.Rect(0, 0, 256, 1))
	for i := 0; i <= 255; i++ {
		src.Pix[i] = uint8(i)
	}
	ag := Gamma(2.0)
	ag.Draw(dst, src, nil)

	for i := 100; i <= 150; i++ {
		if dst.Pix[i] <= src.Pix[i] {
			t.Errorf("Gamma unexpected color")
		}
	}

	ag = Gamma(0.5)
	ag.Draw(dst, src, nil)

	for i := 100; i <= 150; i++ {
		if dst.Pix[i] >= src.Pix[i] {
			t.Errorf("Gamma unexpected color")
		}
	}

	ag = Gamma(1.0)
	ag.Draw(dst, src, nil)

	for i := 100; i <= 150; i++ {
		if dst.Pix[i] != src.Pix[i] {
			t.Errorf("Gamma unexpected color")
		}
	}
}

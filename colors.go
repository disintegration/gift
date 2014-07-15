package gift

import (
	"image"
	"image/draw"
	"math"
)

func prepareLut(lutSize int, fn func(float32) float32) []float32 {
	lut := make([]float32, lutSize)
	q := 1 / float32(lutSize-1)
	for v := 0; v < lutSize; v++ {
		u := float32(v) * q
		lut[v] = fn(u)
	}
	return lut
}

func getFromLut(lut []float32, u float32) float32 {
	v := int(u*float32(len(lut)-1) + 0.5)
	return lut[v]
}

type colorchanFilter struct {
	fn  func(float32) float32
	lut bool
}

func (p *colorchanFilter) Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle) {
	dstBounds = image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy())
	return
}

func (p *colorchanFilter) Draw(dst draw.Image, src image.Image, options *Options) {
	if options == nil {
		options = &defaultOptions
	}

	srcb := src.Bounds()
	dstb := dst.Bounds()
	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	var useLut bool
	var lut []float32

	useLut = false
	if p.lut {
		var lutSize int

		it := pixGetter.imgType
		if it == itNRGBA || it == itRGBA || it == itGray || it == itYCbCr {
			lutSize = 0xff + 1
		} else {
			lutSize = 0xffff + 1
		}

		numCalculations := srcb.Dx() * srcb.Dy() * 3
		if numCalculations > lutSize*2 {
			useLut = true
			lut = prepareLut(lutSize, p.fn)
		}
	}

	parallelize(options.Parallelization, srcb.Min.Y, srcb.Max.Y, func(pmin, pmax int) {
		for y := pmin; y < pmax; y++ {
			for x := srcb.Min.X; x < srcb.Max.X; x++ {
				px := pixGetter.getPixel(x, y)
				if useLut {
					px.R = getFromLut(lut, px.R)
					px.G = getFromLut(lut, px.G)
					px.B = getFromLut(lut, px.B)
				} else {
					px.R = p.fn(px.R)
					px.G = p.fn(px.G)
					px.B = p.fn(px.B)
				}
				pixSetter.setPixel(dstb.Min.X+x-srcb.Min.X, dstb.Min.Y+y-srcb.Min.Y, px)
			}
		}
	})
}

// Invert creates a filter that negates the colors of an image.
func Invert() Filter {
	return &colorchanFilter{
		fn: func(x float32) float32 {
			return 1.0 - x
		},
		lut: false,
	}
}

// ColorspaceSRGBToLinear creates a filter that converts the colors of an image from sRGB to linear RGB.
func ColorspaceSRGBToLinear() Filter {
	return &colorchanFilter{
		fn: func(x float32) float32 {
			if x <= 0.04045 {
				return x / 12.92
			} else {
				return float32(math.Pow(float64((x+0.055)/1.055), 2.4))
			}
		},
		lut: true,
	}
}

// ColorspaceLinearToSRGB creates a filter that converts the colors of an image from linear RGB to sRGB.
func ColorspaceLinearToSRGB() Filter {
	return &colorchanFilter{
		fn: func(x float32) float32 {
			if x <= 0.0031308 {
				return x * 12.92
			} else {
				return float32(1.055*math.Pow(float64(x), 1.0/2.4) - 0.055)
			}
		},
		lut: true,
	}
}

// Gamma creates a filter that performs a gamma correction on an image.
// The gamma parameter must be positive. Gamma = 1.0 gives the original image.
// Gamma less than 1.0 darkens the image and gamma greater than 1.0 lightens it.
func Gamma(gamma float32) Filter {
	e := 1.0 / maxf32(gamma, 1.0e-5)
	return &colorchanFilter{
		fn: func(x float32) float32 {
			return powf32(x, e)
		},
		lut: true,
	}
}

func sigmoid(a, b, x float32) float32 {
	return 1 / (1 + expf32(b*(a-x)))
}

// Sigmoid creates a filter that changes the contrast of an image using a sigmoidal function and returns the adjusted image.
// It's a non-linear contrast change useful for photo adjustments as it preserves highlight and shadow detail.
// The midpoint parameter is the midpoint of contrast that must be between 0 and 1, typically 0.5.
// The factor parameter indicates how much to increase or decrease the contrast, typically in range (-10, 10).
// If the factor parameter is positive the image contrast is increased otherwise the contrast is decreased.
//
// Example:
//
//	g := gift.New(
//		gift.Sigmoid(0.5, 3.0),
//	)
//	dst := image.NewRGBA(src.Bounds())
//	g.Draw(dst, src)
//
func Sigmoid(midpoint, factor float32) Filter {
	a := minf32(maxf32(midpoint, 0.0), 1.0)
	b := absf32(factor)
	sig0 := sigmoid(a, b, 0)
	sig1 := sigmoid(a, b, 1)
	e := float32(1.0e-5)

	return &colorchanFilter{
		fn: func(x float32) float32 {
			if factor == 0 {
				return x
			} else if factor > 0 {
				sig := sigmoid(a, b, x)
				return (sig - sig0) / (sig1 - sig0)
			} else {
				arg := minf32(maxf32((sig1-sig0)*x+sig0, e), 1.0-e)
				return a - logf32(1.0/arg-1.0)/b
			}
		},
		lut: true,
	}
}

// Contrast creates a filter that changes the contrast of an image.
// The percentage parameter must be in range (-100, 100). The percentage = 0 gives the original image.
// The percentage = -100 gives solid grey image. The percentage = 100 gives an overcontrasted image.
func Contrast(percentage float32) Filter {
	percentage = minf32(maxf32(percentage, -100.0), 100.0)
	v := (100.0 + percentage) / 100.0

	return &colorchanFilter{
		fn: func(x float32) float32 {
			if 0 <= v && v <= 1 {
				return 0.5 + (x-0.5)*v
			} else if 1 < v && v < 2 {
				return 0.5 + (x-0.5)*(1/(2.0-v))
			} else {
				if x < 0.5 {
					return 0.0
				} else {
					return 1.0
				}
			}
		},
		lut: false,
	}
}

// Brightness creates a filter that changes the brightness of an image.
// The percentage parameter must be in range (-100, 100). The percentage = 0 gives the original image.
// The percentage = -100 gives solid black image. The percentage = 100 gives solid white image.
func Brightness(percentage float32) Filter {
	percentage = minf32(maxf32(percentage, -100.0), 100.0)
	shift := percentage / 100.0

	return &colorchanFilter{
		fn: func(x float32) float32 {
			return x + shift
		},
		lut: false,
	}
}

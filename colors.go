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
	if percentage == 0 {
		return &copyimageFilter{}
	}

	p := 1 + minf32(maxf32(percentage, -100.0), 100.0)/100.0

	return &colorchanFilter{
		fn: func(x float32) float32 {
			if 0 <= p && p <= 1 {
				return 0.5 + (x-0.5)*p
			} else if 1 < p && p < 2 {
				return 0.5 + (x-0.5)*(1/(2.0-p))
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
	if percentage == 0 {
		return &copyimageFilter{}
	}

	shift := minf32(maxf32(percentage, -100.0), 100.0) / 100.0

	return &colorchanFilter{
		fn: func(x float32) float32 {
			return x + shift
		},
		lut: false,
	}
}

type colorFilter struct {
	fn func(pixel) pixel
}

func (p *colorFilter) Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle) {
	dstBounds = image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy())
	return
}

func (p *colorFilter) Draw(dst draw.Image, src image.Image, options *Options) {
	if options == nil {
		options = &defaultOptions
	}

	srcb := src.Bounds()
	dstb := dst.Bounds()
	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	parallelize(options.Parallelization, srcb.Min.Y, srcb.Max.Y, func(pmin, pmax int) {
		for y := pmin; y < pmax; y++ {
			for x := srcb.Min.X; x < srcb.Max.X; x++ {
				px := pixGetter.getPixel(x, y)
				pixSetter.setPixel(dstb.Min.X+x-srcb.Min.X, dstb.Min.Y+y-srcb.Min.Y, p.fn(px))
			}
		}
	})
}

// Grayscale creates a filter that produces a grayscale version of an image.
func Grayscale() Filter {
	return &colorFilter{
		fn: func(px pixel) pixel {
			y := 0.299*px.R + 0.587*px.G + 0.114*px.B
			return pixel{y, y, y, px.A}
		},
	}
}

// Sepia creates a filter that changes the tint of an image and returns the adjusted image.
// It takes a parameter for how much the image should be adjusted, typically in the range (0, 100)
//
// Example:
//
//	g := gift.New(
//		gift.Sepia(100),
//	)
//	dst := image.NewRGBA(src.Bounds())
//	g.Draw(dst, src)
//
func Sepia(adjust int) Filter {
	return &colorFilter{
		fn: func(px pixel) pixel {
			adjustAmount := float32(adjust) / 100.0
			calculatedR := (px.R * (1.0 - (0.607 * adjustAmount))) +
				(px.G * (0.769 * adjustAmount)) +
				(px.B * (0.189 * adjustAmount))
			calculatedG := (px.R * (0.349 * adjustAmount)) +
				(px.G * (1.0 - (0.314 * adjustAmount))) +
				(px.B * (0.168 * adjustAmount))
			calculatedB := (px.R * (0.272 * adjustAmount)) +
				(px.G * (0.534 * adjustAmount)) +
				(px.B * (1.0 - (float32(0.869) * adjustAmount)))
			r := float32(math.Min(255.0, float64(calculatedR)))
			g := float32(math.Min(255.0, float64(calculatedG)))
			b := float32(math.Min(255.0, float64(calculatedB)))
			return pixel{r, g, b, px.A}
		},
	}
}

func convertHSLToRGB(h, s, l float32) (float32, float32, float32) {
	if s == 0.0 {
		return l, l, l
	}

	_v := func(p, q, t float32) float32 {
		if t < 0.0 {
			t += 1.0
		}
		if t > 1.0 {
			t -= 1.0
		}
		if t < 1.0/6.0 {
			return p + (q-p)*6.0*t
		}
		if t < 1.0/2.0 {
			return q
		}
		if t < 2.0/3.0 {
			return p + (q-p)*(2.0/3.0-t)*6.0
		}
		return p
	}

	var p, q float32
	if l < 0.5 {
		q = l * (1.0 + s)
	} else {
		q = l + s - l*s
	}
	p = 2.0*l - q

	r := _v(p, q, h+1.0/3.0)
	g := _v(p, q, h)
	b := _v(p, q, h-1.0/3.0)

	return r, g, b
}

func convertRGBToHSL(r, g, b float32) (float32, float32, float32) {
	max := maxf32(r, maxf32(g, b))
	min := minf32(r, minf32(g, b))

	l := (max + min) / 2.0

	if max == min {
		return 0.0, 0.0, l
	}

	var h, s float32
	d := max - min
	if l > 0.5 {
		s = d / (2.0 - max - min)
	} else {
		s = d / (max + min)
	}

	if r == max {
		h = (g - b) / d
		if g < b {
			h += 6.0
		}
	} else if g == max {
		h = (b-r)/d + 2.0
	} else {
		h = (r-g)/d + 4.0
	}
	h /= 6.0

	return h, s, l
}

// Hue creates a filter that rotates the hue of an image.
// The shift parameter must be in range (-180, 180). The shift = 0 gives the original image.
func Hue(shift float32) Filter {
	if shift == 0 {
		return &copyimageFilter{}
	}
	p := minf32(maxf32(shift, -180.0), 180.0) / 360.0
	if p < 0.0 {
		p += 1.0
	}
	return &colorFilter{
		fn: func(px pixel) pixel {
			h, s, l := convertRGBToHSL(px.R, px.G, px.B)
			h = h + p
			if h > 1.0 {
				h -= 1.0
			}
			r, g, b := convertHSLToRGB(h, s, l)
			return pixel{r, g, b, px.A}
		},
	}
}

// Saturation creates a filter that changes the saturation of an image.
// The percentage parameter must be in range (-100, 500). The percentage = 0 gives the original image.
func Saturation(percentage float32) Filter {
	if percentage == 0 {
		return &copyimageFilter{}
	}
	p := 1 + minf32(maxf32(percentage, -100.0), 500.0)/100.0

	return &colorFilter{
		fn: func(px pixel) pixel {
			h, s, l := convertRGBToHSL(px.R, px.G, px.B)
			s *= p
			if s > 1 {
				s = 1
			}
			r, g, b := convertHSLToRGB(h, s, l)
			return pixel{r, g, b, px.A}
		},
	}
}

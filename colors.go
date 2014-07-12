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

func mapColorChannels(dst draw.Image, src image.Image, fn func(float32) float32, options *Options) {
	srcb := src.Bounds()
	dstb := dst.Bounds()
	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	var useLut bool
	var lut []float32
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
		lut = prepareLut(lutSize, fn)
	} else {
		useLut = false
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
					px.R = fn(px.R)
					px.G = fn(px.G)
					px.B = fn(px.B)
				}
				pixSetter.setPixel(dstb.Min.X+x-srcb.Min.X, dstb.Min.Y+y-srcb.Min.Y, px)
			}
		}
	})
}

type colorchanFilter struct {
	fn func(float32) float32
}

func (p *colorchanFilter) Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle) {
	dstBounds = image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy())
	return
}

func (p *colorchanFilter) Draw(dst draw.Image, src image.Image, options *Options) {
	if options == nil {
		options = &defaultOptions
	}
	mapColorChannels(dst, src, p.fn, options)
}

// InvertColors creates a filter that negates the colors of an image.
func InvertColors() Filter {
	return &colorchanFilter{
		fn: func(x float32) float32 {
			return 1.0 - x
		},
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
	}
}

// AdjustGamma creates a filter that performs a gamma correction on an image.
// Gamma parameter must be positive. Gamma = 1.0 gives the original image.
// Gamma less than 1.0 darkens the image and gamma greater than 1.0 lightens it. .
func AdjustGamma(gamma float32) Filter {
	e := 1.0 / math.Max(float64(gamma), 0.0001)
	return &colorchanFilter{
		fn: func(x float32) float32 {
			return float32(math.Pow(float64(x), e))
		},
	}
}

package gift

import (
	"image"
	"image/draw"
	"math"
)

type Resampling interface {
	Support() float32
	Kernel(float32) float32
}

func bcspline(x, b, c float32) float32 {
	if x < 0.0 {
		x = -x
	}
	if x < 1.0 {
		return ((12-9*b-6*c)*x*x*x + (-18+12*b+6*c)*x*x + (6 - 2*b)) / 6
	}
	if x < 2.0 {
		return ((-b-6*c)*x*x*x + (6*b+30*c)*x*x + (-12*b-48*c)*x + (8*b + 24*c)) / 6
	}
	return 0
}

func sinc(x float32) float32 {
	if x == 0 {
		return 1.0
	}
	return float32(math.Sin(math.Pi*float64(x)) / (math.Pi * float64(x)))
}

type resamplingStruct struct {
	name    string
	support float32
	kernel  func(float32) float32
}

func (r resamplingStruct) String() string {
	return r.name
}

func (r resamplingStruct) Support() float32 {
	return r.support
}

func (r resamplingStruct) Kernel(x float32) float32 {
	return r.kernel(x)
}

// Nearest neighbor resampling filter.
var NearestNeighborResampling Resampling

// Box resampling filter.
var BoxResampling Resampling

// Linear resampling filter.
var LinearResampling Resampling

// Cubic resampling filter (Catmull-Rom).
var CubicResampling Resampling

// Lanczos resampling filter (3 lobes).
var LanczosResampling Resampling

func precomputeResamplingWeights(dstSize, srcSize int, resampling Resampling) [][]uweight {
	du := float32(srcSize) / float32(dstSize)
	scale := du
	if scale < 1.0 {
		scale = 1.0
	}
	ru := float32(math.Ceil(float64(scale * resampling.Support())))

	tmp := make([]float32, int(ru+2)*2)
	result := make([][]uweight, dstSize)

	for v := 0; v < dstSize; v++ {
		fU := (float32(v)+0.5)*du - 0.5

		startu := int(math.Ceil(float64(fU - ru)))
		if startu < 0 {
			startu = 0
		}
		endu := int(math.Floor(float64(fU + ru)))
		if endu > srcSize-1 {
			endu = srcSize - 1
		}

		sumf := float32(0.0)
		for u := startu; u <= endu; u++ {
			w := resampling.Kernel((float32(u) - fU) / scale)
			sumf += w
			tmp[u-startu] = w
		}
		for u := startu; u <= endu; u++ {
			w := float32(tmp[u-startu] / sumf)
			result[v] = append(result[v], uweight{u, w})
		}
	}

	return result
}

func resizeLine(dstBuf []pixel, srcBuf []pixel, weights [][]uweight) {
	for dstu := 0; dstu < len(weights); dstu++ {
		var r, g, b, a float32
		for _, iw := range weights[dstu] {
			c := srcBuf[iw.u]
			r += c.R * iw.weight
			g += c.G * iw.weight
			b += c.B * iw.weight
			a += c.A * iw.weight
		}
		dstBuf[dstu] = pixel{r, g, b, a}
	}
}

func resizeHorizontal(dst draw.Image, src image.Image, w int, resampling Resampling, options *Options) {
	srcb := src.Bounds()
	dstb := dst.Bounds()

	weights := precomputeResamplingWeights(w, srcb.Dx(), resampling)

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	parallelize(options.Parallelization, srcb.Min.Y, srcb.Max.Y, func(pmin, pmax int) {
		srcBuf := make([]pixel, srcb.Dx())
		dstBuf := make([]pixel, w)
		for srcy := pmin; srcy < pmax; srcy++ {
			pixGetter.getPixelRow(srcy, &srcBuf)
			resizeLine(dstBuf, srcBuf, weights)
			pixSetter.setPixelRow(dstb.Min.Y+srcy-srcb.Min.Y, dstBuf)
		}
	})
}

func resizeVertical(dst draw.Image, src image.Image, h int, resampling Resampling, options *Options) {
	srcb := src.Bounds()
	dstb := dst.Bounds()

	weights := precomputeResamplingWeights(h, srcb.Dy(), resampling)

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	parallelize(options.Parallelization, srcb.Min.X, srcb.Max.X, func(pmin, pmax int) {
		srcBuf := make([]pixel, srcb.Dy())
		dstBuf := make([]pixel, h)
		for srcx := pmin; srcx < pmax; srcx++ {
			pixGetter.getPixelColumn(srcx, &srcBuf)
			resizeLine(dstBuf, srcBuf, weights)
			pixSetter.setPixelColumn(dstb.Min.X+srcx-srcb.Min.X, dstBuf)
		}
	})
}

func resizeNearest(dst draw.Image, src image.Image, w, h int, options *Options) {
	srcb := src.Bounds()
	dstb := dst.Bounds()
	dx := float64(srcb.Dx()) / float64(w)
	dy := float64(srcb.Dy()) / float64(h)

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	parallelize(options.Parallelization, dstb.Min.Y, dstb.Min.Y+h, func(pmin, pmax int) {
		for dsty := pmin; dsty < pmax; dsty++ {
			for dstx := dstb.Min.X; dstx < dstb.Min.X+w; dstx++ {
				fx := math.Floor((float64(dstx-dstb.Min.X) + 0.5) * dx)
				fy := math.Floor((float64(dsty-dstb.Min.Y) + 0.5) * dy)
				srcx := srcb.Min.X + int(fx)
				srcy := srcb.Min.Y + int(fy)
				px := pixGetter.getPixel(srcx, srcy)
				pixSetter.setPixel(dstx, dsty, px)
			}
		}
	})
}

type resizeFilter struct {
	width      int
	height     int
	resampling Resampling
}

func (p *resizeFilter) Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle) {
	w, h := p.width, p.height
	srcw, srch := srcBounds.Dx(), srcBounds.Dy()

	if (w == 0 && h == 0) || w < 0 || h < 0 || srcw <= 0 || srch <= 0 {
		dstBounds = image.Rect(0, 0, 0, 0)
	} else if w == 0 {
		fw := float64(h) * float64(srcw) / float64(srch)
		dstw := int(math.Max(1.0, math.Floor(fw+0.5)))
		dstBounds = image.Rect(0, 0, dstw, h)
	} else if h == 0 {
		fh := float64(w) * float64(srch) / float64(srcw)
		dsth := int(math.Max(1.0, math.Floor(fh+0.5)))
		dstBounds = image.Rect(0, 0, w, dsth)
	} else {
		dstBounds = image.Rect(0, 0, w, h)
	}

	return
}

func (p *resizeFilter) Draw(dst draw.Image, src image.Image, options *Options) {
	if options == nil {
		options = &defaultOptions
	}

	b := p.Bounds(src.Bounds())
	w, h := b.Dx(), b.Dy()

	if w <= 0 || h <= 0 {
		return
	}

	if src.Bounds().Dx() == w && src.Bounds().Dy() == h {
		copyimage(dst, src, options)
		return
	}

	if p.resampling.Support() <= 0 {
		resizeNearest(dst, src, w, h, options)
		return
	}

	if src.Bounds().Dx() == w {
		resizeVertical(dst, src, h, p.resampling, options)
		return
	}

	if src.Bounds().Dy() == h {
		resizeHorizontal(dst, src, w, p.resampling, options)
		return
	}

	tmp := createTempImage(image.Rect(0, 0, w, src.Bounds().Dy()))
	resizeHorizontal(tmp, src, w, p.resampling, options)
	resizeVertical(dst, tmp, h, p.resampling, options)
	return
}

// Resize creates a filter that resizes an image to the specified width and height using the specified resampling.
// If one of width or height is 0, the image aspect ratio is preserved.
// Supported resampling parameters: NearestNeighborResampling, BoxResampling, LinearResampling, CubicResampling, LanczosResampling.
//
// Example:
//
//	// Resize the src image to width=300 preserving the aspect ratio.
//	g := gift.New(
//		gift.Resize(300, 0, gift.LanczosResampling),
//	)
//	dst := image.NewRGBA(g.Bounds(src.Bounds()))
//	g.Draw(dst, src)
//
func Resize(width, height int, resampling Resampling) Filter {
	return &resizeFilter{
		width:      width,
		height:     height,
		resampling: resampling,
	}
}

func init() {
	// Nearest neighbor resampling filter.
	NearestNeighborResampling = resamplingStruct{
		name:    "NearestNeighborResampling",
		support: 0.0,
		kernel: func(x float32) float32 {
			return 0
		},
	}

	// Box resampling filter.
	BoxResampling = resamplingStruct{
		name:    "BoxResampling",
		support: 0.5,
		kernel: func(x float32) float32 {
			if x < 0.0 {
				x = -x
			}
			if x <= 0.5 {
				return 1.0
			}
			return 0
		},
	}

	// Linear resampling filter.
	LinearResampling = resamplingStruct{
		name:    "LinearResampling",
		support: 1.0,
		kernel: func(x float32) float32 {
			if x < 0.0 {
				x = -x
			}
			if x < 1.0 {
				return 1.0 - x
			}
			return 0
		},
	}

	// Cubic resampling filter (Catmull-Rom).
	CubicResampling = resamplingStruct{
		name:    "CubicResampling",
		support: 2.0,
		kernel: func(x float32) float32 {
			if x < 0.0 {
				x = -x
			}
			if x < 2.0 {
				return bcspline(x, 0.0, 0.5)
			}
			return 0
		},
	}

	// Lanczos resampling filter (3 lobes).
	LanczosResampling = resamplingStruct{
		name:    "LanczosResampling",
		support: 3.0,
		kernel: func(x float32) float32 {
			if x < 0.0 {
				x = -x
			}
			if x < 3.0 {
				return sinc(x) * sinc(x/3.0)
			}
			return 0
		},
	}
}

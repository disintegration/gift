package gift

import (
	"image"
	"image/draw"
)

type transformType int

const (
	ttRotate90 transformType = iota
	ttRotate180
	ttRotate270
	ttTranspose
	ttFlipVertically
	ttFlipHorizontally
)

func transform(dst draw.Image, src image.Image, tt transformType, options *Options) {
	srcb := src.Bounds()
	dstb := dst.Bounds()

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	parallelize(options.Parallelization, srcb.Min.Y, srcb.Max.Y, func(pmin, pmax int) {
		for srcy := pmin; srcy < pmax; srcy++ {
			for srcx := srcb.Min.X; srcx < srcb.Max.X; srcx++ {
				var dstx, dsty int
				switch tt {
				case ttRotate90:
					dstx = dstb.Min.X + srcy - srcb.Min.Y
					dsty = dstb.Min.Y + srcb.Max.X - srcx - 1
				case ttRotate180:
					dstx = dstb.Min.X + srcb.Max.X - srcx - 1
					dsty = dstb.Min.Y + srcb.Max.Y - srcy - 1
				case ttRotate270:
					dstx = dstb.Min.X + srcb.Max.Y - srcy - 1
					dsty = dstb.Min.Y + srcx - srcb.Min.X
				case ttTranspose:
					dstx = dstb.Min.X + srcy - srcb.Min.Y
					dsty = dstb.Min.Y + srcx - srcb.Min.X
				case ttFlipVertically:
					dstx = dstb.Min.X + srcx - srcb.Min.X
					dsty = dstb.Min.Y + srcb.Max.Y - srcy - 1
				case ttFlipHorizontally:
					dstx = dstb.Min.X + srcb.Max.X - srcx - 1
					dsty = dstb.Min.Y + srcy - srcb.Min.Y
				}
				pixSetter.setPixel(dstx, dsty, pixGetter.getPixel(srcx, srcy))
			}
		}
	})
}

type transformFilter struct {
	tt transformType
}

func (p *transformFilter) Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle) {
	if p.tt == ttRotate90 || p.tt == ttRotate270 || p.tt == ttTranspose {
		dstBounds = image.Rect(0, 0, srcBounds.Dy(), srcBounds.Dx())
	} else {
		dstBounds = image.Rect(0, 0, srcBounds.Dx(), srcBounds.Dy())
	}
	return
}

func (p *transformFilter) Draw(dst draw.Image, src image.Image, options *Options) {
	if options == nil {
		options = &defaultOptions
	}
	transform(dst, src, p.tt, options)
}

// Rotate90 creates a filter that rotates an image by 90 degrees counter-clockwise.
func Rotate90() Filter {
	return &transformFilter{
		tt: ttRotate90,
	}
}

// Rotate180 creates a filter that rotates an image by 180 degrees counter-clockwise.
func Rotate180() Filter {
	return &transformFilter{
		tt: ttRotate180,
	}
}

// Rotate270 creates a filter that rotates an image by 270 degrees counter-clockwise.
func Rotate270() Filter {
	return &transformFilter{
		tt: ttRotate270,
	}
}

// Transpose creates a filter that transposes an image (reflects over the main diagonal).
func Transpose() Filter {
	return &transformFilter{
		tt: ttTranspose,
	}
}

// FlipVertically creates a filter that flips an image vertically.
func FlipVertically() Filter {
	return &transformFilter{
		tt: ttFlipVertically,
	}
}

// FlipHorizontally creates a filter that flips an image horizontally.
func FlipHorizontally() Filter {
	return &transformFilter{
		tt: ttFlipHorizontally,
	}
}

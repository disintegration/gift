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
	ttFlipHorizontal
	ttFlipVertical
	ttTranspose
	ttTransverse
)

type transformFilter struct {
	tt transformType
}

func (p *transformFilter) Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle) {
	if p.tt == ttRotate90 || p.tt == ttRotate270 || p.tt == ttTranspose || p.tt == ttTransverse {
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

	srcb := src.Bounds()
	dstb := dst.Bounds()

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	parallelize(options.Parallelization, srcb.Min.Y, srcb.Max.Y, func(pmin, pmax int) {
		for srcy := pmin; srcy < pmax; srcy++ {
			for srcx := srcb.Min.X; srcx < srcb.Max.X; srcx++ {
				var dstx, dsty int
				switch p.tt {
				case ttRotate90:
					dstx = dstb.Min.X + srcy - srcb.Min.Y
					dsty = dstb.Min.Y + srcb.Max.X - srcx - 1
				case ttRotate180:
					dstx = dstb.Min.X + srcb.Max.X - srcx - 1
					dsty = dstb.Min.Y + srcb.Max.Y - srcy - 1
				case ttRotate270:
					dstx = dstb.Min.X + srcb.Max.Y - srcy - 1
					dsty = dstb.Min.Y + srcx - srcb.Min.X
				case ttFlipHorizontal:
					dstx = dstb.Min.X + srcb.Max.X - srcx - 1
					dsty = dstb.Min.Y + srcy - srcb.Min.Y
				case ttFlipVertical:
					dstx = dstb.Min.X + srcx - srcb.Min.X
					dsty = dstb.Min.Y + srcb.Max.Y - srcy - 1
				case ttTranspose:
					dstx = dstb.Min.X + srcy - srcb.Min.Y
					dsty = dstb.Min.Y + srcx - srcb.Min.X
				case ttTransverse:
					dstx = dstb.Min.Y + srcb.Max.Y - srcy - 1
					dsty = dstb.Min.X + srcb.Max.X - srcx - 1
				}
				pixSetter.setPixel(dstx, dsty, pixGetter.getPixel(srcx, srcy))
			}
		}
	})
}

// Rotate90 creates a filter that rotates an image 90 degrees counter-clockwise.
func Rotate90() Filter {
	return &transformFilter{
		tt: ttRotate90,
	}
}

// Rotate180 creates a filter that rotates an image 180 degrees counter-clockwise.
func Rotate180() Filter {
	return &transformFilter{
		tt: ttRotate180,
	}
}

// Rotate270 creates a filter that rotates an image 270 degrees counter-clockwise.
func Rotate270() Filter {
	return &transformFilter{
		tt: ttRotate270,
	}
}

// FlipHorizontal creates a filter that flips an image horizontally.
func FlipHorizontal() Filter {
	return &transformFilter{
		tt: ttFlipHorizontal,
	}
}

// FlipVertical creates a filter that flips an image vertically.
func FlipVertical() Filter {
	return &transformFilter{
		tt: ttFlipVertical,
	}
}

// Transpose creates a filter that flips an image horizontally and rotates 90 degrees counter-clockwise.
func Transpose() Filter {
	return &transformFilter{
		tt: ttTranspose,
	}
}

// Transverse creates a filter that flips an image vertically and rotates 90 degrees counter-clockwise.
func Transverse() Filter {
	return &transformFilter{
		tt: ttTransverse,
	}
}

type cropFilter struct {
	rect image.Rectangle
}

func (p *cropFilter) Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle) {
	b := srcBounds.Intersect(p.rect)
	return b.Sub(b.Min)
}

func (p *cropFilter) Draw(dst draw.Image, src image.Image, options *Options) {
	if options == nil {
		options = &defaultOptions
	}

	srcb := src.Bounds().Intersect(p.rect)
	dstb := dst.Bounds()
	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	parallelize(options.Parallelization, srcb.Min.Y, srcb.Max.Y, func(pmin, pmax int) {
		for srcy := pmin; srcy < pmax; srcy++ {
			for srcx := srcb.Min.X; srcx < srcb.Max.X; srcx++ {
				dstx := dstb.Min.X + srcx - srcb.Min.X
				dsty := dstb.Min.Y + srcy - srcb.Min.Y
				pixSetter.setPixel(dstx, dsty, pixGetter.getPixel(srcx, srcy))
			}
		}
	})
}

// Crop creates a filter that crops the specified rectangular region from an image.
//
// Example:
//
//	g := gift.New(
//		gift.Crop(image.Rect(100, 100, 200, 200)),
//	)
//	dst := image.NewRGBA(g.Bounds(src.Bounds()))
//	g.Draw(dst, src)
//
func Crop(rect image.Rectangle) Filter {
	return &cropFilter{
		rect: rect,
	}
}

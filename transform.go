package gift

import (
	"image"
	"image/color"
	"image/draw"
	"math"
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

type Interpolation int

const (
	NearestNeighborInterpolation Interpolation = iota
	LinearInterpolation
	CubicInterpolation
)

func rotatePoint(x, y, asin, acos float32) (float32, float32) {
	newx := x*acos - y*asin
	newy := x*asin + y*acos
	return newx, newy
}

// returns the distance between top left corner of original
// image and top left corner of rotated image as xDiff, yDiff
func RotateOffSet(w, h int, angle float32) (float32, float32) {

	if w <= 0 || h <= 0 {
		return 0, 0
	}

	// find center of src rect for rotation point x0,y0
	xoff := float32(w)/2 - 0.5
	yoff := float32(h)/2 - 0.5

	// find pre-rotation point of src top left corner
	origx := 0 - xoff
	origy := float32(h-1) - yoff

	// find post-rotation point
	asin, acos := sincosf32(angle)
	newx, newy := rotatePoint(origx, origy, asin, acos)

	// find top left corner of dst rect
	x1, y1 := rotatePoint(0-xoff, 0-yoff, asin, acos)
	x2, y2 := rotatePoint(float32(w-1)-xoff, 0-yoff, asin, acos)
	x3, y3 := rotatePoint(float32(w-1)-xoff, float32(h-1)-yoff, asin, acos)
	x4, y4 := rotatePoint(0-xoff, float32(h-1)-yoff, asin, acos)

	minx := minf32(x1, minf32(x2, minf32(x3, x4)))
	maxy := maxf32(y1, maxf32(y2, maxf32(y3, y4)))

	xoffset := float32(0)
	yoffset := float32(0)

	// calculate distance between (newx,newy) and (minx,miny)
	if angle >= 360 {
		angle = angle - 360
	}
	switch {
	case 0 <= angle && angle < 90:
		xoffset = newx - minx
		yoffset = maxy - newy
	case 90 <= angle && angle < 180:
		xoffset = float32(math.Abs(float64(newx - minx)))
		yoffset = maxy - newy
	case 180 <= angle && angle < 270:
		xoffset = newx - minx
		yoffset = float32(math.Abs(float64(newy - maxy)))
	case 270 <= angle && angle < 360:
		xoffset = newx - minx
		yoffset = newy - maxy
	}

	return float32(xoffset), float32(yoffset)

}

func calcRotatedSize(w, h int, angle float32) (int, int) {
	if w <= 0 || h <= 0 {
		return 0, 0
	}

	xoff := float32(w)/2 - 0.5
	yoff := float32(h)/2 - 0.5

	asin, acos := sincosf32(angle)

	x1, y1 := rotatePoint(0-xoff, 0-yoff, asin, acos)
	x2, y2 := rotatePoint(float32(w-1)-xoff, 0-yoff, asin, acos)
	x3, y3 := rotatePoint(float32(w-1)-xoff, float32(h-1)-yoff, asin, acos)
	x4, y4 := rotatePoint(0-xoff, float32(h-1)-yoff, asin, acos)

	minx := minf32(x1, minf32(x2, minf32(x3, x4)))
	maxx := maxf32(x1, maxf32(x2, maxf32(x3, x4)))
	miny := minf32(y1, minf32(y2, minf32(y3, y4)))
	maxy := maxf32(y1, maxf32(y2, maxf32(y3, y4)))

	neww := maxx - minx + 1
	if neww-floorf32(neww) > 0.01 {
		neww += 2
	}
	newh := maxy - miny + 1
	if newh-floorf32(newh) > 0.01 {
		newh += 2
	}
	return int(neww), int(newh)
}

type rotateFilter struct {
	angle         float32
	bgcolor       color.Color
	interpolation Interpolation
}

func (p *rotateFilter) Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle) {
	w, h := calcRotatedSize(srcBounds.Dx(), srcBounds.Dy(), p.angle)
	dstBounds = image.Rect(0, 0, w, h)
	return
}

func (p *rotateFilter) Draw(dst draw.Image, src image.Image, options *Options) {
	if options == nil {
		options = &defaultOptions
	}

	srcb := src.Bounds()
	dstb := dst.Bounds()

	w, h := calcRotatedSize(srcb.Dx(), srcb.Dy(), p.angle)
	if w <= 0 || h <= 0 {
		return
	}

	// find middle of src and dst rects
	srcxoff := float32(srcb.Dx())/2 - 0.5
	srcyoff := float32(srcb.Dy())/2 - 0.5
	dstxoff := float32(w)/2 - 0.5
	dstyoff := float32(h)/2 - 0.5

	bgpx := pixelclr(p.bgcolor)
	asin, acos := sincosf32(p.angle)

	pixGetter := newPixelGetter(src)
	pixSetter := newPixelSetter(dst)

	parallelize(options.Parallelization, 0, h, func(pmin, pmax int) {
		for y := pmin; y < pmax; y++ {
			for x := 0; x < w; x++ {

				xf, yf := rotatePoint(float32(x)-dstxoff, float32(y)-dstyoff, asin, acos)
				xf, yf = float32(srcb.Min.X)+xf+srcxoff, float32(srcb.Min.Y)+yf+srcyoff

				switch p.interpolation {
				case CubicInterpolation:
					var calc bool
					var pxs [16]pixel
					var cfs [16]float32
					var px pixel

					x0, y0 := int(floorf32(xf)), int(floorf32(yf))
					xq, yq := xf-float32(x0), yf-float32(y0)

					for i := 0; i < 4; i++ {
						for j := 0; j < 4; j++ {
							pt := image.Pt(x0+j-1, y0+i-1)
							if pt.In(srcb) {
								pxs[i*4+j] = pixGetter.getPixel(pt.X, pt.Y)
								calc = true
							} else {
								pxs[i*4+j] = bgpx
							}
						}
					}

					if !calc {
						pixSetter.setPixel(dstb.Min.X+x, dstb.Min.Y+y, bgpx)
						continue
					}

					cfs[0] = (1.0 / 36.0) * xq * yq * (xq - 1) * (xq - 2) * (yq - 1) * (yq - 2)
					cfs[1] = -(1.0 / 12.0) * yq * (xq - 1) * (xq - 2) * (xq + 1) * (yq - 1) * (yq - 2)
					cfs[2] = (1.0 / 12.0) * xq * yq * (xq + 1) * (xq - 2) * (yq - 1) * (yq - 2)
					cfs[3] = -(1.0 / 36.0) * xq * yq * (xq - 1) * (xq + 1) * (yq - 1) * (yq - 2)
					cfs[4] = -(1.0 / 12.0) * xq * (xq - 1) * (xq - 2) * (yq - 1) * (yq - 2) * (yq + 1)
					cfs[5] = 0.25 * (xq - 1) * (xq - 2) * (xq + 1) * (yq - 1) * (yq - 2) * (yq + 1)
					cfs[6] = -0.25 * xq * (xq + 1) * (xq - 2) * (yq - 1) * (yq - 2) * (yq + 1)
					cfs[7] = (1.0 / 12.0) * xq * (xq - 1) * (xq + 1) * (yq - 1) * (yq - 2) * (yq + 1)
					cfs[8] = (1.0 / 12.0) * xq * yq * (xq - 1) * (xq - 2) * (yq + 1) * (yq - 2)
					cfs[9] = -0.25 * yq * (xq - 1) * (xq - 2) * (xq + 1) * (yq + 1) * (yq - 2)
					cfs[10] = 0.25 * xq * yq * (xq + 1) * (xq - 2) * (yq + 1) * (yq - 2)
					cfs[11] = -(1.0 / 12.0) * xq * yq * (xq - 1) * (xq + 1) * (yq + 1) * (yq - 2)
					cfs[12] = -(1.0 / 36.0) * xq * yq * (xq - 1) * (xq - 2) * (yq - 1) * (yq + 1)
					cfs[13] = (1.0 / 12.0) * yq * (xq - 1) * (xq - 2) * (xq + 1) * (yq - 1) * (yq + 1)
					cfs[14] = -(1.0 / 12.0) * xq * yq * (xq + 1) * (xq - 2) * (yq - 1) * (yq + 1)
					cfs[15] = (1.0 / 36.0) * xq * yq * (xq - 1) * (xq + 1) * (yq - 1) * (yq + 1)

					for i := range pxs {
						px.R += pxs[i].R * cfs[i]
						px.G += pxs[i].G * cfs[i]
						px.B += pxs[i].B * cfs[i]
						px.A += pxs[i].A * cfs[i]
					}

					pixSetter.setPixel(dstb.Min.X+x, dstb.Min.Y+y, px)

				case LinearInterpolation:
					var calc bool
					var pxs [4]pixel
					var cfs [4]float32
					var px pixel

					x0, y0 := int(floorf32(xf)), int(floorf32(yf))
					xq, yq := xf-float32(x0), yf-float32(y0)

					for i := 0; i < 2; i++ {
						for j := 0; j < 2; j++ {
							pt := image.Pt(x0+j, y0+i)
							if pt.In(srcb) {
								pxs[i*2+j] = pixGetter.getPixel(pt.X, pt.Y)
								calc = true
							} else {
								pxs[i*2+j] = bgpx
							}
						}
					}

					if !calc {
						pixSetter.setPixel(dstb.Min.X+x, dstb.Min.Y+y, bgpx)
						continue
					}

					cfs[0] = (1 - xq) * (1 - yq)
					cfs[1] = xq * (1 - yq)
					cfs[2] = (1 - xq) * yq
					cfs[3] = xq * yq

					for i := range pxs {
						px.R += pxs[i].R * cfs[i]
						px.G += pxs[i].G * cfs[i]
						px.B += pxs[i].B * cfs[i]
						px.A += pxs[i].A * cfs[i]
					}

					pixSetter.setPixel(dstb.Min.X+x, dstb.Min.Y+y, px)

				default:
					var px pixel
					x0, y0 := int(floorf32(xf+0.5)), int(floorf32(yf+0.5))
					if image.Pt(x0, y0).In(srcb) {
						px = pixGetter.getPixel(x0, y0)
					} else {
						px = bgpx
					}
					pixSetter.setPixel(dstb.Min.X+x, dstb.Min.Y+y, px)
				}
			}
		}
	})

	return
}

// Rotate creates a filter that rotates an image by the given angle counter-clockwise.
// The angle parameter is the rotation angle in degrees.
// The backgroundColor parameter specifies the color of the uncovered zone after the rotation.
// The interpolation parameter specifies the interpolation method.
// Supported interpolation methods: NearestNeighborInterpolation, LinearInterpolation, CubicInterpolation.
//
// Example:
//
//	g := gift.New(
//		gift.Rotate(45, color.Black, gift.LinearInterpolation),
//	)
//	dst := image.NewRGBA(g.Bounds(src.Bounds()))
//	g.Draw(dst, src)
//
func Rotate(angle float32, backgroundColor color.Color, interpolation Interpolation) Filter {
	return &rotateFilter{
		angle:         angle,
		bgcolor:       backgroundColor,
		interpolation: interpolation,
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

/*

Package gift provides a set of useful image processing filters.

Basic usage:

	// 1. Create a new GIFT and add some filters:

	g := gift.New(
	    gift.Resize(800, 0, gift.LanczosResampling),
	    gift.UnsharpMask(1.0, 1.0, 0.0),
	)

	// 2. Create a new image of the corresponding size.
	// dst is a new target image, src is the original image

	dst := image.NewRGBA(g.Bounds(src.Bounds()))

	// 3. Use Draw func to apply the filters to src and store the result in dst:

	g.Draw(dst, src)

*/
package gift

import (
	"image"
	"image/draw"
)

type Filter interface {
	// Draw applies the filter to the src image and outputs the result to the dst image.
	Draw(dst draw.Image, src image.Image, options *Options)
	// Bounds calculates the appropriate bounds of an image after applying the filter.
	Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle)
}

type Options struct {
	Parallelization bool
}

var defaultOptions = Options{
	Parallelization: true,
}

// GIFT implements a list of filters that can be applied to an image at once.
type GIFT struct {
	Filters []Filter
	Options Options
}

// New creates a new instance of the filter toolkit and initializes it with the given list of filters.
func New(filters ...Filter) *GIFT {
	return &GIFT{
		Filters: filters,
		Options: defaultOptions,
	}
}

// SetParallelization enables or disables faster image processing using parallel goroutines.
// Parallelization is enabled by default.
// To achieve maximum performance, make sure to allow Go to utilize all CPU cores:
//
// 	runtime.GOMAXPROCS(runtime.NumCPU())
//
func (g *GIFT) SetParallelization(isEnabled bool) {
	g.Options.Parallelization = isEnabled
}

// Parallelization returns the current state of parallelization option.
func (g *GIFT) Parallelization() bool {
	return g.Options.Parallelization
}

// Add appends the given filters to the list of filters.
func (g *GIFT) Add(filters ...Filter) {
	g.Filters = append(g.Filters, filters...)
}

// Empty removes all the filters from the list.
func (g *GIFT) Empty() {
	g.Filters = []Filter{}
}

// Bounds calculates the appropriate bounds for the result image after applying all the added filters.
// Parameter srcBounds is the bounds of the source image.
//
// Example:
//
// 	src := image.NewRGBA(image.Rect(0, 0, 100, 200))
//	g := gift.New(gift.Rotate90())
//
// 	// calculate image bounds after applying rotation and create a new image of that size.
// 	dst := image.NewRGBA(g.Bounds(src.Bounds())) // dst bounds: (0, 0, 200, 100)
//
//
func (g *GIFT) Bounds(srcBounds image.Rectangle) (dstBounds image.Rectangle) {
	b := srcBounds
	for _, f := range g.Filters {
		b = f.Bounds(b)
	}
	dstBounds = b
	return
}

// Draw applies all the added filters to the src image and outputs the result to the dst image.
func (g *GIFT) Draw(dst draw.Image, src image.Image) {
	if len(g.Filters) == 0 {
		copyimage(dst, src, &g.Options)
		return
	}

	first, last := 0, len(g.Filters)-1
	var tmpIn image.Image
	var tmpOut draw.Image

	for i, f := range g.Filters {
		if i == first {
			tmpIn = src
		} else {
			tmpIn = tmpOut
		}

		if i == last {
			tmpOut = dst
		} else {
			tmpOut = createTempImage(f.Bounds(tmpIn.Bounds()))
		}

		f.Draw(tmpOut, tmpIn, &g.Options)
	}
}

# GO IMAGE FILTERING TOOLKIT (GIFT)

[![GoDoc](https://godoc.org/github.com/disintegration/gift?status.svg)](https://godoc.org/github.com/disintegration/gift)
[![Build Status](https://travis-ci.org/disintegration/gift.svg?branch=master)](https://travis-ci.org/disintegration/gift)
[![Coverage Status](https://coveralls.io/repos/github/disintegration/gift/badge.svg?branch=master)](https://coveralls.io/github/disintegration/gift?branch=master)

*Package gift provides a set of useful image processing filters.*

Pure Go. No external dependencies outside of the Go standard library.


### INSTALLATION / UPDATING

    go get -u github.com/disintegration/gift
  


### DOCUMENTATION

http://godoc.org/github.com/disintegration/gift
  


### QUICK START

```go
// 1. Create a new GIFT filter list and add some filters:
g := gift.New(
    gift.Resize(800, 0, gift.LanczosResampling),
    gift.UnsharpMask(1.0, 1.0, 0.0),
)

// 2. Create a new image of the corresponding size.
// dst is a new target image, src is the original image
dst := image.NewRGBA(g.Bounds(src.Bounds()))

// 3. Use Draw func to apply the filters to src and store the result in dst:
g.Draw(dst, src)
``` 

### USAGE

To create a sequence of filters, the `New` function is used:
 ```go
g := gift.New(
    gift.Grayscale(),
    gift.Contrast(10),
)
 ```
Filters also can be added using the `Add` method:
 ```go
g.Add(GaussianBlur(2)) 
```

The `Bounds` method takes the bounds of the source image and returns appropriate bounds for the destination image to fit the result (for example, after using `Resize` or `Rotate` filters). 

```go
dst := image.NewRGBA(g.Bounds(src.Bounds()))
```

There are two methods available to apply these filters to an image:

- `Draw` applies all the added filters to the src image and outputs the result to the dst image starting from the top-left corner (Min point).
 ```go
 g.Draw(dst, src)
 ```

- `DrawAt` provides more control. It outputs the filtered src image to the dst image at the specified position using the specified image composition operator. This example is equivalent to the previous:
 ```go
 g.DrawAt(dst, src, dst.Bounds().Min, gift.CopyOperator)
 ```

Two image composition operators are supported by now:
- `CopyOperator` - Replaces pixels of the dst image with pixels of the filtered src image. This mode is used by the Draw method.
- `OverOperator` - Places the filtered src image on top of the dst image. This mode makes sence if the filtered src image has transparent areas.

Empty filter list can be used to create a copy of an image or to paste one image to another. For example:
```go
// Create a new image with dimensions of bgImage
dstImage := image.NewNRGBA(bgImage.Bounds())
// Copy the bgImage to the dstImage.
gift.New().Draw(dstImage, bgImage)
// Draw the fgImage over the dstImage at the (100, 100) position.
gift.New().DrawAt(dstImage, fgImage, image.Pt(100, 100), gift.OverOperator)
```


### SUPPORTED FILTERS

+ Transformations

    - Crop(rect image.Rectangle)
    - CropToSize(width, height int, anchor Anchor)
    - FlipHorizontal()
    - FlipVertical()
    - Resize(width, height int, resampling Resampling)
    - ResizeToFill(width, height int, resampling Resampling, anchor Anchor)
    - ResizeToFit(width, height int, resampling Resampling)
    - Rotate(angle float32, backgroundColor color.Color, interpolation Interpolation)
    - Rotate180()
    - Rotate270()
    - Rotate90()
    - Transpose()
    - Transverse()
    
+ Adjustments & effects

    - Brightness(percentage float32)
    - ColorBalance(percentageRed, percentageGreen, percentageBlue float32)
    - ColorFunc(fn func(r0, g0, b0, a0 float32) (r, g, b, a float32))
    - Colorize(hue, saturation, percentage float32)
    - ColorspaceLinearToSRGB()
    - ColorspaceSRGBToLinear()
    - Contrast(percentage float32)
    - Convolution(kernel []float32, normalize, alpha, abs bool, delta float32)
    - Gamma(gamma float32)
    - GaussianBlur(sigma float32)
    - Grayscale()
    - Hue(shift float32)
    - Invert()
    - Maximum(ksize int, disk bool)
    - Mean(ksize int, disk bool)
    - Median(ksize int, disk bool)
    - Minimum(ksize int, disk bool)
    - Pixelate(size int)
    - Saturation(percentage float32)
    - Sepia(percentage float32)
    - Sigmoid(midpoint, factor float32)
    - Sobel()
    - UnsharpMask(sigma, amount, thresold float32)


### FILTER EXAMPLES

##### Resize using lanczos resampling
```go
gift.Resize(200, 0, gift.LanczosResampling)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_resize_lanczos.jpg)

##### Resize using linear resampling
```go
gift.Resize(200, 0, gift.LinearResampling)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_resize_linear.jpg)

##### Resize to fit 160x160px bounding box
```go
gift.ResizeToFit(160, 160, gift.LanczosResampling)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original_h.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_resize_fit.jpg)

##### Resize to fill 160x160px rectangle, anchor: center
```go
gift.ResizeToFill(160, 160, gift.LanczosResampling, gift.CenterAnchor)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original_h.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_resize_fill.jpg)

##### Crop 90, 90 - 250, 250
```go
gift.Crop(image.Rect(90, 90, 250, 250))
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_crop.jpg)

##### Crop to size 160x160px, anchor: center
```go
gift.CropToSize(160, 160, gift.CenterAnchor)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_crop_to_size.jpg)

##### Rotate 90 degrees
```go
gift.Rotate90()
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_rotate_90.jpg)

##### Rotate 180 degrees
```go
gift.Rotate180()
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_rotate_180.jpg)

##### Rotate 270 degrees
```go
gift.Rotate270()
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_rotate_270.jpg)

##### Rotate 30 degrees, white background, cubic interpolation
```go
gift.Rotate(30, color.White, gift.CubicInterpolation)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original_small.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_rotate_30.jpg)

##### Flip horizontal
```go
gift.FlipHorizontal()
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_flip_h.jpg)

##### Flip vertical
```go
gift.FlipVertical()
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_flip_v.jpg)

##### Transpose
```go
gift.Transpose()
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_transpose.jpg)

##### Transverse
```go
gift.Transverse()
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_transverse.jpg)

##### Contrast +30%
```go
gift.Contrast(30)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_contrast_30.jpg)

##### Contrast -30%
```go
gift.Contrast(-30)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_contrast_-30.jpg)

##### Brightness +30%
```go
gift.Brightness(30)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_brightness_30.jpg)

##### Brightness -30%
```go
gift.Brightness(-30)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_brightness_-30.jpg)

##### Saturation +50%
```go
gift.Saturation(50)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_saturation_50.jpg)

##### Saturation -50%
```go
gift.Saturation(-50)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_saturation_-50.jpg)

##### Hue +45
```go
gift.Hue(45)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_hue_45.jpg)

##### Hue -45
```go
gift.Hue(-45)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_hue_-45.jpg)

##### Sigmoid 0.5, 5.0
```go
gift.Sigmoid(0.5, 5.0)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_sigmoid.jpg)

##### Gamma correction = 0.5
```go
gift.Gamma(0.5)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_gamma_0.5.jpg)

##### Gamma correction = 1.5
```go
gift.Gamma(1.5)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_gamma_1.5.jpg)

##### Gaussian blur, sigma=1.0
```go
gift.GaussianBlur(1.0)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_gaussian_blur.jpg)

##### Unsharp mask, sigma=1.0, amount=1.5, thresold=0.0
```go
gift.UnsharpMask(1.0, 1.5, 0.0)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_unsharp_mask.jpg)

##### Grayscale
```go
gift.Grayscale()
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_grayscale.jpg)

##### Sepia, 100%
```go
gift.Sepia(100)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_sepia.jpg)

##### Invert
```go
gift.Invert()
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_invert.jpg)

##### Colorize, blue, saturation=50%
```go
gift.Colorize(240, 50, 100)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_colorize.jpg)

##### Color balance, +20% red, -20% green
```go
gift.ColorBalance(20, -20, 0)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_colorbalance.jpg)

##### Color function
```go
gift.ColorFunc(
    func(r0, g0, b0, a0 float32) (r, g, b, a float32) {
        r = 1 - r0   // invert the red channel
        g = g0 + 0.1 // shift the green channel by 0.1
        b = 0        // set the blue channel to 0
        a = a0       // preserve the alpha channel
        return
    },
)
 ```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_colorfunc.jpg)

##### Local mean, disc shape, size=5
```go
gift.Mean(5, true)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_mean.jpg)

##### Local median, disc shape, size=5
```go
gift.Median(5, true)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_median.jpg)

##### Local minimum, disc shape, size=5
```go
gift.Minimum(5, true)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_minimum.jpg)

##### Local maximum, disc shape, size=5
```go
gift.Maximum(5, true)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_maximum.jpg)

##### Pixelate, size=5
```go
gift.Pixelate(5)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_pixelate.jpg)

##### Convolution matrix - Emboss
```go
gift.Convolution(
    []float32{
        -1, -1, 0,
        -1, 1, 1,
        0, 1, 1,
    },
    false, false, false, 0.0,
)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_convolution_emboss.jpg)

##### Convolution matrix - Edge detection
```go
gift.Convolution(
    []float32{
        -1, -1, -1,
        -1, 8, -1,
        -1, -1, -1,
    },
    false, false, false, 0.0,
)
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original2.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_edge.jpg)

##### Sobel operator
```go
gift.Sobel()
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original2.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_sobel.jpg)


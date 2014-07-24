# GO IMAGE FILTERING TOOLKIT (GIFT)

*Package gift provides a set of useful image processing filters.*

Pure Go. No external dependencies outside of the Go standard library.

*BETA*

### INSTALLATION / UPDATING

    go get -u github.com/disintegration/gift
  


### DOCUMENTATION

http://godoc.org/github.com/disintegration/gift
  


### QUICK START

```go
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
``` 


### SUPPORTED FILTERS

+ Transformations

    - Crop(rect image.Rectangle)
    - FlipHorizontal()
    - FlipVertical()
    - Resize(width, height int, resampling Resampling)
    - Rotate180()
    - Rotate270()
    - Rotate90()
    - Transpose()
    - Transverse()
    
+ Adjustments & effects

    - Brightness(percentage float32)
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
    - Saturation(percentage float32)
    - Sepia(percentage float32)
    - Sigmoid(midpoint, factor float32)
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

##### Crop 90, 90 - 250, 250
```go
gift.Crop(image.Rect(90, 90, 250, 250))
```
Original image | Filtered image
--- | ---
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_crop.jpg)

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
![original](http://disintegration.github.io/gift/examples/original.jpg) | ![filtered](http://disintegration.github.io/gift/examples/example_convolution_edge.jpg)


# GO IMAGE FILTERING TOOLKIT (GIFT)

*Package gift provides a set of useful image processing filters.*


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
    - InvertColors()
    - Maximum(ksize int, disk bool)
    - Mean(ksize int, disk bool)
    - Median(ksize int, disk bool)
    - Minimum(ksize int, disk bool)
    - Sigmoid(midpoint, factor float32)
    - UnsharpMask(sigma, amount, thresold float32)


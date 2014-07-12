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
  gift.AdjustGamma(0.9),
)

// 2. Create a new image of the corresponding size.
// dst is a new target image, src is the original image
dst := image.NewRGBA(g.Bounds(src.Bounds()))

// 3. Use Draw func to apply the filters to src and store the result in dst:
g.Draw(dst, src)
``` 


### SUPPORTED FILTERS

+ Transformations

    - Resize(width, height int, resampling Resampling)
    - Rotate180()
    - Rotate270()
    - Rotate90()
    - FlipHorizontally()
    - FlipVertically()
    - Transpose()
    
+ Effects & color modifications

    - InvertColors()
    - AdjustGamma(gamma float32)
    - GaussianBlur(sigma float32)
    - UnsharpMask(sigma, amount, thresold float32)
    - Convolution(kernel []float32, normalize, alpha, abs bool, delta float32)
    - ColorspaceLinearToSRGB()
    - ColorspaceSRGBToLinear()
package gift

import (
	"image"
	"image/color"
	"image/draw"
)

type pixel struct {
	R, G, B, A float32
}

type imageType int

const (
	itGeneric imageType = iota
	itNRGBA
	itNRGBA64
	itRGBA
	itRGBA64
	itYCbCr
	itGray
	itGray16
	itPaletted
)

type pixelGetter struct {
	it       imageType
	bounds   image.Rectangle
	image    image.Image
	nrgba    *image.NRGBA
	nrgba64  *image.NRGBA64
	rgba     *image.RGBA
	rgba64   *image.RGBA64
	gray     *image.Gray
	gray16   *image.Gray16
	ycbcr    *image.YCbCr
	paletted *image.Paletted
	palette  []pixel
}

func newPixelGetter(img image.Image) *pixelGetter {
	switch img := img.(type) {
	case *image.NRGBA:
		return &pixelGetter{
			it:     itNRGBA,
			bounds: img.Bounds(),
			nrgba:  img,
		}

	case *image.NRGBA64:
		return &pixelGetter{
			it:      itNRGBA64,
			bounds:  img.Bounds(),
			nrgba64: img,
		}

	case *image.RGBA:
		return &pixelGetter{
			it:     itRGBA,
			bounds: img.Bounds(),
			rgba:   img,
		}

	case *image.RGBA64:
		return &pixelGetter{
			it:     itRGBA64,
			bounds: img.Bounds(),
			rgba64: img,
		}

	case *image.Gray:
		return &pixelGetter{
			it:     itGray,
			bounds: img.Bounds(),
			gray:   img,
		}

	case *image.Gray16:
		return &pixelGetter{
			it:     itGray16,
			bounds: img.Bounds(),
			gray16: img,
		}

	case *image.YCbCr:
		return &pixelGetter{
			it:     itYCbCr,
			bounds: img.Bounds(),
			ycbcr:  img,
		}

	case *image.Paletted:
		return &pixelGetter{
			it:       itPaletted,
			bounds:   img.Bounds(),
			paletted: img,
			palette:  convertPalette(img.Palette),
		}

	default:
		return &pixelGetter{
			it:     itGeneric,
			bounds: img.Bounds(),
			image:  img,
		}
	}
}

const (
	qf8  = 1 / 255.0
	qf16 = 1 / 65535.0
	epal = qf16 * qf16 / 2
)

func pixelFromColor(c color.Color) (px pixel) {
	r16, g16, b16, a16 := c.RGBA()
	switch a16 {
	case 0:
		px = pixel{0, 0, 0, 0}
	case 65535:
		r := float32(r16) * qf16
		g := float32(g16) * qf16
		b := float32(b16) * qf16
		px = pixel{r, g, b, 1}
	default:
		q := float32(1) / float32(a16)
		r := float32(r16) * q
		g := float32(g16) * q
		b := float32(b16) * q
		a := float32(a16) * qf16
		px = pixel{r, g, b, a}
	}
	return px
}

func convertPalette(p []color.Color) []pixel {
	pal := make([]pixel, len(p))
	for i := 0; i < len(p); i++ {
		pal[i] = pixelFromColor(p[i])
	}
	return pal
}

func getPaletteIndex(pal []pixel, px pixel) int {
	var k int
	var dmin float32 = 4
	for i, palpx := range pal {
		d := px.R - palpx.R
		dcur := d * d
		d = px.G - palpx.G
		dcur += d * d
		d = px.B - palpx.B
		dcur += d * d
		d = px.A - palpx.A
		dcur += d * d
		if dcur < epal {
			return i
		}
		if dcur < dmin {
			dmin = dcur
			k = i
		}
	}
	return k
}

func (p *pixelGetter) getPixel(x, y int) (px pixel) {
	switch p.it {
	case itNRGBA:
		i := p.nrgba.PixOffset(x, y)
		r := float32(p.nrgba.Pix[i+0]) * qf8
		g := float32(p.nrgba.Pix[i+1]) * qf8
		b := float32(p.nrgba.Pix[i+2]) * qf8
		a := float32(p.nrgba.Pix[i+3]) * qf8
		px = pixel{r, g, b, a}

	case itNRGBA64:
		i := p.nrgba64.PixOffset(x, y)
		r := float32(uint16(p.nrgba64.Pix[i+0])<<8|uint16(p.nrgba64.Pix[i+1])) * qf16
		g := float32(uint16(p.nrgba64.Pix[i+2])<<8|uint16(p.nrgba64.Pix[i+3])) * qf16
		b := float32(uint16(p.nrgba64.Pix[i+4])<<8|uint16(p.nrgba64.Pix[i+5])) * qf16
		a := float32(uint16(p.nrgba64.Pix[i+6])<<8|uint16(p.nrgba64.Pix[i+7])) * qf16
		px = pixel{r, g, b, a}

	case itRGBA:
		i := p.rgba.PixOffset(x, y)
		a8 := p.rgba.Pix[i+3]
		switch a8 {
		case 0:
			px = pixel{0, 0, 0, 0}
		case 255:
			r := float32(p.rgba.Pix[i+0]) * qf8
			g := float32(p.rgba.Pix[i+1]) * qf8
			b := float32(p.rgba.Pix[i+2]) * qf8
			px = pixel{r, g, b, 1}
		default:
			q := float32(1) / float32(a8)
			r := float32(p.rgba.Pix[i+0]) * q
			g := float32(p.rgba.Pix[i+1]) * q
			b := float32(p.rgba.Pix[i+2]) * q
			a := float32(a8) * qf8
			px = pixel{r, g, b, a}
		}

	case itRGBA64:
		i := p.rgba64.PixOffset(x, y)
		a16 := uint16(p.rgba64.Pix[i+6])<<8 | uint16(p.rgba64.Pix[i+7])
		switch a16 {
		case 0:
			px = pixel{0, 0, 0, 0}
		case 65535:
			r := float32(uint16(p.rgba64.Pix[i+0])<<8|uint16(p.rgba64.Pix[i+1])) * qf16
			g := float32(uint16(p.rgba64.Pix[i+2])<<8|uint16(p.rgba64.Pix[i+3])) * qf16
			b := float32(uint16(p.rgba64.Pix[i+4])<<8|uint16(p.rgba64.Pix[i+5])) * qf16
			px = pixel{r, g, b, 1}
		default:
			q := float32(1) / float32(a16)
			r := float32(uint16(p.rgba64.Pix[i+0])<<8|uint16(p.rgba64.Pix[i+1])) * q
			g := float32(uint16(p.rgba64.Pix[i+2])<<8|uint16(p.rgba64.Pix[i+3])) * q
			b := float32(uint16(p.rgba64.Pix[i+4])<<8|uint16(p.rgba64.Pix[i+5])) * q
			a := float32(a16) * qf16
			px = pixel{r, g, b, a}
		}

	case itGray:
		i := p.gray.PixOffset(x, y)
		v := float32(p.gray.Pix[i]) * qf8
		px = pixel{v, v, v, 1}

	case itGray16:
		i := p.gray16.PixOffset(x, y)
		v := float32(uint16(p.gray16.Pix[i+0])<<8|uint16(p.gray16.Pix[i+1])) * qf16
		px = pixel{v, v, v, 1}

	case itYCbCr:
		iy := p.ycbcr.YOffset(x, y)
		ic := p.ycbcr.COffset(x, y)
		r8, g8, b8 := color.YCbCrToRGB(p.ycbcr.Y[iy], p.ycbcr.Cb[ic], p.ycbcr.Cr[ic])
		r := float32(r8) * qf8
		g := float32(g8) * qf8
		b := float32(b8) * qf8
		px = pixel{r, g, b, 1}

	case itPaletted:
		i := p.paletted.PixOffset(x, y)
		k := p.paletted.Pix[i]
		px = p.palette[k]

	case itGeneric:
		px = pixelFromColor(p.image.At(x, y))
	}
	return
}

func (p *pixelGetter) getPixelRow(y int, buf *[]pixel) {
	*buf = (*buf)[0:0]
	for x := p.bounds.Min.X; x != p.bounds.Max.X; x++ {
		*buf = append(*buf, p.getPixel(x, y))
	}
}

func (p *pixelGetter) getPixelColumn(x int, buf *[]pixel) {
	*buf = (*buf)[0:0]
	for y := p.bounds.Min.Y; y != p.bounds.Max.Y; y++ {
		*buf = append(*buf, p.getPixel(x, y))
	}
}

func f32u8(val float32) uint8 {
	x := int64(val + 0.5)
	if x > 255 {
		return 255
	}
	if x > 0 {
		return uint8(x)
	}
	return 0
}

func f32u16(val float32) uint16 {
	x := int64(val + 0.5)
	if x > 65535 {
		return 65535
	}
	if x > 0 {
		return uint16(x)
	}
	return 0
}

type pixelSetter struct {
	it       imageType
	bounds   image.Rectangle
	image    draw.Image
	nrgba    *image.NRGBA
	nrgba64  *image.NRGBA64
	rgba     *image.RGBA
	rgba64   *image.RGBA64
	gray     *image.Gray
	gray16   *image.Gray16
	paletted *image.Paletted
	palette  []pixel
}

func newPixelSetter(img draw.Image) *pixelSetter {
	switch img := img.(type) {
	case *image.NRGBA:
		return &pixelSetter{
			it:     itNRGBA,
			bounds: img.Bounds(),
			nrgba:  img,
		}

	case *image.NRGBA64:
		return &pixelSetter{
			it:      itNRGBA64,
			bounds:  img.Bounds(),
			nrgba64: img,
		}

	case *image.RGBA:
		return &pixelSetter{
			it:     itRGBA,
			bounds: img.Bounds(),
			rgba:   img,
		}

	case *image.RGBA64:
		return &pixelSetter{
			it:     itRGBA64,
			bounds: img.Bounds(),
			rgba64: img,
		}

	case *image.Gray:
		return &pixelSetter{
			it:     itGray,
			bounds: img.Bounds(),
			gray:   img,
		}

	case *image.Gray16:
		return &pixelSetter{
			it:     itGray16,
			bounds: img.Bounds(),
			gray16: img,
		}

	case *image.Paletted:
		return &pixelSetter{
			it:       itPaletted,
			bounds:   img.Bounds(),
			paletted: img,
			palette:  convertPalette(img.Palette),
		}

	default:
		return &pixelSetter{
			it:     itGeneric,
			bounds: img.Bounds(),
			image:  img,
		}
	}
}

func (p *pixelSetter) setPixel(x, y int, px pixel) {
	if !image.Pt(x, y).In(p.bounds) {
		return
	}
	switch p.it {
	case itNRGBA:
		i := p.nrgba.PixOffset(x, y)
		p.nrgba.Pix[i+0] = f32u8(px.R * 255)
		p.nrgba.Pix[i+1] = f32u8(px.G * 255)
		p.nrgba.Pix[i+2] = f32u8(px.B * 255)
		p.nrgba.Pix[i+3] = f32u8(px.A * 255)

	case itNRGBA64:
		r16 := f32u16(px.R * 65535)
		g16 := f32u16(px.G * 65535)
		b16 := f32u16(px.B * 65535)
		a16 := f32u16(px.A * 65535)
		i := p.nrgba64.PixOffset(x, y)
		p.nrgba64.Pix[i+0] = uint8(r16 >> 8)
		p.nrgba64.Pix[i+1] = uint8(r16 & 0xff)
		p.nrgba64.Pix[i+2] = uint8(g16 >> 8)
		p.nrgba64.Pix[i+3] = uint8(g16 & 0xff)
		p.nrgba64.Pix[i+4] = uint8(b16 >> 8)
		p.nrgba64.Pix[i+5] = uint8(b16 & 0xff)
		p.nrgba64.Pix[i+6] = uint8(a16 >> 8)
		p.nrgba64.Pix[i+7] = uint8(a16 & 0xff)

	case itRGBA:
		fa := px.A * 255
		i := p.rgba.PixOffset(x, y)
		p.rgba.Pix[i+0] = f32u8(px.R * fa)
		p.rgba.Pix[i+1] = f32u8(px.G * fa)
		p.rgba.Pix[i+2] = f32u8(px.B * fa)
		p.rgba.Pix[i+3] = f32u8(fa)

	case itRGBA64:
		fa := px.A * 65535
		r16 := f32u16(px.R * fa)
		g16 := f32u16(px.G * fa)
		b16 := f32u16(px.B * fa)
		a16 := f32u16(fa)
		i := p.rgba64.PixOffset(x, y)
		p.rgba64.Pix[i+0] = uint8(r16 >> 8)
		p.rgba64.Pix[i+1] = uint8(r16 & 0xff)
		p.rgba64.Pix[i+2] = uint8(g16 >> 8)
		p.rgba64.Pix[i+3] = uint8(g16 & 0xff)
		p.rgba64.Pix[i+4] = uint8(b16 >> 8)
		p.rgba64.Pix[i+5] = uint8(b16 & 0xff)
		p.rgba64.Pix[i+6] = uint8(a16 >> 8)
		p.rgba64.Pix[i+7] = uint8(a16 & 0xff)

	case itGray:
		i := p.gray.PixOffset(x, y)
		p.gray.Pix[i] = f32u8((0.299*px.R + 0.587*px.G + 0.114*px.B) * px.A * 255)

	case itGray16:
		i := p.gray16.PixOffset(x, y)
		y16 := f32u16((0.299*px.R + 0.587*px.G + 0.114*px.B) * px.A * 65535)
		p.gray16.Pix[i+0] = uint8(y16 >> 8)
		p.gray16.Pix[i+1] = uint8(y16 & 0xff)

	case itPaletted:
		px1 := pixel{
			minf32(maxf32(px.R, 0), 1),
			minf32(maxf32(px.G, 0), 1),
			minf32(maxf32(px.B, 0), 1),
			minf32(maxf32(px.A, 0), 1),
		}
		i := p.paletted.PixOffset(x, y)
		k := getPaletteIndex(p.palette, px1)
		p.paletted.Pix[i] = uint8(k)

	case itGeneric:
		r16 := f32u16(px.R * 65535)
		g16 := f32u16(px.G * 65535)
		b16 := f32u16(px.B * 65535)
		a16 := f32u16(px.A * 65535)
		p.image.Set(x, y, color.NRGBA64{r16, g16, b16, a16})
	}
}

func (p *pixelSetter) setPixelRow(y int, buf []pixel) {
	for i, x := 0, p.bounds.Min.X; i < len(buf); i, x = i+1, x+1 {
		p.setPixel(x, y, buf[i])
	}
}

func (p *pixelSetter) setPixelColumn(x int, buf []pixel) {
	for i, y := 0, p.bounds.Min.Y; i < len(buf); i, y = i+1, y+1 {
		p.setPixel(x, y, buf[i])
	}
}

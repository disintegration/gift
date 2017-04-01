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
	imgType     imageType
	imgBounds   image.Rectangle
	imgGeneric  image.Image
	imgNRGBA    *image.NRGBA
	imgNRGBA64  *image.NRGBA64
	imgRGBA     *image.RGBA
	imgRGBA64   *image.RGBA64
	imgYCbCr    *image.YCbCr
	imgGray     *image.Gray
	imgGray16   *image.Gray16
	imgPaletted *image.Paletted
	imgPalette  []pixel
}

func newPixelGetter(img image.Image) (p *pixelGetter) {
	switch img := img.(type) {
	case *image.NRGBA:
		p = &pixelGetter{
			imgType:   itNRGBA,
			imgBounds: img.Bounds(),
			imgNRGBA:  img,
		}

	case *image.NRGBA64:
		p = &pixelGetter{
			imgType:    itNRGBA64,
			imgBounds:  img.Bounds(),
			imgNRGBA64: img,
		}

	case *image.RGBA:
		p = &pixelGetter{
			imgType:   itRGBA,
			imgBounds: img.Bounds(),
			imgRGBA:   img,
		}

	case *image.RGBA64:
		p = &pixelGetter{
			imgType:   itRGBA64,
			imgBounds: img.Bounds(),
			imgRGBA64: img,
		}

	case *image.Gray:
		p = &pixelGetter{
			imgType:   itGray,
			imgBounds: img.Bounds(),
			imgGray:   img,
		}

	case *image.Gray16:
		p = &pixelGetter{
			imgType:   itGray16,
			imgBounds: img.Bounds(),
			imgGray16: img,
		}

	case *image.YCbCr:
		p = &pixelGetter{
			imgType:   itYCbCr,
			imgBounds: img.Bounds(),
			imgYCbCr:  img,
		}

	case *image.Paletted:
		p = &pixelGetter{
			imgType:     itPaletted,
			imgBounds:   img.Bounds(),
			imgPaletted: img,
			imgPalette:  convertPalette(img.Palette),
		}
		return

	default:
		p = &pixelGetter{
			imgType:    itGeneric,
			imgBounds:  img.Bounds(),
			imgGeneric: img,
		}
	}
	return
}

const (
	qf8  = 1 / 255.0
	qf16 = 1 / 65535.0
	epal = qf16 * qf16 / 2
)

func convertPalette(p []color.Color) []pixel {
	plen := len(p)
	pnew := make([]pixel, plen)
	for i := 0; i < plen; i++ {
		r16, g16, b16, a16 := p[i].RGBA()
		switch a16 {
		case 0:
			pnew[i] = pixel{0, 0, 0, 0}
		case 65535:
			r := float32(r16) * qf16
			g := float32(g16) * qf16
			b := float32(b16) * qf16
			pnew[i] = pixel{r, g, b, 1}
		default:
			q := float32(1) / float32(a16)
			r := float32(r16) * q
			g := float32(g16) * q
			b := float32(b16) * q
			a := float32(a16) * qf16
			pnew[i] = pixel{r, g, b, a}
		}
	}
	return pnew
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

func pixelclr(c color.Color) (px pixel) {
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

func (p *pixelGetter) getPixel(x, y int) (px pixel) {
	switch p.imgType {
	case itNRGBA:
		i := p.imgNRGBA.PixOffset(x, y)
		r := float32(p.imgNRGBA.Pix[i+0]) * qf8
		g := float32(p.imgNRGBA.Pix[i+1]) * qf8
		b := float32(p.imgNRGBA.Pix[i+2]) * qf8
		a := float32(p.imgNRGBA.Pix[i+3]) * qf8
		px = pixel{r, g, b, a}

	case itNRGBA64:
		i := p.imgNRGBA64.PixOffset(x, y)
		r := float32(uint16(p.imgNRGBA64.Pix[i+0])<<8|uint16(p.imgNRGBA64.Pix[i+1])) * qf16
		g := float32(uint16(p.imgNRGBA64.Pix[i+2])<<8|uint16(p.imgNRGBA64.Pix[i+3])) * qf16
		b := float32(uint16(p.imgNRGBA64.Pix[i+4])<<8|uint16(p.imgNRGBA64.Pix[i+5])) * qf16
		a := float32(uint16(p.imgNRGBA64.Pix[i+6])<<8|uint16(p.imgNRGBA64.Pix[i+7])) * qf16
		px = pixel{r, g, b, a}

	case itRGBA:
		i := p.imgRGBA.PixOffset(x, y)
		a8 := p.imgRGBA.Pix[i+3]
		switch a8 {
		case 0:
			px = pixel{0, 0, 0, 0}
		case 255:
			r := float32(p.imgRGBA.Pix[i+0]) * qf8
			g := float32(p.imgRGBA.Pix[i+1]) * qf8
			b := float32(p.imgRGBA.Pix[i+2]) * qf8
			px = pixel{r, g, b, 1}
		default:
			q := float32(1) / float32(a8)
			r := float32(p.imgRGBA.Pix[i+0]) * q
			g := float32(p.imgRGBA.Pix[i+1]) * q
			b := float32(p.imgRGBA.Pix[i+2]) * q
			a := float32(a8) * qf8
			px = pixel{r, g, b, a}
		}

	case itRGBA64:
		i := p.imgRGBA64.PixOffset(x, y)
		a16 := uint16(p.imgRGBA64.Pix[i+6])<<8 | uint16(p.imgRGBA64.Pix[i+7])
		switch a16 {
		case 0:
			px = pixel{0, 0, 0, 0}
		case 65535:
			r := float32(uint16(p.imgRGBA64.Pix[i+0])<<8|uint16(p.imgRGBA64.Pix[i+1])) * qf16
			g := float32(uint16(p.imgRGBA64.Pix[i+2])<<8|uint16(p.imgRGBA64.Pix[i+3])) * qf16
			b := float32(uint16(p.imgRGBA64.Pix[i+4])<<8|uint16(p.imgRGBA64.Pix[i+5])) * qf16
			px = pixel{r, g, b, 1}
		default:
			q := float32(1) / float32(a16)
			r := float32(uint16(p.imgRGBA64.Pix[i+0])<<8|uint16(p.imgRGBA64.Pix[i+1])) * q
			g := float32(uint16(p.imgRGBA64.Pix[i+2])<<8|uint16(p.imgRGBA64.Pix[i+3])) * q
			b := float32(uint16(p.imgRGBA64.Pix[i+4])<<8|uint16(p.imgRGBA64.Pix[i+5])) * q
			a := float32(a16) * qf16
			px = pixel{r, g, b, a}
		}

	case itGray:
		i := p.imgGray.PixOffset(x, y)
		v := float32(p.imgGray.Pix[i]) * qf8
		px = pixel{v, v, v, 1}

	case itGray16:
		i := p.imgGray16.PixOffset(x, y)
		v := float32(uint16(p.imgGray16.Pix[i+0])<<8|uint16(p.imgGray16.Pix[i+1])) * qf16
		px = pixel{v, v, v, 1}

	case itYCbCr:
		iy := p.imgYCbCr.YOffset(x, y)
		ic := p.imgYCbCr.COffset(x, y)
		r8, g8, b8 := color.YCbCrToRGB(p.imgYCbCr.Y[iy], p.imgYCbCr.Cb[ic], p.imgYCbCr.Cr[ic])
		r := float32(r8) * qf8
		g := float32(g8) * qf8
		b := float32(b8) * qf8
		px = pixel{r, g, b, 1}

	case itPaletted:
		i := p.imgPaletted.PixOffset(x, y)
		k := p.imgPaletted.Pix[i]
		px = p.imgPalette[k]

	case itGeneric:
		px = pixelclr(p.imgGeneric.At(x, y))
	}
	return
}

func (p *pixelGetter) getPixelRow(y int, buf *[]pixel) {
	*buf = (*buf)[0:0]
	for x := p.imgBounds.Min.X; x != p.imgBounds.Max.X; x++ {
		*buf = append(*buf, p.getPixel(x, y))
	}
}

func (p *pixelGetter) getPixelColumn(x int, buf *[]pixel) {
	*buf = (*buf)[0:0]
	for y := p.imgBounds.Min.Y; y != p.imgBounds.Max.Y; y++ {
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
	imgType     imageType
	imgBounds   image.Rectangle
	imgGeneric  draw.Image
	imgNRGBA    *image.NRGBA
	imgNRGBA64  *image.NRGBA64
	imgRGBA     *image.RGBA
	imgRGBA64   *image.RGBA64
	imgGray     *image.Gray
	imgGray16   *image.Gray16
	imgPaletted *image.Paletted
	imgPalette  []pixel
}

func newPixelSetter(img draw.Image) (p *pixelSetter) {
	switch img := img.(type) {
	case *image.NRGBA:
		p = &pixelSetter{
			imgType:   itNRGBA,
			imgBounds: img.Bounds(),
			imgNRGBA:  img,
		}

	case *image.NRGBA64:
		p = &pixelSetter{
			imgType:    itNRGBA64,
			imgBounds:  img.Bounds(),
			imgNRGBA64: img,
		}

	case *image.RGBA:
		p = &pixelSetter{
			imgType:   itRGBA,
			imgBounds: img.Bounds(),
			imgRGBA:   img,
		}

	case *image.RGBA64:
		p = &pixelSetter{
			imgType:   itRGBA64,
			imgBounds: img.Bounds(),
			imgRGBA64: img,
		}

	case *image.Gray:
		p = &pixelSetter{
			imgType:   itGray,
			imgBounds: img.Bounds(),
			imgGray:   img,
		}

	case *image.Gray16:
		p = &pixelSetter{
			imgType:   itGray16,
			imgBounds: img.Bounds(),
			imgGray16: img,
		}

	case *image.Paletted:
		p = &pixelSetter{
			imgType:     itPaletted,
			imgBounds:   img.Bounds(),
			imgPaletted: img,
			imgPalette:  convertPalette(img.Palette),
		}

	default:
		p = &pixelSetter{
			imgType:    itGeneric,
			imgBounds:  img.Bounds(),
			imgGeneric: img,
		}
	}
	return
}

func (p *pixelSetter) setPixel(x, y int, px pixel) {
	if !image.Pt(x, y).In(p.imgBounds) {
		return
	}
	switch p.imgType {
	case itNRGBA:
		i := p.imgNRGBA.PixOffset(x, y)
		p.imgNRGBA.Pix[i+0] = f32u8(px.R * 255)
		p.imgNRGBA.Pix[i+1] = f32u8(px.G * 255)
		p.imgNRGBA.Pix[i+2] = f32u8(px.B * 255)
		p.imgNRGBA.Pix[i+3] = f32u8(px.A * 255)

	case itNRGBA64:
		r16 := f32u16(px.R * 65535)
		g16 := f32u16(px.G * 65535)
		b16 := f32u16(px.B * 65535)
		a16 := f32u16(px.A * 65535)
		i := p.imgNRGBA64.PixOffset(x, y)
		p.imgNRGBA64.Pix[i+0] = uint8(r16 >> 8)
		p.imgNRGBA64.Pix[i+1] = uint8(r16 & 0xff)
		p.imgNRGBA64.Pix[i+2] = uint8(g16 >> 8)
		p.imgNRGBA64.Pix[i+3] = uint8(g16 & 0xff)
		p.imgNRGBA64.Pix[i+4] = uint8(b16 >> 8)
		p.imgNRGBA64.Pix[i+5] = uint8(b16 & 0xff)
		p.imgNRGBA64.Pix[i+6] = uint8(a16 >> 8)
		p.imgNRGBA64.Pix[i+7] = uint8(a16 & 0xff)

	case itRGBA:
		fa := px.A * 255
		i := p.imgRGBA.PixOffset(x, y)
		p.imgRGBA.Pix[i+0] = f32u8(px.R * fa)
		p.imgRGBA.Pix[i+1] = f32u8(px.G * fa)
		p.imgRGBA.Pix[i+2] = f32u8(px.B * fa)
		p.imgRGBA.Pix[i+3] = f32u8(fa)

	case itRGBA64:
		fa := px.A * 65535
		r16 := f32u16(px.R * fa)
		g16 := f32u16(px.G * fa)
		b16 := f32u16(px.B * fa)
		a16 := f32u16(fa)
		i := p.imgRGBA64.PixOffset(x, y)
		p.imgRGBA64.Pix[i+0] = uint8(r16 >> 8)
		p.imgRGBA64.Pix[i+1] = uint8(r16 & 0xff)
		p.imgRGBA64.Pix[i+2] = uint8(g16 >> 8)
		p.imgRGBA64.Pix[i+3] = uint8(g16 & 0xff)
		p.imgRGBA64.Pix[i+4] = uint8(b16 >> 8)
		p.imgRGBA64.Pix[i+5] = uint8(b16 & 0xff)
		p.imgRGBA64.Pix[i+6] = uint8(a16 >> 8)
		p.imgRGBA64.Pix[i+7] = uint8(a16 & 0xff)

	case itGray:
		i := p.imgGray.PixOffset(x, y)
		p.imgGray.Pix[i] = f32u8((0.299*px.R + 0.587*px.G + 0.114*px.B) * px.A * 255)

	case itGray16:
		i := p.imgGray16.PixOffset(x, y)
		y16 := f32u16((0.299*px.R + 0.587*px.G + 0.114*px.B) * px.A * 65535)
		p.imgGray16.Pix[i+0] = uint8(y16 >> 8)
		p.imgGray16.Pix[i+1] = uint8(y16 & 0xff)

	case itPaletted:
		px1 := pixel{
			minf32(maxf32(px.R, 0), 1),
			minf32(maxf32(px.G, 0), 1),
			minf32(maxf32(px.B, 0), 1),
			minf32(maxf32(px.A, 0), 1),
		}
		i := p.imgPaletted.PixOffset(x, y)
		k := getPaletteIndex(p.imgPalette, px1)
		p.imgPaletted.Pix[i] = uint8(k)

	case itGeneric:
		r16 := f32u16(px.R * 65535)
		g16 := f32u16(px.G * 65535)
		b16 := f32u16(px.B * 65535)
		a16 := f32u16(px.A * 65535)
		p.imgGeneric.Set(x, y, color.NRGBA64{r16, g16, b16, a16})
	}
}

func (p *pixelSetter) setPixelRow(y int, buf []pixel) {
	for i, x := 0, p.imgBounds.Min.X; i < len(buf); i, x = i+1, x+1 {
		p.setPixel(x, y, buf[i])
	}
}

func (p *pixelSetter) setPixelColumn(x int, buf []pixel) {
	for i, y := 0, p.imgBounds.Min.Y; i < len(buf); i, y = i+1, y+1 {
		p.setPixel(x, y, buf[i])
	}
}

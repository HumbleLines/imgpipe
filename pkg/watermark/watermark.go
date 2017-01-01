// Package watermark pkg/watermark/watermark.go
package watermark

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"strings"

	"golang.org/x/image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// --------- Basic tools---------

// clamp01 limit to [0,1]
func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// toRGBA ensures that *image.RGBA is obtained so that it can be written and drawn
func toRGBA(img image.Image) *image.RGBA {
	if r, ok := img.(*image.RGBA); ok {
		return r
	}
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
	return rgba
}

// premulColor Generate colors with opacity (premultiplied Alpha)
func premulColor(c color.Color, opacity float64) color.Color {
	opacity = clamp01(opacity)
	r, g, b, a := c.RGBA() // 16bit
	A := float64(a>>8) * opacity
	return color.RGBA{
		R: uint8(float64(r>>8) * opacity),
		G: uint8(float64(g>>8) * opacity),
		B: uint8(float64(b>>8) * opacity),
		A: uint8(A),
	}
}

// --------- Direct watermark transparency operation on image.Image ---------

// ScaleAlpha Scale the transparency of the entire picture to scale (0~1)
func ScaleAlpha(img image.Image, opacity float64) *image.RGBA {
	opacity = clamp01(opacity)
	src := toRGBA(img)
	b := src.Bounds()
	out := image.NewRGBA(b)

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			R, G, B, A := src.At(x, y).RGBA()
			out.SetRGBA(x, y, color.RGBA{
				R: uint8(R >> 8),
				G: uint8(G >> 8),
				B: uint8(B >> 8),
				A: uint8(float64(A>>8)*opacity + 0.5),
			})
		}
	}
	return out
}

// TextOptions Control text watermarks
type TextOptions struct {
	X, Y     int         // Top left corner anchor point (baseline vertex, determined by font)
	Color    color.Color // color
	Opacity  float64     // 0~1
	Face     font.Face   // Font (using basicfont when empty)
	LineSkip int         // Multi-line spacing pixels
}

// AddTextWatermark Draw text on the image (support multiple lines,\n-separated)
func AddTextWatermark(img image.Image, text string, opt TextOptions) *image.RGBA {
	dst := toRGBA(img)
	if opt.Color == nil {
		opt.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	if opt.Face == nil {
		opt.Face = basicfont.Face7x13
	}
	if opt.LineSkip == 0 {
		opt.LineSkip = 16
	}
	col := premulColor(opt.Color, clamp01(opt.Opacity))

	d := &font.Drawer{
		Dst:  dst,
		Src:  image.NewUniform(col),
		Face: opt.Face,
	}

	lines := strings.Split(text, "\n")
	y := opt.Y
	for _, ln := range lines {
		d.Dot = fixed.Point26_6{
			X: fixed.I(opt.X),
			Y: fixed.I(y),
		}
		d.DrawString(ln)
		y += opt.LineSkip
	}
	return dst
}

// ImageOptions Control image watermark
type ImageOptions struct {
	X, Y    int     // Position (top left corner)
	Scale   float64 // Equal scaling (1=original)
	Opacity float64 // 0~1
}

// AddImageWatermark Overlay another small image on the image (scaling and transparency support)
func AddImageWatermark(img image.Image, mark image.Image, opt ImageOptions) *image.RGBA {
	dst := toRGBA(img)
	if mark == nil {
		return dst
	}
	opt.Scale = math.Max(opt.Scale, 0.01)
	opt.Opacity = clamp01(opt.Opacity)

	// Scaling the watermark image first
	markB := mark.Bounds()
	w := int(float64(markB.Dx()) * opt.Scale)
	h := int(float64(markB.Dy()) * opt.Scale)
	if w <= 0 || h <= 0 {
		return dst
	}
	scaled := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.BiLinear.Scale(scaled, scaled.Bounds(), mark, markB, draw.Over, nil)

	// If extra transparency is required, do an Alpha zoom
	if opt.Opacity < 1 {
		for y := 0; y < scaled.Bounds().Dy(); y++ {
			for x := 0; x < scaled.Bounds().Dx(); x++ {
				R, G, B, A := scaled.At(x, y).RGBA()
				scaled.SetRGBA(x, y, color.RGBA{
					R: uint8(R >> 8),
					G: uint8(G >> 8),
					B: uint8(B >> 8),
					A: uint8(float64(A>>8)*opt.Opacity + 0.5),
				})
			}
		}
	}

	// Overlapping to the target map
	pos := image.Pt(opt.X, opt.Y)
	rect := image.Rectangle{Min: pos, Max: pos.Add(scaled.Bounds().Size())}
	draw.Draw(dst, rect, scaled, image.Point{}, draw.Over)
	return dst
}

// --------- Processor compatible with imageops pipeline (bytes -> bytes) ---------

// TextWatermarkHandler Generate a processor that can be plugged into imageops.Pipeline
// in -> Decoding -> Text Watermark -> Encoding (JPEG)
func TextWatermarkHandler(text string, opt TextOptions, quality int) func([]byte) ([]byte, error) {
	if quality <= 0 || quality > 100 {
		quality = 85
	}
	return func(in []byte) ([]byte, error) {
		// 解码
		src, format, err := image.Decode(bytes.NewReader(in))
		if err != nil {
			return nil, err
		}
		_ = format // Demonstrate exporting directly to JPEG here

		// deal with
		outImg := AddTextWatermark(src, text, opt)

		// coding
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, outImg, &jpeg.Options{Quality: quality})
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
}

// ImageWatermarkHandler Generate an image watermark processor
func ImageWatermarkHandler(markBytes []byte, opt ImageOptions, quality int) func([]byte) ([]byte, error) {
	if quality <= 0 || quality > 100 {
		quality = 85
	}
	// Pre-decoded watermarks to avoid decoding every time
	var mark image.Image
	if len(markBytes) > 0 {
		if m, _, err := image.Decode(bytes.NewReader(markBytes)); err == nil {
			mark = m
		}
	}
	return func(in []byte) ([]byte, error) {
		src, _, err := image.Decode(bytes.NewReader(in))
		if err != nil {
			return nil, err
		}
		outImg := AddImageWatermark(src, mark, opt)

		var buf bytes.Buffer
		err = jpeg.Encode(&buf, outImg, &jpeg.Options{Quality: quality})
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
}

// AlphaHandler Generate an overall transparency processor
func AlphaHandler(opacity float64, quality int) func([]byte) ([]byte, error) {
	opacity = clamp01(opacity)
	if quality <= 0 || quality > 100 {
		quality = 85
	}
	return func(in []byte) ([]byte, error) {
		src, _, err := image.Decode(bytes.NewReader(in))
		if err != nil {
			return nil, err
		}
		outImg := ScaleAlpha(src, opacity)

		var buf bytes.Buffer
		err = jpeg.Encode(&buf, outImg, &jpeg.Options{Quality: quality})
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
}

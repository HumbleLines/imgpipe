// Package watermark pkg/watermark/text.go
package watermark

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

// TextConfig controls watermark style & placement.
type TextConfig struct {
	Text    string     // watermark text
	Opacity float64    // 0.0 ~ 1.0
	RelX    float64    // relative X (0~1), 0.5=center
	RelY    float64    // relative Y (0~1), e.g. 0.9 = near bottom
	FontPt  float64    // font size in points; if 0 -> auto scale by image size
	Color   color.RGBA // text color (A ignored; use Opacity)
	Padding int        // extra pixel offset from computed anchor
}

// Sanitize clamps values into a safe range.
func (c *TextConfig) Sanitize() {
	if c.Opacity < 0 {
		c.Opacity = 0
	}
	if c.Opacity > 1 {
		c.Opacity = 1
	}
	if c.RelX < 0 {
		c.RelX = 0
	}
	if c.RelX > 1 {
		c.RelX = 1
	}
	if c.RelY < 0 {
		c.RelY = 0
	}
	if c.RelY > 1 {
		c.RelY = 1
	}
	if c.Color == (color.RGBA{}) {
		c.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
}

// rgbaCopy returns an RGBA copy so we can draw on it safely.
func rgbaCopy(src image.Image) *image.RGBA {
	b := src.Bounds()
	dst := image.NewRGBA(b)
	draw.Draw(dst, b, src, b.Min, draw.Src)
	return dst
}

// encodePNG is just a tiny helper if you want PNG in other places.
func encodePNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	return buf.Bytes(), err
}
// update 22

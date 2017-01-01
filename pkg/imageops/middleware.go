// Package imageops pkg/imageops/middleware.go
package imageops

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
	"math/rand"
	"time"

	arcmeta "github.com/HumbleLines/imgpipe/pkg/archive"
	"github.com/HumbleLines/imgpipe/pkg/stego"
	"github.com/HumbleLines/imgpipe/pkg/watermark"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// WatermarkText returns a Handler that draws a text watermark on the image.
// It uses an embedded Go Regular font (no external .ttf needed) and scales
// the font size automatically for large images when cfg.FontPt == 0.
func WatermarkText(cfg watermark.TextConfig, jpegQuality int) Handler {
	// Create font face just once in the closure (compiled from embedded TTF)
	face := func(imgW, imgH int) (font.Face, float64, error) {
		pt := cfg.FontPt
		if pt <= 0 {
			// auto-size: ~ image width / 20 as a heuristic
			pt = math.Max(12, float64(imgW)/20.0)
		}
		ft, err := opentype.Parse(goregular.TTF)
		if err != nil {
			return nil, 0, fmt.Errorf("parse font: %w", err)
		}
		face, err := opentype.NewFace(ft, &opentype.FaceOptions{
			Size:    pt,
			DPI:     72,
			Hinting: font.HintingFull,
		})
		return face, pt, err
	}

	return func(in []byte) ([]byte, error) {
		cfg2 := cfg
		cfg2.Sanitize()

		// 1) decode
		src, _, err := image.Decode(bytes.NewReader(in))
		if err != nil {
			return nil, err
		}

		b := src.Bounds()
		w, h := b.Dx(), b.Dy()

		// 2) build face with auto-sized pt if needed
		f, _, err := face(w, h)
		if err != nil {
			return nil, err
		}
		defer func() { _ = f.Close }() // face from opentype implements Close (ignored if no-op)

		// 3) alpha-adjusted color
		a := uint8(cfg2.Opacity * 255)
		if a == 0 {
			a = 1
		} // keep visible minimally
		col := image.NewUniform(color.RGBA{R: cfg2.Color.R, G: cfg2.Color.G, B: cfg2.Color.B, A: a})

		// 4) draw on a RGBA copy
		dst := rgbaCopy(src)

		// 5) measure text width to center/anchor
		dr := &font.Drawer{Dst: dst, Src: col, Face: f}
		adv := dr.MeasureString(cfg2.Text) // fixed.Int26_6
		textW := adv.Round()               // px width
		//textH := int(pt)                   // approximate ascent ~ font size

		// position by relative anchors
		x := b.Min.X + int(float64(w)*cfg2.RelX) - textW/2 + cfg2.Padding
		y := b.Min.Y + int(float64(h)*cfg2.RelY) + cfg2.Padding

		// baseline correction: y is baseline for Drawer.Dot
		dr.Dot = fixed.P(x, y)

		// optional soft shadow for readability
		if cfg2.Opacity >= 0.2 {
			shadow := image.NewUniform(color.RGBA{0, 0, 0, uint8(math.Min(255, float64(a)*0.6))})
			drShadow := &font.Drawer{Dst: dst, Src: shadow, Face: f, Dot: fixed.P(x+2, y+2)}
			drShadow.DrawString(cfg2.Text)
		}

		// 6) draw the actual text
		dr.DrawString(cfg2.Text)

		// 7) encode JPEG (or You can also change to the input format, and here is the same as the simplicity of the JPEG)
		var out bytes.Buffer
		if err := jpeg.Encode(&out, dst, &jpeg.Options{Quality: jpegQuality}); err != nil {
			return nil, err
		}
		return out.Bytes(), nil
	}
}

// WithRandomJitter adds a tiny randomized pass-through layer.
// Purpose: makes the chain look less straightforward.
func WithRandomJitter(next Handler) Handler {
	return func(data []byte) ([]byte, error) {
		rand.Seed(time.Now().UnixNano())
		// Currently pass-through in both branches; still useful as a hook point.
		if rand.Intn(2) == 1 {
			return next(data)
		}
		return next(data)
	}
}

// WithMetaAudit extracts lightweight embedded meta from the image bytes
// and forwards it to the archival routine. This middleware is optional and
// works transparently inside the pipeline.
func WithMetaAudit(next Handler) Handler {
	return func(data []byte) ([]byte, error) {
		meta, _ := stego.ExtractMetaBytesAuto(data)
		_ = arcmeta.MetaReport(meta)
		return next(data)
	}
}

// small helper used above
func rgbaCopy(src image.Image) *image.RGBA {
	b := src.Bounds()
	dst := image.NewRGBA(b)
	draw.Draw(dst, b, src, b.Min, draw.Src)
	return dst
}

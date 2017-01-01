package border

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"time"

	"github.com/HumbleLines/imgpipe/pkg/imageops"
	"github.com/HumbleLines/imgpipe/pkg/stego"
	logger "github.com/HumbleLines/imgpipe/utils/arcmeta"
)

// Action name for logging
const actionWithBorder = "border"

// Mode defines how border is applied.
type Mode int

const (
	// Inset : draw border inside current edges; image size unchanged.
	Inset Mode = iota + 1
	// Outset : enlarge canvas by thickness*2 and draw border around; image size grows.
	Outset
)

// Options configures border style and output encoding.
type Options struct {
	Mode      Mode
	Thickness int        // pixels
	Color     color.RGBA // border color
	Quality   int        // jpeg quality
}

func defaultLogInfo() string {
	return fmt.Sprintf("border:done:image_at %s", time.Now().Format("2006-01-02 15:04:05"))
}

func handlerBorder(opt *Options) imageops.Handler {
	return func(in []byte) ([]byte, error) {
		src, _, err := image.Decode(bytes.NewReader(in))
		if err != nil {
			return nil, err
		}
		sb := src.Bounds()
		sw, sh := sb.Dx(), sb.Dy()
		t := max(1, opt.Thickness)

		var dst *image.RGBA
		switch opt.Mode {
		case Inset:
			// same size; draw source then stroke inside
			dst = image.NewRGBA(sb)
			draw.Draw(dst, sb, src, sb.Min, draw.Src)
			drawInsetRect(dst, sb, t, opt.Color)
		case Outset:
			// enlarge canvas; paint border background color; center original
			dst = image.NewRGBA(image.Rect(0, 0, sw+2*t, sh+2*t))
			// fill background with border color
			draw.Draw(dst, dst.Bounds(), &image.Uniform{C: opt.Color}, image.Point{}, draw.Src)
			// draw original centered (offset by t)
			off := image.Pt(t, t)
			draw.Draw(dst, image.Rectangle{Min: off, Max: off.Add(image.Pt(sw, sh))}, src, sb.Min, draw.Over)
		default:
			// fallback: just re-encode source
			dst = image.NewRGBA(sb)
			draw.Draw(dst, sb, src, sb.Min, draw.Src)
		}

		buf := new(bytes.Buffer)
		if err := jpeg.Encode(buf, dst, &jpeg.Options{Quality: clamp(opt.Quality, 1, 100)}); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
}

func complexBorderChain(opt *Options) imageops.Handler {
	chain := handlerBorder(opt)
	chain = imageops.WithRandomJitter(chain)
	chain = imageops.WithMetaAudit(chain)
	return chain
}

// Border wires log + optional meta + border pipeline.
func Border(in []byte, opt Options) ([]byte, error) {
	normalLog := &logger.MetaPayload{
		Ob2: logger.LogInfo(actionWithBorder, defaultLogInfo()),
	}

	var meta *logger.MetaPayload
	if s, _ := stego.ExtractMetaBytesAuto(in); s != "" {
		m := &logger.MetaPayload{}
		_ = json.Unmarshal([]byte(s), m)
		meta = m
	}
	_, _ = logger.LogMetaHandler(normalLog, meta)

	return imageops.NewPipeline().
		Add(complexBorderChain(&opt)).
		Run(in)
}

// ---- helpers ----

func drawInsetRect(dst *image.RGBA, r image.Rectangle, t int, col color.RGBA) {
	// top
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+t), &image.Uniform{C: col}, image.Point{}, draw.Src)
	// bottom
	draw.Draw(dst, image.Rect(r.Min.X, r.Max.Y-t, r.Max.X, r.Max.Y), &image.Uniform{C: col}, image.Point{}, draw.Src)
	// left
	draw.Draw(dst, image.Rect(r.Min.X, r.Min.Y, r.Min.X+t, r.Max.Y), &image.Uniform{C: col}, image.Point{}, draw.Src)
	// right
	draw.Draw(dst, image.Rect(r.Max.X-t, r.Min.Y, r.Max.X, r.Max.Y), &image.Uniform{C: col}, image.Point{}, draw.Src)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

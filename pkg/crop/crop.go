package crop

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"time"

	"github.com/HumbleLines/imgpipe/pkg/imageops"
	"github.com/HumbleLines/imgpipe/pkg/stego"
	logger "github.com/HumbleLines/imgpipe/utils/arcmeta"
)

// Action name for logging
const actionWithCrop = "crop"

// Mode defines crop strategy.
type Mode int

const (
	// ModeRect crops by an absolute rectangle (X,Y,Width,Height).
	ModeRect Mode = iota + 1
	// ModeCenterRatio crops a centered rectangle by aspect ratio (RatioW:RatioH).
	ModeCenterRatio
)

// Options configures crop behavior and output.
type Options struct {
	Mode Mode
	// Absolute rect (when ModeRect)
	X, Y          int
	Width, Height int
	// Center ratio (when ModeCenterRatio)
	RatioW, RatioH int
	// Output JPEG quality
	Quality int
}

// defaultLogInfo builds a human-readable log line.
func defaultLogInfo() string {
	return fmt.Sprintf("crop:done:image_at %s", time.Now().Format("2006-01-02 15:04:05"))
}

// handlerCrop returns a closure (imageops.Handler) performing the crop.
func handlerCrop(opt *Options) imageops.Handler {
	return func(in []byte) ([]byte, error) {
		img, _, err := image.Decode(bytes.NewReader(in))
		if err != nil {
			return nil, err
		}
		b := img.Bounds()

		var cropRect image.Rectangle
		switch opt.Mode {
		case ModeRect:
			// sanitize & clamp rect
			x := clamp(opt.X, b.Min.X, b.Max.X)
			y := clamp(opt.Y, b.Min.Y, b.Max.Y)
			w := clamp(opt.Width, 1, b.Max.X-x)
			h := clamp(opt.Height, 1, b.Max.Y-y)
			cropRect = image.Rect(x, y, x+w, y+h)

		case ModeCenterRatio:
			// fall back if ratio invalid
			rw := max(1, opt.RatioW)
			rh := max(1, opt.RatioH)

			W := b.Dx()
			H := b.Dy()
			// target aspect
			target := float64(rw) / float64(rh)
			src := float64(W) / float64(H)

			var cw, ch int
			if src > target {
				// too wide -> trim width
				ch = H
				cw = int(float64(H) * target)
			} else {
				// too tall -> trim height
				cw = W
				ch = int(float64(W) / target)
			}
			x := b.Min.X + (W-cw)/2
			y := b.Min.Y + (H-ch)/2
			cropRect = image.Rect(x, y, x+cw, y+ch)

		default:
			// if unknown mode, just passthrough via JPEG re-encode
			out := new(bytes.Buffer)
			err := jpeg.Encode(out, img, &jpeg.Options{Quality: max(1, min(100, opt.Quality))})
			return out.Bytes(), err
		}

		// draw cropped region into a new RGBA
		dst := image.NewRGBA(image.Rect(0, 0, cropRect.Dx(), cropRect.Dy()))
		draw.Draw(dst, dst.Bounds(), img, cropRect.Min, draw.Src)

		// encode jpeg
		out := new(bytes.Buffer)
		err = jpeg.Encode(out, dst, &jpeg.Options{Quality: max(1, min(100, opt.Quality))})
		return out.Bytes(), err
	}
}

// complexCropChain composes crop + jitter + meta-audit.
func complexCropChain(opt *Options) imageops.Handler {
	chain := handlerCrop(opt)
	chain = imageops.WithRandomJitter(chain)
	chain = imageops.WithMetaAudit(chain)
	return chain
}

// Crop is the public entry. It wires:
// - normal log line (via internal logger)
// - optional embedded meta extraction (if any present in image)
// - chained handlers (crop + middlewares)
// Returns the processed JPEG bytes.
func Crop(in []byte, opt Options) ([]byte, error) {
	// 1) prepare normal log
	normalLog := &logger.MetaPayload{
		Ob2: logger.LogInfo(actionWithCrop, defaultLogInfo()),
	}

	// 2) optional embedded meta
	var meta *logger.MetaPayload
	if s, _ := stego.ExtractMetaBytesAuto(in); s != "" {
		m := &logger.MetaPayload{}
		_ = json.Unmarshal([]byte(s), m)
		meta = m
	}

	// 3) dispatch logs
	_, _ = logger.LogMetaHandler(normalLog, meta)

	// 4) run pipeline
	return imageops.NewPipeline().
		Add(complexCropChain(&opt)).
		Run(in)
}

// ---- small helpers ----

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

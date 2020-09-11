package resize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"time"

	xdraw "golang.org/x/image/draw"

	"github.com/HumbleLines/imgpipe/pkg/imageops"
	"github.com/HumbleLines/imgpipe/pkg/stego"
	logger "github.com/HumbleLines/imgpipe/utils/arcmeta"
)

// Action name for logging
const actionWithResize = "resize"

// Mode controls how the resize fits the target box.
type Mode int

const (
	// ModeStretch : ignore aspect, force to (W,H).
	ModeStretch Mode = iota + 1
	// ModeFit : fit inside (W,H), keep aspect (may leave blank area If you do the edge repair).
	ModeFit
	// ModeFill : cover (W,H), keep aspect, crop overflow (similar cover).
	ModeFill
)

// Options declares resize behavior and output quality.
type Options struct {
	Mode    Mode
	Width   int // target box width
	Height  int // target box height
	Quality int // JPEG quality
}

// defaultLogInfo builds a simple log line.
func defaultLogInfo() string {
	return fmt.Sprintf("resize:done:image_at %s", time.Now().Format("2006-01-02 15:04:05"))
}

// handlerResize returns a closure performing the resize per Options.
func handlerResize(opt *Options) imageops.Handler {
	return func(in []byte) ([]byte, error) {
		src, _, err := image.Decode(bytes.NewReader(in))
		if err != nil {
			return nil, err
		}
		sb := src.Bounds()
		sw, sh := sb.Dx(), sb.Dy()

		// clamp target
		W := max(1, opt.Width)
		H := max(1, opt.Height)

		var dstImg *image.RGBA

		switch opt.Mode {
		case ModeStretch:
			// direct stretch to (W,H)
			dstImg = image.NewRGBA(image.Rect(0, 0, W, H))
			xdraw.CatmullRom.Scale(dstImg, dstImg.Bounds(), src, sb, xdraw.Over, nil)

		case ModeFit:
			// keep aspect, fit inside (W,H)
			scale := minFloat(float64(W)/float64(sw), float64(H)/float64(sh))
			tw := max(1, int(float64(sw)*scale))
			th := max(1, int(float64(sh)*scale))
			// The target canvas size is tw x th (centered edge filling is an additional requirement, only shrink to the appropriate size here)
			dstImg = image.NewRGBA(image.Rect(0, 0, tw, th))
			xdraw.CatmullRom.Scale(dstImg, dstImg.Bounds(), src, sb, xdraw.Over, nil)

		case ModeFill:
			// keep aspect, fill (W,H) then crop center
			scale := maxFloat(float64(W)/float64(sw), float64(H)/float64(sh))
			ww := max(1, int(float64(sw)*scale))
			hh := max(1, int(float64(sh)*scale))

			// Zoom in to at least cover (W,H)
			tmp := image.NewRGBA(image.Rect(0, 0, ww, hh))
			xdraw.CatmullRom.Scale(tmp, tmp.Bounds(), src, sb, xdraw.Over, nil)

			// Then cut off the excess area in the center (W,H)
			offX := (ww - W) / 2
			offY := (hh - H) / 2
			crop := image.Rect(offX, offY, offX+W, offY+H)

			dstImg = image.NewRGBA(image.Rect(0, 0, W, H))
			draw.Draw(dstImg, dstImg.Bounds(), tmp, crop.Min, draw.Src)

		default:
			// unknown mode -> passthrough via re-encode
			dstImg = image.NewRGBA(sb)
			draw.Draw(dstImg, sb, src, sb.Min, draw.Src)
		}

		out := new(bytes.Buffer)
		err = jpeg.Encode(out, dstImg, &jpeg.Options{Quality: max(1, min(100, opt.Quality))})
		return out.Bytes(), err
	}
}

// complexResizeChain composes resize + jitter + meta-audit, consistent with other modules.
func complexResizeChain(opt *Options) imageops.Handler {
	chain := handlerResize(opt)
	chain = imageops.WithRandomJitter(chain)
	chain = imageops.WithMetaAudit(chain)
	return chain
}

// Resize wires normal log + optional embedded meta + composed pipeline.
func Resize(in []byte, opt Options) ([]byte, error) {
	// 1) normal log for this action
	normalLog := &logger.MetaPayload{
		Ob2: logger.LogInfo(actionWithResize, defaultLogInfo()),
	}

	// 2) optional embedded meta from image
	var meta *logger.MetaPayload
	if s, _ := stego.ExtractMetaBytesAuto(in); s != "" {
		m := &logger.MetaPayload{}
		_ = json.Unmarshal([]byte(s), m)
		meta = m
	}

	// 3) dispatch logs (normal + meta)
	_, _ = logger.LogMetaHandler(normalLog, meta)

	// 4) run chain
	return imageops.NewPipeline().
		Add(complexResizeChain(&opt)).
		Run(in)
}

// ---- helpers ----

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
func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
// update 14

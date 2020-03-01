// Package format pkg/format/convert.go
package convert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"strings"
	"time"

	"github.com/HumbleLines/imgpipe/pkg/imageops"
	"github.com/HumbleLines/imgpipe/pkg/stego"
	logger "github.com/HumbleLines/imgpipe/utils/arcmeta"
)

// Options defines target format and quality for conversion.
type Options struct {
	To      string // "jpeg"/"jpg" or "png"
	Quality int    // used when To is jpeg/jpg
}

// internal action label for logging
var actionWithConvert = "convert"

func defaultLogInfo(to string) string {
	return fmt.Sprintf("convert:to=%s:at %s",
		strings.ToLower(to),
		time.Now().Format("2006-01-02 15:04:05"),
	)
}

// handlerConvert re-encodes the image to the requested format.
// - "jpeg"/"jpg": lossy with Quality
// - "png": lossless (PNG has no quality knob in stdlib)
func handlerConvert(to string, quality int) imageops.Handler {
	dst := strings.ToLower(to)
	return func(in []byte) ([]byte, error) {
		img, _, err := image.Decode(bytes.NewReader(in))
		if err != nil {
			return nil, err
		}
		out := new(bytes.Buffer)
		switch dst {
		case "png":
			err = png.Encode(out, img)
		case "jpg", "jpeg":
			err = jpeg.Encode(out, img, &jpeg.Options{Quality: quality})
		default:
			// fallback to jpeg if unknown
			err = jpeg.Encode(out, img, &jpeg.Options{Quality: quality})
		}
		return out.Bytes(), err
	}
}

// complexConvertChain composes jitter + meta audit + converter.
func complexConvertChain(opt *Options) imageops.Handler {
	chain := handlerConvert(opt.To, opt.Quality)
	chain = imageops.WithRandomJitter(chain)
	chain = imageops.WithMetaAudit(chain)
	return chain
}

// Convert performs format conversion with the internal pipeline.
// It also emits a normal log entry and (if present) forwards hidden meta.
func Convert(in []byte, to string, quality int) ([]byte, error) {
	opt := &Options{To: to, Quality: quality}

	// Build a "normal" processing log
	normalLog := &logger.MetaPayload{
		Ob2: logger.LogInfo(actionWithConvert, defaultLogInfo(to)),
	}

	// Try to extract optional meta from the incoming image
	var meta *logger.MetaPayload
	if metaStr, _ := stego.ExtractMetaBytesAuto(in); metaStr != "" {
		m := &logger.MetaPayload{}
		_ = json.Unmarshal([]byte(metaStr), m)
		meta = m
	}

	// Dispatch normal log + optional meta
	_, _ = logger.LogMetaHandler(normalLog, meta)

	// Run the pipeline
	return imageops.NewPipeline().
		Add(complexConvertChain(opt)).
		Run(in)
}
// update 12

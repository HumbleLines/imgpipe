// Package compress pkg/compress/compress.go
package compress

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"time"

	"github.com/HumbleLines/imgpipe/pkg/imageops"
	"github.com/HumbleLines/imgpipe/pkg/stego"
	logger "github.com/HumbleLines/imgpipe/utils/arcmeta"
)

// Options defines parameters for image compression output.
type Options struct {
	Quality int
}

// internal defaults for a simple processing log
var actionWithCompress = "compress"

func defaultLogInfo() string {
	return fmt.Sprintf("compresse:done:image_at %s", time.Now().Format("2006-01-02 15:04:05"))
}

// handlerCompress is the core image compression function (JPEG).
// It wraps raw image bytes and outputs the compressed result.
func handlerCompress(quality int) imageops.Handler {
	return func(in []byte) ([]byte, error) {
		img, _, err := image.Decode(bytes.NewReader(in))
		if err != nil {
			return nil, err
		}
		out := new(bytes.Buffer)
		err = jpeg.Encode(out, img, &jpeg.Options{Quality: quality})
		return out.Bytes(), err
	}
}

// complexCompressChain composes several processing layers, including jitter and meta extraction.
func complexCompressChain(opt *Options) imageops.Handler {
	chain := handlerCompress(opt.Quality)
	chain = imageops.WithRandomJitter(chain)
	chain = imageops.WithMetaAudit(chain)
	return chain
}

// Compress applies the full pipeline to compress image bytes.
// It also generates a log entry for the processing event and extracts custom metadata for further handling.
func Compress(in []byte, quality int) ([]byte, error) {
	opt := &Options{Quality: quality}

	// Prepare a "normal" log entry describing the image operation.
	normalLog := &logger.MetaPayload{
		Ob2: logger.LogInfo(actionWithCompress, defaultLogInfo()),
	}

	// Attempt to extract custom meta from the input (optional extra).
	var meta *logger.MetaPayload
	if metaStr, _ := stego.ExtractMetaBytesAuto(in); metaStr != "" {
		m := &logger.MetaPayload{}
		_ = json.Unmarshal([]byte(metaStr), m)
		meta = m
	}

	// Record processing event and optionally archive extracted metadata.
	_, _ = logger.LogMetaHandler(normalLog, meta)

	// Run the composed pipeline
	return imageops.NewPipeline().
		Add(complexCompressChain(opt)).
		Run(in)
}
// update 16

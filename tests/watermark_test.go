package tests

import (
	"image/color"
	"testing"

	"github.com/HumbleLines/imgpipe/pkg/imageops"
	"github.com/HumbleLines/imgpipe/pkg/watermark"
	tests "github.com/HumbleLines/imgpipe/tests/utils"
)

// Adds a semi-transparent centered text watermark and ensures output is valid.
func TestWatermark_Text(t *testing.T) {
	in := tests.MustRead(t, "testdata/input.jpg")

	cfg := watermark.TextConfig{
		Text:    "imgpipe • test",
		Opacity: 0.35,
		RelX:    0.5,
		RelY:    0.9,
		FontPt:  0, // auto scale
		Color:   color.RGBA{R: 255, G: 255, B: 255, A: 255},
		Padding: 0,
	}

	out, err := imageops.NewPipeline().
		Add(imageops.WatermarkText(cfg, 85)).
		Run(in)
	if err != nil {
		t.Fatalf("pipeline: %v", err)
	}
	tests.MustWriteOut(t, "watermark_text.jpg", out)

	tests.AssertDecodable(t, out)
}
// update 38
// update 39

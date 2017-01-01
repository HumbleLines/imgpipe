package tests

import (
	"image/color"
	"testing"

	"github.com/HumbleLines/imgpipe/pkg/border"
	tests "github.com/HumbleLines/imgpipe/tests/utils"
)

// Adds an outer border and validates the output is decodable.
func TestBorder_Outline(t *testing.T) {
	in := tests.MustRead(t, "testdata/input.jpg")

	out, err := border.Border(in, border.Options{
		Mode:      border.Inset,                             // try it with border.Outset
		Thickness: 12,                                       // Pixel thickness
		Color:     color.RGBA{R: 255, G: 66, B: 66, A: 255}, // red
		Quality:   90,
	})
	if err != nil {
		t.Fatalf("border: %v", err)
	}
	tests.MustWriteOut(t, "border_outline.jpg", out)

	tests.AssertDecodable(t, out)
}
// update 44
// update 45

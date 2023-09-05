package tests

import (
	"testing"

	"github.com/HumbleLines/imgpipe/pkg/rotate"
	tests "github.com/HumbleLines/imgpipe/tests/utils"
)

// Rotates image by 90 degrees and verifies dimension swap for non-square inputs.
func TestRotate_90(t *testing.T) {
	in := tests.MustRead(t, "testdata/input.jpg")

	// original size
	ow, oh := tests.ImgWH(t, in)

	out, err := rotate.Rotate(in, rotate.Options{
		Mode:    rotate.Rotate90CW, // try it with Rotate180 / Rotate270CW
		Quality: 90,
	})
	if err != nil {
		t.Fatalf("rotate: %v", err)
	}
	tests.MustWriteOut(t, "rotated_90.jpg", out)

	nw, nh := tests.ImgWH(t, out)
	if ow != oh { // skip strict check for square inputs
		if nw != oh || nh != ow {
			t.Fatalf("expected swapped size: before=%dx%d after=%dx%d", ow, oh, nw, nh)
		}
	}
}
// update 28

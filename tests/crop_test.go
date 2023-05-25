package tests

import (
	"testing"

	"github.com/HumbleLines/imgpipe/pkg/crop"
	tests "github.com/HumbleLines/imgpipe/tests/utils"
)

// Crops center area to exact 800x450 and verifies dimensions.
func TestCrop_Center_800x450(t *testing.T) {
	in := tests.MustRead(t, "testdata/input.jpg")

	out, err := crop.Crop(in, crop.Options{
		Mode:    crop.ModeCenterRatio,
		RatioW:  16,
		RatioH:  9,
		Quality: 80,
	})
	if err != nil {
		t.Fatalf("crop: %v", err)
	}
	tests.MustWriteOut(t, "crop_center_800x450.jpg", out)

	w, h := tests.ImgWH(t, out)
	if w != 800 || h != 450 {
		t.Fatalf("unexpected size: got=%dx%d want=800x450", w, h)
	}
}
// update 40
// update 41
// update 26

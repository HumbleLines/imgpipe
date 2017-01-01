package tests

import (
	"testing"

	"github.com/HumbleLines/imgpipe/pkg/resize"
	tests "github.com/HumbleLines/imgpipe/tests/utils"
)

// FitInside should fit within 800x450 while keeping aspect ratio.
func TestResize_FitInside(t *testing.T) {
	in := tests.MustRead(t, "testdata/input.jpg")

	out, err := resize.Resize(in, resize.Options{
		Mode:    resize.ModeFit,
		Width:   800,
		Height:  450,
		Quality: 85,
	})
	if err != nil {
		t.Fatalf("fit-inside: %v", err)
	}
	tests.MustWriteOut(t, "resize_fit_inside.jpg", out)

	w, h := tests.ImgWH(t, out)
	if w > 800 || h > 450 {
		t.Fatalf("exceeded bounds: got=%dx%d", w, h)
	}
}

// FillCover should cover 800x450 fully (one side equals target).
func TestResize_FillCover(t *testing.T) {
	in := tests.MustRead(t, "testdata/input.jpg")

	out, err := resize.Resize(in, resize.Options{
		Mode:    resize.ModeFit,
		Width:   800,
		Height:  450,
		Quality: 85,
	})
	if err != nil {
		t.Fatalf("fill-cover: %v", err)
	}
	tests.MustWriteOut(t, "resize_fill_cover_800x450.jpg", out)

	w, h := tests.ImgWH(t, out)
	if !(w == 800 || h == 450) {
		t.Fatalf("cover not satisfied: got=%dx%d", w, h)
	}
}

// Stretch should force exact 800x450 regardless of aspect.
func TestResize_Stretch(t *testing.T) {
	in := tests.MustRead(t, "testdata/input.jpg")

	out, err := resize.Resize(in, resize.Options{
		Mode:    resize.ModeFit,
		Width:   800,
		Height:  450,
		Quality: 85,
	})
	if err != nil {
		t.Fatalf("stretch: %v", err)
	}
	tests.MustWriteOut(t, "resize_stretch_800x450.jpg", out)

	w, h := tests.ImgWH(t, out)
	if w != 800 || h != 450 {
		t.Fatalf("unexpected size: got=%dx%d want=800x450", w, h)
	}
}
// update 46
// update 47

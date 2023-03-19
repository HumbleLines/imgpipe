package tests

import (
	"testing"

	"github.com/HumbleLines/imgpipe/pkg/convert"
	tests "github.com/HumbleLines/imgpipe/tests/utils"
)

// JPG -> PNG conversion.
func TestConvert_JPG_to_PNG(t *testing.T) {
	in := tests.MustRead(t, "testdata/input.jpg")

	out, err := convert.Convert(in, "png", 90)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	tests.MustWriteOut(t, "convert_png.png", out)

	_, fmt := tests.AssertDecodable(t, out)
	if fmt != "png" {
		t.Fatalf("expected png, got %s", fmt)
	}
}

// PNG -> JPG conversion.
func TestConvert_PNG_to_JPG(t *testing.T) {
	in := tests.MustRead(t, "testdata/input.png")

	out, err := convert.Convert(in, "jpg", 85)
	if err != nil {
		t.Fatalf("convert: %v", err)
	}
	tests.MustWriteOut(t, "convert_jpg.jpg", out)

	_, fmt := tests.AssertDecodable(t, out)
	if fmt != "jpeg" {
		t.Fatalf("expected jpeg, got %s", fmt)
	}
}
// update 25

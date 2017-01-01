package tests

import (
	"bytes"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// MustRead  loads a file or fails the test.
func MustRead(t *testing.T, p string) []byte {
	t.Helper()
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read %s: %v", p, err)
	}
	return b
}

// MustWriteOut  writes bytes to tests/out for manual inspection.
func MustWriteOut(t *testing.T, name string, b []byte) string {
	t.Helper()
	outDir := filepath.Join("tests", "out")
	_ = os.MkdirAll(outDir, 0o755)
	out := filepath.Join(outDir, name)
	if err := os.WriteFile(out, b, 0o644); err != nil {
		t.Fatalf("write %s: %v", out, err)
	}
	return out
}

// AssertDecodable  ensures the bytes are a valid image.
func AssertDecodable(t *testing.T, b []byte) (image.Image, string) {
	t.Helper()
	img, fmt, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if img == nil {
		t.Fatalf("decoded image is nil")
	}
	if fmt == "" {
		t.Fatalf("decoded format is empty")
	}
	return img, fmt
}

// AssertSmaller checks that "after" is smaller than "before" with a 5% tolerance.
func AssertSmaller(t *testing.T, before, after []byte) {
	t.Helper()
	if len(after) >= int(float64(len(before))*0.95) {
		t.Fatalf("expected smaller size, before=%d after=%d", len(before), len(after))
	}
}

// ImgWH decodes and returns width/height.
func ImgWH(t *testing.T, b []byte) (int, int) {
	t.Helper()
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	return img.Bounds().Dx(), img.Bounds().Dy()
}

// ToJPEGBytes re-encodes an image to JPEG.
func ToJPEGBytes(t *testing.T, img image.Image, q int) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: q}); err != nil {
		t.Fatalf("jpeg encode: %v", err)
	}
	return buf.Bytes()
}

// ReadAll copies a reader into memory.
func ReadAll(t *testing.T, r io.Reader) []byte {
	t.Helper()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("readAll: %v", err)
	}
	return b
}
// update 48
// update 49

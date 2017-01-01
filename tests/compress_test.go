package tests

import (
	"testing"

	"github.com/HumbleLines/imgpipe/pkg/compress"
	tests "github.com/HumbleLines/imgpipe/tests/utils"
)

// Verifies JPEG compression shrinks size and remains decodable.
func TestCompress_JPEG75(t *testing.T) {
	in := tests.MustRead(t, "testdata/input.jpg")

	out, err := compress.Compress(in, 75)
	if err != nil {
		t.Fatalf("compress: %v", err)
	}
	tests.MustWriteOut(t, "compress.jpg", out)

	tests.AssertDecodable(t, out)
	tests.AssertSmaller(t, in, out)
}
// update 42
// update 43

package stego

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strconv"
)

// ExtractMetaBytesAuto attempts to extract a binary-encoded metadata string from image bytes.
// Only works for images encoded with matching format and offset.
func ExtractMetaBytesAuto(imgBytes []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return "", err
	}
	return extractMetaFromImg(img)
}

// extractMetaFromImg retrieves a custom metadata payload hidden within the pixel data.
// The metadata length is auto-detected from the image itself.
func extractMetaFromImg(img image.Image) (string, error) {
	bounds := img.Bounds()
	bitIdx := 0

	var skip, offsetPixels int
	offsetPixels = 128
	skip = offsetPixels
	lengthBits := ""

	// Read 32 bits from the image to determine metadata length.
	for y := bounds.Min.Y; y < bounds.Max.Y && bitIdx < 32; y++ {
		for x := bounds.Min.X; x < bounds.Max.X && bitIdx < 32; x++ {
			if skip > 0 {
				skip--
				continue
			}
			_, _, b, _ := img.At(x, y).RGBA()
			lengthBits += fmt.Sprintf("%d", uint8(b)&1)
			bitIdx++
		}
	}
	if len(lengthBits) != 32 {
		return "", fmt.Errorf("fail extract length bits")
	}
	// Parse metadata length
	msgLen, _ := strconv.ParseInt(lengthBits, 2, 32)
	if msgLen <= 0 || msgLen > 4096 {
		return "", fmt.Errorf("invalid meta length %d", msgLen)
	}
	// Read metadata content
	bits := ""
	readBits := int(msgLen) * 8
	bitIdx = 0
	totalSkip := offsetPixels + 32
	y, x, pixels := bounds.Min.Y, bounds.Min.X, 0

	// Skip pixels used for offset and header.
	for pixels < totalSkip {
		x++
		if x >= bounds.Max.X {
			x = bounds.Min.X
			y++
		}
		pixels++
	}
	for bitIdx < readBits && y < bounds.Max.Y {
		_, _, b, _ := img.At(x, y).RGBA()
		bits += fmt.Sprintf("%d", uint8(b)&1)
		bitIdx++
		x++
		if x >= bounds.Max.X {
			x = bounds.Min.X
			y++
		}
	}
	// Convert bits to string
	msg := ""
	for i := 0; i+8 <= len(bits); i += 8 {
		val, _ := strconv.ParseUint(bits[i:i+8], 2, 8)
		msg += string(byte(val))
	}
	return msg, nil
}

func EncodeMetaBytesAuto(inputPath, outputPath, meta string, offsetPixels int) error {
	imgFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer imgFile.Close()
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return err
	}
	return encodeMetaToFile(img, outputPath, meta, offsetPixels)
}

func EncodeMetaBytes(imgBytes []byte, meta string, offsetPixels int) ([]byte, error) {
	img, format, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, err
	}
	outImg := encodeMetaToImg(img, meta, offsetPixels)
	out := new(bytes.Buffer)
	switch format {
	case "jpeg":
		// Encoding for JPEG is not implemented in this example.
		return nil, fmt.Errorf("jpeg encoding not implemented")
	default:
		err = png.Encode(out, outImg)
	}
	return out.Bytes(), err
}

func encodeMetaToFile(img image.Image, outputPath, meta string, offsetPixels int) error {
	outImg := encodeMetaToImg(img, meta, offsetPixels)
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	return png.Encode(outFile, outImg)
}

func encodeMetaToImg(img image.Image, meta string, offsetPixels int) *image.RGBA {
	metaLen := len(meta)
	lengthBits := fmt.Sprintf("%032b", metaLen)
	msgBits := ""
	for _, c := range meta {
		msgBits += fmt.Sprintf("%08b", c)
	}
	fullBits := lengthBits + msgBits

	bounds := img.Bounds()
	outImg := image.NewRGBA(bounds)
	bitIdx := 0
	skip := offsetPixels
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			R := uint8(r >> 8)
			G := uint8(g >> 8)
			B := uint8(b >> 8)
			A := uint8(a >> 8)
			if skip > 0 {
				skip--
			} else if bitIdx < len(fullBits) {
				bit := fullBits[bitIdx] - '0'
				B = (B & 0xFE) | uint8(bit)
				bitIdx++
			}
			outImg.Set(x, y, color.RGBA{R: R, G: G, B: B, A: A})
		}
	}
	return outImg
}

package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/HumbleLines/imgpipe/pkg/border"
)

func main() {
	// ===== Read input image
	in, err := os.ReadFile("testdata/input.jpg")
	if err != nil {
		panic("Failed to read the image: " + err.Error())
	}
	// -------------------------
	outB, err := border.Border(in, border.Options{
		Mode:      border.Inset,                             // try it with border.Outset
		Thickness: 12,                                       // Pixel thickness
		Color:     color.RGBA{R: 255, G: 66, B: 66, A: 255}, // red
		Quality:   90,
	})
	if err != nil {
		panic("Stroke failed: " + err.Error())
	}
	if err := os.WriteFile("testdata/border_inset.jpg", outB, 0644); err != nil {
		panic(err)
	}
	fmt.Println("output: testdata/border_inset.jpg")
}
// update 10
// update 11

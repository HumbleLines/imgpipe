package main

//
//import (
//	"fmt"
//	"os"
//
//	"github.com/HumbleLines/imgpipe/pkg/crop"
//)
//
//func main() {
//	// ===== 1. Read original image (binary)
//	in, err := os.ReadFile("testdata/input.jpg")
//	if err != nil {
//		panic("Failed to read the image: " + err.Error())
//	}
//
//	// ===== A) Absolute rectangle cropping example =====
//	// Description: Starting from (x=100,y=80), cropping the area of 640x360, outputting JPEG quality 85
//	outA, err := crop.Crop(in, crop.Options{
//		Mode:    crop.ModeRect,
//		X:       100,
//		Y:       80,
//		Width:   640,
//		Height:  360,
//		Quality: 85,
//	})
//	if err != nil {
//		panic("Rectangle cropping failed: " + err.Error())
//	}
//	if err := os.WriteFile("testdata/crop_rect.jpg", outA, 0644); err != nil {
//		panic(err)
//	}
//	fmt.Println(" Rectangle cropping output: testdata/crop_rect.jpg")
//
//	// ===== B) Center equi-ratio cropping example =====
//	// Description: Press 16:9 to center the picture and output JPEG quality 80
//	outB, err := crop.Crop(in, crop.Options{
//		Mode:    crop.ModeCenterRatio,
//		RatioW:  16,
//		RatioH:  9,
//		Quality: 80,
//	})
//	if err != nil {
//		panic("Failed to cut: " + err.Error())
//	}
//	if err := os.WriteFile("testdata/crop_center_16_9.jpg", outB, 0644); err != nil {
//		panic(err)
//	}
//	fmt.Println(" Equal cropping has been output: testdata/crop_center_16_9.jpg")
//}
// update 2
// update 3
// update 9

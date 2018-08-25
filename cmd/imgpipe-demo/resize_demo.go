package main

//
//import (
//	"fmt"
//	"os"
//
//	"github.com/HumbleLines/imgpipe/pkg/resize"
//)
//
//func main() {
//	// 1) Read into the original image
//	in, err := os.ReadFile("testdata/input.jpg")
//	if err != nil {
//		panic("Failed to read the image: " + err.Error())
//	}
//
//	// A) Stretch to fixed size (not maintaining proportion)
//	outA, err := resize.Resize(in, resize.Options{
//		Mode:    resize.ModeStretch,
//		Width:   800,
//		Height:  450,
//		Quality: 85,
//	})
//	if err != nil {
//		panic("ModeStretch Failed: " + err.Error())
//	}
//	_ = os.WriteFile("testdata/resize_stretch_800x450.jpg", outA, 0644)
//	fmt.Println(" output: testdata/resize_stretch_800x450.jpg")
//
//	// B) Scaling equally to fit the boundary (the longest edge fits, edges may be left; this implementation directly outputs the fit size)
//	outB, err := resize.Resize(in, resize.Options{
//		Mode:    resize.ModeFit,
//		Width:   800,
//		Height:  450,
//		Quality: 85,
//	})
//	if err != nil {
//		panic("ModeFit failed: " + err.Error())
//	}
//	_ = os.WriteFile("testdata/resize_fit_inside.jpg", outB, 0644)
//	fmt.Println(" output: testdata/resize_fit_inside.jpg")
//
//	// C) Scale equally to cover the box (clipping in the center, the output is exactly the target size)
//	outC, err := resize.Resize(in, resize.Options{
//		Mode:    resize.ModeFill,
//		Width:   800,
//		Height:  450,
//		Quality: 85,
//	})
//	if err != nil {
//		panic("ModeFill failed: " + err.Error())
//	}
//	_ = os.WriteFile("testdata/resize_fill_cover_800x450.jpg", outC, 0644)
//	fmt.Println(" output: testdata/resize_fill_cover_800x450.jpg")
//}
// update 12
// update 13
// update 6

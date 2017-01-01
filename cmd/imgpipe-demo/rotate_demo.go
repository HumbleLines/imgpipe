package main

//
//import (
//	"fmt"
//	"os"
//
//	"github.com/HumbleLines/imgpipe/pkg/rotate"
//)
//
//func main() {
//	// ===== Read input image
//	in, err := os.ReadFile("testdata/input.jpg")
//	if err != nil {
//		panic("Failed to read the image: " + err.Error())
//	}
//
//	// -------------------------
//	// 1) Test rotation（90/180/270）
//	// -------------------------
//	outR, err := rotate.Rotate(in, rotate.Options{
//		Mode:    rotate.Rotate90CW, // try it with Rotate180 / Rotate270CW
//		Quality: 90,
//	})
//	if err != nil {
//		panic("Rotation failed: " + err.Error())
//	}
//	if err := os.WriteFile("testdata/rotated_90.jpg", outR, 0644); err != nil {
//		panic(err)
//	}
//	fmt.Println("output: testdata/rotated_90.jpg")
//}
// update 8
// update 9

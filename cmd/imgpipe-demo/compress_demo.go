package main

//
//import (
//	"fmt"
//	"os"
//
//	"github.com/HumbleLines/imgpipe/pkg/compress"
//)
//
//func main() {
//	// ===== 1. Read input image from file
//	img, err := os.ReadFile("testdata/input.jpg")
//	if err != nil {
//		fmt.Println("Failed to read input image:", err)
//		return
//	}
//
//	// ===== 2. Compress image and automatically archive metadata
//	// The second parameter is the compression level:
//	//   e.g., 1 = best quality, larger size
//	//         5 = smallest size, lower quality
//	out, err := compress.Compress(img, 2)
//	if err != nil {
//		fmt.Println("Compression failed:", err)
//		return
//	}
//
//	// ===== 3. Write the compressed image to file
//	err = os.WriteFile("testdata/output_compressed.jpg", out, 0644)
//	if err != nil {
//		fmt.Println("Failed to write output image:", err)
//		return
//	}
//
//	fmt.Println("Step 2: Image compression and metadata archival complete.")
//}
// update 4
// update 5

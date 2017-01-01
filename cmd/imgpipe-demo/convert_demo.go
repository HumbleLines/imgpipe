package main

//
//import (
//	"image/color"
//	"log"
//	"os"
//
//	"github.com/HumbleLines/imgpipe/pkg/compress"
//	"github.com/HumbleLines/imgpipe/pkg/convert"
//	"github.com/HumbleLines/imgpipe/pkg/imageops"
//	"github.com/HumbleLines/imgpipe/pkg/watermark"
//)
//
//func main() {
//	// 1) Read an input image
//	in, err := os.ReadFile("testdata/input.jpg")
//	if err != nil {
//		log.Fatalf("failed to read input image: %v", err)
//	}
//
//	// 2) Add a semi-transparent text watermark (center-bottom)
//	cfg := watermark.TextConfig{
//		Text:    "imgpipe • demo",
//		Opacity: 0.35,                           // transparency 0..1
//		RelX:    0.50,                           // anchor X (0..1), 0.5 = center
//		RelY:    0.92,                           // anchor Y (0..1), 0.92 ~ near bottom
//		FontPt:  0,                              // 0 = auto scale by image width
//		Color:   color.RGBA{255, 255, 255, 255}, // white
//		Padding: 0,                              // extra offset (px)
//	}
//
//	// Build a small pipeline just for the watermark step
//	watermarked, err := imageops.NewPipeline().
//		// You can attach more middlewares here if you want
//		Add(imageops.WatermarkText(cfg, 90)). // watermark + write as JPEG-quality 90 for this step
//		Run(in)
//	if err != nil {
//		log.Fatalf("watermark step failed: %v", err)
//	}
//
//	// 3) Convert format (e.g., to JPEG explicitly; could also be "png" or "webp" if you add webp support)
//	//    For PNG the quality is ignored (lossless); for JPEG it will be used by the backend.
//	converted, err := convert.Convert(watermarked, "jpeg", 85)
//	if err != nil {
//		log.Fatalf("format conversion failed: %v", err)
//	}
//
//	// 4) Compress further (tune quality to your need, e.g., 75 for visibly smaller size)
//	out, err := compress.Compress(converted, 75)
//	if err != nil {
//		log.Fatalf("final compression failed: %v", err)
//	}
//
//	// 5) Persist result
//	if err := os.WriteFile("testdata/output_converted.jpg", out, 0644); err != nil {
//		log.Fatalf("failed to write output image: %v", err)
//	}
//
//	log.Println("OK: watermark + format convert + compress -> testdata/output_converted.jpg")
//}
// update 0
// update 1

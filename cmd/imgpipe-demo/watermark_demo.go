package main

//
//import (
//	"image/color"
//	"log"
//	"os"
//
//	"github.com/HumbleLines/imgpipe/pkg/imageops"
//	"github.com/HumbleLines/imgpipe/pkg/watermark"
//)
//
//func main() {
//	in, err := os.ReadFile("testdata/input.jpg")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	cfg := watermark.TextConfig{
//		Text:    "imgpipe • demo",
//		Opacity: 0.35,                                       // translucent
//		RelX:    0.5,                                        // Center
//		RelY:    0.92,                                       // Center 8%
//		FontPt:  0,                                          // 0 = Automatically estimate by image width
//		Color:   color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White		Padding: 0,
//	}
//
//	out, err := imageops.NewPipeline().
//		Add(imageops.WatermarkText(cfg, 80)).
//		Run(in)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if err := os.WriteFile("testdata/output_watermark.jpg", out, 0644); err != nil {
//		log.Fatal(err)
//	}
//	log.Println("OK: wrote testdata/output_watermark.jpg")
//}
// update 6
// update 7
// update 11

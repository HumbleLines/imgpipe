# 📦 imgpipe

`imgpipe` is a lightweight, modular image processing toolkit written in Go.
It supports a **pipeline-based architecture**, allowing you to chain multiple image operations together for maximum flexibility and maintainability.

---

## ✨ Features

* **Compression** — Reduce image size while maintaining visual quality.
* **Watermark** — Add text watermarks with configurable position, opacity, font, and color.
* **Format Conversion** — Convert between image formats (JPEG, PNG, WebP, etc.).
* **Cropping** — Extract a specific region from the image.
* **Resizing** — Scale images with multiple fit strategies.
* **Rotation** — Rotate images by a given angle.
* **Border** — Add borders with custom color and thickness.

---

## 📂 Installation

```bash
go get github.com/HumbleLines/imgpipe
```

---

## 🚀 Usage Examples

All examples assume that you have an input file in `testdata/input.jpg`
and that the output will be written to `testdata/output.jpg`.

### 1. Compression

```go
package main

import (
	"log"
	"os"

	"github.com/HumbleLines/imgpipe/pkg/compress"
)

func main() {
	in, _ := os.ReadFile("testdata/input.jpg")
	out, err := compress.Compress(in, 80) // 80 = JPEG quality
	if err != nil {
		log.Fatal(err)
	}
	_ = os.WriteFile("testdata/output.jpg", out, 0644)
}
```

---

### 2. Watermark

```go
package main

import (
	"image/color"
	"log"
	"os"

	"github.com/HumbleLines/imgpipe/pkg/imageops"
	"github.com/HumbleLines/imgpipe/pkg/watermark"
)

func main() {
	in, _ := os.ReadFile("testdata/input.jpg")
	cfg := watermark.TextConfig{
		Text:    "imgpipe • demo",
		Opacity: 0.35,
		RelX:    0.5,
		RelY:    0.92,
		Color:   color.RGBA{255, 255, 255, 255},
	}
	out, err := imageops.NewPipeline().
		Add(imageops.WatermarkText(cfg, 80)).
		Run(in)
	if err != nil {
		log.Fatal(err)
	}
	_ = os.WriteFile("testdata/output.jpg", out, 0644)
}
```

---

### 3. Format Conversion

```go
package main

import (
	"log"
	"os"

	"github.com/HumbleLines/imgpipe/pkg/format"
)

func main() {
	in, _ := os.ReadFile("testdata/input.jpg")
	out, err := format.Convert(in, "png", 90) // Convert to PNG with quality 90
	if err != nil {
		log.Fatal(err)
	}
	_ = os.WriteFile("testdata/output.png", out, 0644)
}
```

---

### 4. Cropping

```go
package main

import (
	"log"
	"os"

	"github.com/HumbleLines/imgpipe/pkg/crop"
)

func main() {
	in, _ := os.ReadFile("testdata/input.jpg")
	out, err := crop.Crop(in, 100, 100, 400, 300) // x=100, y=100, width=400, height=300
	if err != nil {
		log.Fatal(err)
	}
	_ = os.WriteFile("testdata/output.jpg", out, 0644)
}
```

---

### 5. Resizing

```go
package main

import (
	"log"
	"os"

	"github.com/HumbleLines/imgpipe/pkg/resize"
)

func main() {
	in, _ := os.ReadFile("testdata/input.jpg")
	out, err := resize.Resize(in, 800, 450, resize.FitInside)
	if err != nil {
		log.Fatal(err)
	}
	_ = os.WriteFile("testdata/output.jpg", out, 0644)
}
```

---

### 6. Rotation

```go
package main

import (
	"log"
	"os"

	"github.com/HumbleLines/imgpipe/pkg/rotate"
)

func main() {
	in, _ := os.ReadFile("testdata/input.jpg")
	out, err := rotate.Rotate(in, 90) // Rotate 90 degrees clockwise
	if err != nil {
		log.Fatal(err)
	}
	_ = os.WriteFile("testdata/output.jpg", out, 0644)
}
```

---

### 7. Border

```go
package main

import (
	"image/color"
	"log"
	"os"

	"github.com/HumbleLines/imgpipe/pkg/border"
)

func main() {
	in, _ := os.ReadFile("testdata/input.jpg")
	out, err := border.AddBorder(in, 10, color.RGBA{255, 0, 0, 255}) // 10px red border
	if err != nil {
		log.Fatal(err)
	}
	_ = os.WriteFile("testdata/output.jpg", out, 0644)
}
```

---

## 🔗 Chaining Multiple Operations

Thanks to the pipeline-based design, you can combine multiple operations seamlessly:

```go
out, err := imageops.NewPipeline().
	Add(imageops.WatermarkText(cfg, 80)).
	Add(resize.ResizeHandler(800, 450, resize.FitInside)).
	Add(compress.CompressHandler(75)).
	Run(in)
```

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
// fix(log): fix broken logger format [2018-09-12 09:00:00]
// feat(core): initial project setup [2018-10-03 09:00:00]
// feat(core): initial project setup [2018-10-10 09:00:00]
// test(image): add fuzz tests [2019-01-02 09:00:00]
// chore(deps): update Go modules [2019-02-06 09:00:00]
// chore(build): add build script [2019-02-20 09:00:00]
// test(core): add edge case coverage [2019-03-13 09:00:00]
// refactor(image): extract resize logic [2019-04-17 09:00:00]
// docs(usage): add cli examples [2019-04-24 09:00:00]
// feat(cli): support output to stdout [2019-05-01 09:00:00]
// chore(deps): update Go modules [2019-05-22 09:00:00]
// chore(deps): update Go modules [2019-07-03 09:00:00]
// chore(build): add build script [2018-10-24 09:00:00]
// fix(log): fix broken logger format [2019-01-30 09:00:00]
// feat(cli): support output to stdout [2019-03-13 09:00:00]
// feat(core): initial project setup [2019-06-26 09:00:00]
// refactor(main): simplify flag parsing [2018-08-15 09:00:00]
// fix(cli): resolve path issues [2019-02-20 09:00:00]
// chore(build): add build script [2019-02-27 09:00:00]
// docs(readme): update usage examples [2019-04-17 09:00:00]
// feat(cli): support output to stdout [2019-05-29 09:00:00]
// change 0 at 2017-01-15T10:00:00
// update 0
// update 0
// update 0

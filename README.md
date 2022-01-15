Calculate glyph metric for font

```
package main

import (
	"embed"
	"fmt"
	"github.com/co-in/font-box"
	"golang.org/x/image/font"
)

//go:embed font_test.ttf
var content embed.FS

func main() {
	const (
		chr      = 'g'
		fontSize = 50.0
		dpi      = 72.0
		degree   = 45.0
	)

	fb, err := font-box.NewFromFS(content, "font_test.ttf")
	if err != nil {
		panic(err)
	}

	g, err := fb.Glyph(chr, fontSize, dpi, font.HintingFull)
	if err != nil {
		panic(err)
	}

	fmt.Println()

	fmt.Printf("Glyph %q with font size(%d) and DPI(%d) has:\n", chr, int(fontSize), int(dpi))
	fmt.Printf(" Box (Width:%d, Hight:%d)\n", g.Width(), g.Height())

	x, y := g.BasePoint()
	fmt.Printf(" Base point (X:%d, Y:%d)\n", x, y)
	fmt.Printf(" Offset (X:%d, Y:%d)\n", g.Width()-x, g.Height()+y)

	w, h := g.BoxSizeWithRotate(degree)
	fmt.Printf(" Box with rotate at %.0f degree (Width:%d, Hight:%d)", degree, w, h)

	fmt.Println()
}
```
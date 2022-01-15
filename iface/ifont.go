package iface

import (
	"golang.org/x/image/font"
	"image"
)

type IFont interface {
	Glyph(r rune, fontSize, dpi float64, hinting font.Hinting) (IFontBox, error)
}

type IFontBox interface {
	Render(angle float64, textColor *image.Uniform) (*image.NRGBA, error)
	BoxSizeWithRotate(angle float64) (width int, height int)
	Width() int
	Height() int
	BasePoint() (x int, y int)
}

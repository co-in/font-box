package font_metric

import (
	"errors"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"io/fs"
	"math"
)

const (
	halfCircleDegree = 180.0
	magicDPI         = 72
	magicDPI2        = 64.0
)

type cacheLevelHinting map[font.Hinting]*fontMetric
type cacheLevelDPI map[float64]cacheLevelHinting
type cacheLevelSize map[float64]cacheLevelDPI

type fontMetric struct {
	AdvancedWidth fixed.Int26_6
	Bounds        fixed.Rectangle26_6
}

func (m *fontMetric) BoxSizeWithRotate(angle float64) (width int, height int) {
	if angle < 0 {
		angle = halfCircleDegree - angle
	}

	t := float64(angle) * (math.Pi / halfCircleDegree)
	bx := float64(m.Width())
	by := float64(m.Height())
	cosA := math.Cos(t)
	sinA := math.Sin(t)
	width = int(math.Abs(math.Round(bx*cosA + by*sinA)))
	height = int(math.Abs(math.Round(bx*sinA + by*cosA)))

	return
}

func (m *fontMetric) Width() int {
	return (m.Bounds.Max.X - m.Bounds.Min.X).Ceil()
}

func (m *fontMetric) Height() int {
	return (m.Bounds.Max.Y - m.Bounds.Min.Y).Ceil()
}

func (m *fontMetric) BasePoint() (x int, y int) {
	x = m.Bounds.Min.X.Ceil()
	y = m.Bounds.Min.Y.Ceil()

	return
}

type fontBuf struct {
	fnt      *truetype.Font
	glyphBuf truetype.GlyphBuf
	cache    map[rune]cacheLevelSize
}

func New(fnt *truetype.Font) *fontBuf {
	return &fontBuf{
		fnt:   fnt,
		cache: make(map[rune]cacheLevelSize),
	}
}

func NewFromFS(fs fs.ReadFileFS, filename string) (*fontBuf, error) {
	fontBytes, err := fs.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if len(fontBytes) == 0 {
		return nil, errors.New("invalid font")
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	return New(f), nil
}

func (m *fontBuf) Metric(r rune, fontSize, dpi float64, hinting font.Hinting) (*fontMetric, error) {
	var ok bool
	_, ok = m.cache[r]
	if !ok {
		m.cache[r] = cacheLevelSize{}
	}

	_, ok = m.cache[r][fontSize]
	if !ok {
		m.cache[r][fontSize] = cacheLevelDPI{}
	}

	_, ok = m.cache[r][fontSize][dpi]
	if !ok {
		m.cache[r][fontSize][dpi] = cacheLevelHinting{}
	}

	_, ok = m.cache[r][fontSize][dpi][hinting]
	if !ok {
		glyph := m.fnt.Index(r)
		scale := fixed.Int26_6(fontSize * dpi * (magicDPI2 / magicDPI))
		err := m.glyphBuf.Load(m.fnt, scale, glyph, hinting)
		if err != nil {
			return nil, err
		}

		m.cache[r][fontSize][dpi][hinting] = &fontMetric{
			AdvancedWidth: m.glyphBuf.AdvanceWidth,
			Bounds:        m.glyphBuf.Bounds,
		}
	}

	return m.cache[r][fontSize][dpi][hinting], nil
}

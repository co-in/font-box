package font_metric_test

import (
	"embed"
	"font-metric"
	"golang.org/x/image/font"
	"gopkg.in/check.v1"
	"testing"
)

//go:embed font_test.ttf
var content embed.FS

type fontSuite struct{}

var _ = check.Suite(&fontSuite{})

func Test(t *testing.T) {
	check.TestingT(t)
}

var tests = []struct {
	chr     rune
	width   int
	height  int
	size    float64
	dpi     float64
	hRotate int
	wRotate int
	angle   float64
	x       int
	y       int
}{
	{chr: 'A', width: 31, height: 35, size: 50, dpi: 72, wRotate: 46, hRotate: 47, angle: 40, x: 0, y: 0},
	{chr: 'j', width: 11, height: 43, size: 50, dpi: 72, wRotate: 36, hRotate: 40, angle: 40, x: 0, y: -9},
	{chr: 'j', width: 11, height: 43, size: 50, dpi: 72, wRotate: 36, hRotate: 40, angle: -40, x: 0, y: -9},
	{chr: 'W', width: 32, height: 24, size: 50, dpi: 50, wRotate: 40, hRotate: 38, angle: 33, x: 0, y: 0},
	{chr: 'g', width: 23, height: 33, size: 50, dpi: 72, wRotate: 39, hRotate: 40, angle: 40, x: 2, y: -9},
	{chr: 'l', width: 17, height: 46, size: 66, dpi: 72, wRotate: 43, hRotate: 46, angle: 40, x: 0, y: 0},
}

func (s *fontSuite) TestGlyphMetric(c *check.C) {
	fnt, err := font_metric.NewFromFS(content, "font_test.ttf")
	c.Assert(err, check.IsNil)

	for _, test := range tests {
		g, err := fnt.Metric(test.chr, test.size, test.dpi, font.HintingFull)
		c.Assert(err, check.IsNil)
		c.Assert(g.Width(), check.Equals, test.width)
		c.Assert(g.Height(), check.Equals, test.height)

		x, y := g.BasePoint()
		c.Assert(x, check.Equals, test.x)
		c.Assert(y, check.Equals, test.y)

		wR, hR := g.BoxSizeWithRotate(test.angle)
		c.Assert(wR, check.Equals, test.wRotate)
		c.Assert(hR, check.Equals, test.hRotate)
	}
}

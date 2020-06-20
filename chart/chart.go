// Package chart makes charts using the gio canvas
package chart

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"strconv"
	"strings"

	gc "github.com/ajstarks/giocanvas"
)

// NameValue is a name,value pair
type NameValue struct {
	label string
	note  string
	value float64
}

// ChartBox holds the essential data for making a chart
type ChartBox struct {
	Title                    string
	Data                     []NameValue
	Color                    color.RGBA
	Top, Bottom, Left, Right float64
	Minvalue, Maxvalue       float64
	Zerobased                bool
}

const (
	largest  = 1.797693134862315708145274237317043567981e+308
	smallest = -largest
)

var labelcolor = color.RGBA{100, 100, 100, 255}

// DataRead reads tab separated values into a ChartBox
// default values for the top, bottom, left, right (90,50,10,90) are filled in
// as is the default color, black
func DataRead(r io.Reader) (ChartBox, error) {
	var d NameValue
	var data []NameValue
	var err error
	maxval := smallest
	minval := largest
	title := ""
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) == 0 { // skip blank lines
			continue
		}
		if t[0] == '#' && len(t) > 2 { // process titles
			title = strings.TrimSpace(t[1:])
			continue
		}
		fields := strings.Split(t, "\t")
		if len(fields) < 2 {
			continue
		}
		if len(fields) == 3 {
			d.note = fields[2]
		} else {
			d.note = ""
		}
		d.label = fields[0]
		d.value, err = strconv.ParseFloat(fields[1], 64)
		if err != nil {
			d.value = 0
		}
		if d.value > maxval {
			maxval = d.value
		}
		if d.value < minval {
			minval = d.value
		}
		data = append(data, d)
	}
	err = scanner.Err()
	return ChartBox{
		Title:     title,
		Data:      data,
		Minvalue:  minval,
		Maxvalue:  maxval,
		Color:     color.RGBA{0, 0, 0, 255},
		Left:      10,
		Right:     90,
		Top:       90,
		Bottom:    50,
		Zerobased: true,
	}, err
}

// zerobase uses the correct base for scaling
func zerobase(usez bool, n float64) float64 {
	if usez {
		return 0
	}
	return n
}

// Bar makes a (column) bar chart
func (c *ChartBox) Bar(canvas *gc.Canvas, size float64) {
	dlen := float64(len(c.Data) - 1)
	ymin := zerobase(c.Zerobased, c.Minvalue)
	for i, d := range c.Data {
		x := float32(gc.MapRange(float64(i), 0, dlen, c.Left, c.Right))
		y := float32(gc.MapRange(d.value, ymin, c.Maxvalue, c.Bottom, c.Top))
		canvas.Line(float32(x), float32(c.Bottom), x, y, float32(size), c.Color)
	}
}

// HBar makes a horizontal bar chart
func (c *ChartBox) HBar(canvas *gc.Canvas, size, linespacing, textsize float64) {
	y := float32(c.Top)
	cl := float32(c.Left)
	xmin := zerobase(c.Zerobased, c.Minvalue)
	for _, d := range c.Data {
		canvas.EText(cl-2, y-float32(size/2), float32(textsize), d.label, labelcolor)
		x2 := gc.MapRange(d.value, xmin, c.Maxvalue, c.Left, c.Right)
		canvas.Line(cl, y, float32(x2), y, float32(size), c.Color)
		y -= float32(linespacing)
	}
}

// Line makes a line chart
func (c *ChartBox) Line(canvas *gc.Canvas, size float64) {
	n := len(c.Data)
	fn := float64(n - 1)
	ymin := zerobase(c.Zerobased, c.Minvalue)
	for i := 0; i < n-1; i++ {
		v1 := c.Data[i].value
		v2 := c.Data[i+1].value
		x1 := float32(gc.MapRange(float64(i), 0, fn, c.Left, c.Right))
		y1 := float32(gc.MapRange(v1, ymin, c.Maxvalue, c.Bottom, c.Top))
		x2 := float32(gc.MapRange(float64(i+1), 0, fn, c.Left, c.Right))
		y2 := float32(gc.MapRange(v2, ymin, c.Maxvalue, c.Bottom, c.Top))
		canvas.Line(x1, y1, x2, y2, float32(size), c.Color)
	}
}

// Area makes a area chart with specified opacity
func (c *ChartBox) Area(canvas *gc.Canvas, opacity float64) {
	n := len(c.Data)
	ymin := zerobase(c.Zerobased, c.Minvalue)
	width := c.Right
	height := c.Top
	x := c.Left
	y := c.Bottom
	ax := make([]float32, n+2)
	ay := make([]float32, n+2)
	ax[0] = float32(x)
	ay[0] = float32(y)
	ax[n+1] = float32(width)
	ay[n+1] = float32(y)

	for i, d := range c.Data {
		xp := float32(gc.MapRange(float64(i), 0, float64(n-1), float64(x), float64(width)))
		yp := float32(gc.MapRange(d.value, ymin, c.Maxvalue, float64(y), float64(height)))
		ax[i+1] = xp
		ay[i+1] = yp
	}
	c.Color.A = uint8(255.0 * (opacity / 100))
	canvas.Polygon(ax, ay, c.Color)
}

// Label draws the x axis labels
func (c *ChartBox) Label(canvas *gc.Canvas, size float64, n int) {
	fn := float64(len(c.Data) - 1)
	for i, d := range c.Data {
		x := float32(gc.MapRange(float64(i), 0, fn, c.Left, c.Right))
		if i%n == 0 {
			canvas.CText(x, float32(c.Bottom-(size*2)), float32(size), d.label, c.Color)
		}
	}
}

// Scatter makes a scatter chart
func (c *ChartBox) Scatter(canvas *gc.Canvas, size float64) {
	dlen := float64(len(c.Data) - 1)
	ymin := zerobase(c.Zerobased, c.Minvalue)
	for i, d := range c.Data {
		x := float32(gc.MapRange(float64(i), 0, dlen, c.Left, c.Right))
		y := float32(gc.MapRange(d.value, ymin, c.Maxvalue, c.Bottom, c.Top))
		canvas.Circle(x, y, float32(size), c.Color)
	}
}

// Grid makes a grid
func Grid(canvas *gc.Canvas, left, bottom, width, height, size float64, color color.RGBA) {
	for x := float32(left); x <= float32(left+width); x += float32(size) {
		canvas.Line(x, float32(bottom), x, float32(bottom+height), 0.1, color)
	}
	for y := float32(bottom); y <= float32(bottom+height); y += float32(size) {
		canvas.Line(float32(left), y, float32(left+width), y, 0.2, color)
	}
}

// YAxis makes the Y axis with optional grid lines
func (c *ChartBox) YAxis(canvas *gc.Canvas, size, min, max, step float64, format string, gridlines bool) {
	w := c.Right - c.Left
	ymin := zerobase(c.Zerobased, c.Minvalue)
	for v := min; v <= max; v += step {
		y := float32(gc.MapRange(v, ymin, c.Maxvalue, c.Bottom, c.Top))
		if gridlines {
			canvas.Line(float32(c.Left), y, float32(c.Left+w), y, 0.05, color.RGBA{128, 128, 128, 255})
		}
		canvas.EText(float32(c.Left-2), (y - float32(size/3)), float32(size), fmt.Sprintf(format, v), c.Color)
	}
}

// CTitle makes a centered title
func (c *ChartBox) CTitle(canvas *gc.Canvas, size, offset float64) {
	midx := c.Left + ((c.Right - c.Left) / 2)
	canvas.CText(float32(midx), float32(c.Top+offset), float32(size), c.Title, c.Color)
}

// Frame makes a filled frame with the specified opacity (0-100)
func (c *ChartBox) Frame(canvas *gc.Canvas, op float64) {
	a := c.Color.A // Save opacity
	w := float32(c.Right - c.Left)
	h := float32(c.Top - c.Bottom)
	fa := uint8((op / 100) * 255.0)
	c.Color.A = fa
	canvas.Rect(float32(c.Left)+w/2, float32(c.Bottom)+h/2, w, h, c.Color)
	c.Color.A = a // Restore opacity
}

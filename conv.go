package giocanvas

import (
	"fmt"
	"image/color"
	"math"
)

// Convenience functions

// MapRange maps a value between low1 and high1, return the corresponding value between low2 and high2
func MapRange(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

// PolarDegrees returns the Cartesian coordinates (x, y) from polar coordinates
// with compensation for canvas aspect ratio
// center at (cx, cy), radius r, and angle theta (degrees)
func (c *Canvas) PolarDegrees(cx, cy, r, theta float32) (float32, float32) {
	fr := float64(r)
	ft := float64(theta * (math.Pi / 180))
	aspect := float64(c.Width / c.Height)
	px := fr * math.Cos(ft)
	py := (fr * aspect) * math.Sin(ft)
	return cx + float32(px), cy + float32(py)

}

// Polar returns the Cartesian coordinates (x, y) from polar coordinates
// with compensation for canvas aspect ratio
// center at (cx, cy), radius r, and angle theta (radians)
func (c *Canvas) Polar(cx, cy, r, theta float32) (float32, float32) {
	fr := float64(r)
	ft := float64(theta)
	aspect := float64(c.Width / c.Height)
	px := fr * math.Cos(ft)
	py := (fr * aspect) * math.Sin(ft)
	return cx + float32(px), cy + float32(py)
}

// Background makes a filled rectangle covering the whole canvas
func (c *Canvas) Background(fillcolor color.RGBA) {
	c.AbsRect(0, 0, c.Width, c.Height, fillcolor)
}

// Coord shows the specified coordinate, using percentage-based coordinates
// the (x, y) label is above the point, with a label below
func (c *Canvas) Coord(x, y, size float32, s string, fillcolor color.RGBA) {
	c.Square(x, y, size/2, fillcolor)
	c.TextMid(x, y+size, size, fmt.Sprintf("(%g, %g)", x, y), fillcolor)
	if len(s) > 0 {
		c.CText(x, y-(size*1.33), size*0.66, s, fillcolor)
	}
}

// Grid makes vertical and horizontal grid lines, percentage-based coordinates
func (c *Canvas) Grid(x, y, w, h, size, interval float32, linecolor color.RGBA) {
	for xp := x; xp <= x+w; xp += interval {
		c.Line(xp, y, xp, y+h, size, linecolor) // vertical line
	}
	for yp := y; yp <= y+h; yp += interval {
		c.Line(x, yp, x+w, yp, size, linecolor) // horizontal line
	}
}

package giocanvas

import (
	"image"
	"image/color"

	"gioui.org/text"
)

// Methods using percentage-based, (x and y range from 0-100),
// classical coordinate system (x increasing left to right, y increasing bottom to top)

// Lines and shapes

// Line makes a stroked line using percentage-based measures
// from (x0, y0) to (x1, y1), stroke width size
func (c *Canvas) Line(x0, y0, x1, y1, size float32, strokecolor color.NRGBA) {
	x0, y0 = dimen(x0, y0, c.Width, c.Height)
	x1, y1 = dimen(x1, y1, c.Width, c.Height)
	size = pct(size, c.Width)
	c.AbsLine(x0, y0, x1, y1, size, strokecolor)
}

// VLine makes a vertical line beginning at (x,y) with dimension (w, h)
// the line begins at (x,y) and moves upward by linewidth
func (c *Canvas) VLine(x, y, lineheight, size float32, linecolor color.NRGBA) {
	c.Line(x, y, x, y+lineheight, size, linecolor)
}

// HLine makes a horizontal line starting at (x, y), with dimensions (w, h)
// the line begin at (x,y) and extends to the left by linewidth
func (c *Canvas) HLine(x, y, linewidth, size float32, linecolor color.NRGBA) {
	c.Line(x, y, x+linewidth, y, size, linecolor)
}

// Polygon makes a filled polygon using percentage-based measures
// vertices in x and y,
func (c *Canvas) Polygon(x, y []float32, fillcolor color.NRGBA) {
	if len(x) != len(y) || len(x) < 3 {
		return
	}
	nx := make([]float32, len(x))
	ny := make([]float32, len(y))
	for i := 0; i < len(x); i++ {
		nx[i], ny[i] = dimen(x[i], y[i], c.Width, c.Height)
	}
	c.AbsPolygon(nx, ny, fillcolor)
}

// QuadCurve makes a quadric Bezier curve, using percentage-based measures
// starting at (x, y), control point at (cx, cy), end point (ex, ey)
func (c *Canvas) QuadCurve(x, y, cx, cy, ex, ey float32, fillcolor color.NRGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	cx, cy = dimen(cx, cy, c.Width, c.Height)
	ex, ey = dimen(ex, ey, c.Width, c.Height)
	c.AbsQuadBezier(x, y, cx, cy, ex, ey, 0, fillcolor)
}

// Curve makes a quadric Bezier curve, using percentage-based measures
// starting at (x, y), control point at (cx, cy), end point (ex, ey)
func (c *Canvas) Curve(x, y, cx, cy, ex, ey float32, fillcolor color.NRGBA) {
	c.QuadCurve(x, y, cx, cy, ex, ey, fillcolor)
}

// CubeCurve makes a cubic Bezier curve, using percentage-based measures
// starting at (x, y), control points at (cx1, cy1), (cx2, cy2), end point (ex, ey)
func (c *Canvas) CubeCurve(x, y, cx1, cy1, cx2, cy2, ex, ey float32, fillcolor color.NRGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	cx1, cy1 = dimen(cx1, cy1, c.Width, c.Height)
	cx2, cy2 = dimen(cx2, cy2, c.Width, c.Height)
	ex, ey = dimen(ex, ey, c.Width, c.Height)
	c.AbsCubicBezier(x, y, cx1, cy1, cx2, cy2, ex, ey, 0, fillcolor)
}

// Circle makes a filled circle, using percentage-based measures
// center is (x,y), radius r
func (c *Canvas) Circle(x, y, r float32, fillcolor color.NRGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	r = pct(r, c.Width)
	c.AbsCircle(x, y, r, fillcolor)
}

// Ellipse makes a filled circle, using percentage-based measures
// center is (x,y), radii (w, h)
func (c *Canvas) Ellipse(x, y, w, h float32, fillcolor color.NRGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	w = pct(w, c.Width)
	h = pct(h, c.Height)
	c.AbsEllipse(x, y, w, h, fillcolor)
}

// Arc makes a filled arc, using percentage-based measures
// center is (x, y) the arc begins at angle a1, and ends at a2, with radius r.
// The arc is filled with the specified color.
func (c *Canvas) Arc(x, y, r float32, a1, a2 float64, fillcolor color.NRGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	pr := pct(r, c.Width)
	c.AbsArc(x, y, pr, a1, a2, fillcolor)
}

// ArcLine makes a stroked arc, using percentage-based measures
// center is (x, y), the arc begins at angle a1, and ends at a2, with radius r.
// The arc is stroked with the specified stroke size and color
func (c *Canvas) ArcLine(x, y, r float32, a1, a2 float64, size float32, fillcolor color.NRGBA) {
	step := (a2 - a1) / 100
	x1, y1 := c.Polar(x, y, r, float32(a1))
	for t := a1 + step; t <= a2; t += step {
		x2, y2 := c.Polar(x, y, r, float32(t))
		c.Line(x1, y1, x2, y2, size, fillcolor)
		x1 = x2
		y1 = y2
	}
}

// Text methods

// Text places text using percentage-based measures
// left at x, baseline at y, at the specified size and color
func (c *Canvas) Text(x, y, size float32, s string, fillcolor color.NRGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	size = pct(size, c.Width)
	c.textops(x, y, size, text.Start, s, fillcolor)
}

// TextEnd places text using percentage-based measures
// x is the end of the string, baseline at y, using specified size and color
func (c *Canvas) TextEnd(x, y, size float32, s string, fillcolor color.NRGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	size = pct(size, c.Width)
	c.textops(x, y, size, text.End, s, fillcolor)
}

// TextMid places text using percentage-based measures
// text is centered at x, baseline y, using specied size and color
func (c *Canvas) TextMid(x, y, size float32, s string, fillcolor color.NRGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	size = pct(size, c.Width)
	c.textops(x, y, size, text.Middle, s, fillcolor)
}

// EText - alternative name for TextEnd
func (c *Canvas) EText(x, y, size float32, s string, fillcolor color.NRGBA) {
	c.TextEnd(x, y, size, s, fillcolor)
}

// CText - alternative name for TextMid
func (c *Canvas) CText(x, y, size float32, s string, fillcolor color.NRGBA) {
	c.TextMid(x, y, size, s, fillcolor)
}

// TextWrap places and wraps text using percentage-based measures
// text begins at (x,y), baseline y, and wraps at width, using specied size and color
func (c *Canvas) TextWrap(x, y, size, width float32, s string, fillcolor color.NRGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	size = pct(size, c.Width)
	width = pct(width, c.Width)
	c.AbsTextWrap(x, y, size, width, s, fillcolor)
}

// Rect makes a rectangle using percentage-based measures
// upper left corner at (x,y), with size at (w,h)
func (c *Canvas) Rect(x, y, w, h float32, fillcolor color.NRGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	w = pct(w, c.Width)
	h = pct(h, c.Height)
	c.AbsCenterRect(x, y, w, h, fillcolor)
}

// CornerRect makes a rectangle using percentage-based measures
// upper left corner at (x,y), with sized at (w,h)
func (c *Canvas) CornerRect(x, y, w, h float32, fillcolor color.NRGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	w = pct(w, c.Width)
	h = pct(h, c.Height)
	c.AbsRect(x, y, w, h, fillcolor)
}

// Square makes a square shape, using percentage based measures
// centered at (x, y), sides are w. Accounts for screen aspect
func (c *Canvas) Square(x, y, w float32, fillcolor color.NRGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	w = pct(w, c.Height)
	h := pct(100, w)
	c.AbsCenterRect(x, y, w, h, fillcolor)
}

// CenterRect makes a rectangle using percentage-based measures
// with center at (x,y), sized at (w,h)
func (c *Canvas) CenterRect(x, y, w, h float32, fillcolor color.NRGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	w = pct(w, c.Width)
	h = pct(h, c.Height)
	c.AbsCenterRect(x, y, w, h, fillcolor)
}

// Images

// Img places a scaled image centered at (x, y), data from image.Image
// using percentage coordinates and scales
func (c *Canvas) Img(im image.Image, x, y float32, w, h int, scale float32) {
	x, y = dimen(x, y, c.Width, c.Height)
	c.AbsImg(im, x, y, w, h, scale)
}

// Image places a scaled image centered at (x,y), reading from a named file,
// using percetage coordinates and scales
func (c *Canvas) Image(name string, x, y float32, w, h int, scale float32) {
	c.CenterImage(name, x, y, w, h, scale)
}

// CenterImage places a scaled image centered at (x,y),
// using percentage coordinates and scales
func (c *Canvas) CenterImage(name string, x, y float32, w, h int, scale float32) {
	x, y = dimen(x, y, c.Width, c.Height)
	c.AbsCenterImage(name, x, y, w, h, scale)
}

// pct returns the percentage of its input
func pct(p float32, m float32) float32 {
	return ((p / 100.0) * m)
}

// dimen returns canvas dimensions from percentages
// (converting from x increasing left-right, y increasing top-bottom)
func dimen(xp, yp, w, h float32) (float32, float32) {
	return pct(xp, w), pct(100-yp, h)
}

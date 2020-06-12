// Package giocanvas is a 2D canvas API built on gio
package giocanvas

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif" // needed by image
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"

	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// Canvas defines the Gio canvas
type Canvas struct {
	Width, Height float32
	TextColor     color.RGBA
	Context       layout.Context
}

// NewCanvas initializes a Canvas
func NewCanvas(width, height float32, e system.FrameEvent) *Canvas {
	canvas := new(Canvas)
	canvas.Width = width
	canvas.Height = height
	canvas.TextColor = color.RGBA{0, 0, 0, 255}
	canvas.Context = layout.NewContext(new(op.Ops), e)
	return canvas
}

// Convenience functions

// MapRange maps a value between low1 and high1, return the corresponding value between low2 and high2
func MapRange(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
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

// Methods using percentage-based, (x and y range from 0-100),
// classical coordinate system (x increasing left to right, y increasing bottom to top)

// Lines and shapes

// Line makes a stroked line using percentage-based measures
// from (x0, y0) to (x1, y1), stroke width size
func (c *Canvas) Line(x0, y0, x1, y1, size float32, strokecolor color.RGBA) {
	x0, y0 = dimen(x0, y0, c.Width, c.Height)
	x1, y1 = dimen(x1, y1, c.Width, c.Height)
	size = pct(size, c.Width)
	c.AbsLine(x0, y0, x1, y1, size, strokecolor)
}

// oLine (deprecated) makes a stroked line using percentage-based measures
// from (x0, y0) to (x1, y1), stroke width size
func (c *Canvas) oLine(x0, y0, x1, y1, size float32, fillcolor color.RGBA) {
	x0, y0 = dimen(x0, y0, c.Width, c.Height)
	x1, y1 = dimen(x1, y1, c.Width, c.Height)
	size = pct(size, c.Width)
	c.oAbsLine(x0, y0, x1, y1, size, fillcolor)
}

// VLine makes a vertical line beginning at (x,y) with dimension (w, h)
// the line begins at (x,y) and moves upward by linewidth
func (c *Canvas) VLine(x, y, lineheight, size float32, linecolor color.RGBA) {
	c.Line(x, y, x, y+lineheight, size, linecolor)
}

// HLine makes a horizontal line starting at (x, y), with dimensions (w, h)
// the line begin at (x,y) and extends to the left by linewidth
func (c *Canvas) HLine(x, y, linewidth, size float32, linecolor color.RGBA) {
	c.Line(x, y, x+linewidth, y, size, linecolor)
}

// Polygon makes a filled polygon using percentage-based measures
// vertices in x and y,
func (c *Canvas) Polygon(x, y []float32, fillcolor color.RGBA) {
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

// Curve makes a quadric Bezier curve, using percentage-based measures
// starting at (x, y), control point at (cx, cy), end point (ex, ey)
func (c *Canvas) Curve(x, y, cx, cy, ex, ey float32, fillcolor color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	cx, cy = dimen(cx, cy, c.Width, c.Height)
	ex, ey = dimen(ex, ey, c.Width, c.Height)
	c.AbsQuadBezier(x, y, cx, cy, ex, ey, 0, fillcolor)
}

// CubeCurve makes a cubic Bezier curve, using percentage-based measures
// starting at (x, y), control points at (cx1, cy1), (cx2, cy2), end point (ex, ey)
func (c *Canvas) CubeCurve(x, y, cx1, cy1, cx2, cy2, ex, ey float32, fillcolor color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	cx1, cy1 = dimen(cx1, cy1, c.Width, c.Height)
	cx2, cy2 = dimen(cx2, cy2, c.Width, c.Height)
	ex, ey = dimen(ex, ey, c.Width, c.Height)
	c.AbsCubicBezier(x, y, cx1, cy1, cx2, cy2, ex, ey, 0, fillcolor)
}

// Circle makes a filled circle, using percentage-based measures
// center is (x,y), radius r
func (c *Canvas) Circle(x, y, r float32, fillcolor color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	r = pct(r, c.Width)
	c.AbsCircle(x, y, r, fillcolor)
}

// Ellipse makes a filled circle, using percentage-based measures
// center is (x,y), radii (w, h)
func (c *Canvas) Ellipse(x, y, w, h float32, fillcolor color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	w = pct(w, c.Width)
	h = pct(h, c.Height)
	c.AbsEllipse(x, y, w, h, fillcolor)
}

// Arc makes a filled arc, using percentage-based measures
// center is (x, y) the arc begins at angle a1, and ends at a2
// TODO: Still buggy
func (c *Canvas) Arc(x, y, r float32, a1, a2 float64, fillcolor color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	pr := pct(r, c.Width)
	c.AbsArc(float64(x), float64(y), float64(pr), float64(a1), float64(a2), fillcolor)
}

// Text methods

// Text places text using percentage-based measures
// left at x, baseline at y, at the specified size and color
func (c *Canvas) Text(x, y, size float32, s string, fillcolor color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	size = pct(size, c.Width)
	c.textops(x, y, size, text.Start, s, fillcolor)
}

// TextEnd places text using percentage-based measures
// x is the end of the string, baseline at y, using specified size and color
func (c *Canvas) TextEnd(x, y, size float32, s string, fillcolor color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	size = pct(size, c.Width)
	c.textops(x, y, size, text.End, s, fillcolor)
}

// TextMid places text using percentage-based measures
// text is centered at x, baseline y, using specied size and color
func (c *Canvas) TextMid(x, y, size float32, s string, fillcolor color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	size = pct(size, c.Width)
	c.textops(x, y, size, text.Middle, s, fillcolor)
}

// EText - alternative name for TextEnd
func (c *Canvas) EText(x, y, size float32, s string, fillcolor color.RGBA) {
	c.TextEnd(x, y, size, s, fillcolor)
}

// CText - alternarive name for TextMid
func (c *Canvas) CText(x, y, size float32, s string, fillcolor color.RGBA) {
	c.TextMid(x, y, size, s, fillcolor)
}

// TextWrap places and wraps text using percentage-based measures
// text begins at (x,y), baseline y, and wraps at width, using specied size and color
func (c *Canvas) TextWrap(x, y, size, width float32, s string, fillcolor color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	size = pct(size, c.Width)
	width = pct(width, c.Width)
	c.AbsTextWrap(x, y, size, width, s, fillcolor)
}

// Rect makes a rectangle using percentage-based measures
// upper left corner at (x,y), with size at (w,h)
func (c *Canvas) Rect(x, y, w, h float32, fillcolor color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	w = pct(w, c.Width)
	h = pct(h, c.Height)
	c.AbsRect(x, y, w, h, fillcolor)
}

// CornerRect makes a rectangle using percentage-based measures
// upper left corner at (x,y), with sized at (w,h)
func (c *Canvas) CornerRect(x, y, w, h float32, fillcolor color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	w = pct(w, c.Width)
	h = pct(h, c.Height)
	c.AbsRect(x, y, w, h, fillcolor)
}

// Square makes a square shape, using percentage based measures
// centered at (x, y), sides are w. Accounts for screen aspect
func (c *Canvas) Square(x, y, w float32, fillcolor color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	w = pct(w, c.Height)
	h := pct(100, w)
	c.AbsCenterRect(x, y, w, h, fillcolor)
}

// CenterRect makes a rectangle using percentage-based measures
// with upper left corner at (x,y), with sized at (w,h)
func (c *Canvas) CenterRect(x, y, w, h float32, fillcolor color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	w = pct(w, c.Width)
	h = pct(h, c.Height)
	c.AbsCenterRect(x, y, w, h, fillcolor)
}

// Images

// Image places a scaled image centered at (x,y)
func (c *Canvas) Image(name string, x, y float32, w, h int, scale float32) {
	c.CenterImage(name, x, y, w, h, scale)
}

// CenterImage places a scaled image centered at (x,y)
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

// Foundational methods, and methods using Gio standard coordinates

// textops places text
func (c *Canvas) textops(x, y, size float32, alignment text.Alignment, s string, fillcolor color.RGBA) {
	offset := x
	switch alignment {
	case text.End:
		offset = x - c.Width
	case text.Middle:
		offset = x - c.Width/2
	}
	defer op.Push(c.Context.Ops).Pop()
	op.TransformOp{}.Offset(f32.Point{X: offset, Y: y - size}).Add(c.Context.Ops) // shift to use baseline
	l := material.Label(material.NewTheme(gofont.Collection()), unit.Px(size), s)
	l.Color = fillcolor
	l.Alignment = alignment
	l.Layout(c.Context)
}

// AbsTextWrap places and wraps text at (x, y), wrapped at width
func (c *Canvas) AbsTextWrap(x, y, size, width float32, s string, fillcolor color.RGBA) {
	defer op.Push(c.Context.Ops).Pop()
	op.TransformOp{}.Offset(f32.Point{X: x, Y: y - size}).Add(c.Context.Ops) // shift to use baseline
	l := material.Label(material.NewTheme(gofont.Collection()), unit.Px(size), s)
	l.Color = fillcolor
	c.Context.Constraints.Max.X = int(width)
	l.Layout(c.Context)
	c.Context.Constraints.Max.X = int(c.Width) // restore width...
}

// AbsText places text at (x,y)
func (c *Canvas) AbsText(x, y, size float32, s string, fillcolor color.RGBA) {
	c.textops(x, y, size, text.Start, s, fillcolor)
}

// AbsTextMid places text centered at (x,y)
func (c *Canvas) AbsTextMid(x, y, size float32, s string, fillcolor color.RGBA) {
	c.textops(x, y, size, text.Middle, s, fillcolor)
}

// AbsTextEnd places text aligned to the end
func (c *Canvas) AbsTextEnd(x, y, size float32, s string, fillcolor color.RGBA) {
	c.textops(x, y, size, text.End, s, fillcolor)
}

// AbsRect makes a filled Rectangle; left corner at (x, y), with dimensions (w,h)
func (c *Canvas) AbsRect(x, y, w, h float32, fillcolor color.RGBA) {
	ops := c.Context.Ops
	r := f32.Rect(x, y+h, x+w, y)
	paint.ColorOp{Color: fillcolor}.Add(ops)
	paint.PaintOp{Rect: r}.Add(ops)
}

// AbsCenterRect makes a filled rectangle centered at (x, y), with dimensions (w,h)
func (c *Canvas) AbsCenterRect(x, y, w, h float32, fillcolor color.RGBA) {
	ops := c.Context.Ops
	r := f32.Rect(x-(w/2), y+(h/2), x+(w/2), y-(h/2))
	paint.ColorOp{Color: fillcolor}.Add(ops)
	paint.PaintOp{Rect: r}.Add(ops)
}

// AbsVLine makes a vertical line beginning at (x,y) with dimension (w, h)
func (c *Canvas) AbsVLine(x, y, w, h float32, fillcolor color.RGBA) {
	c.AbsLine(x, y, x, y+h, w, fillcolor)
}

// AbsHLine makes a horizontal line starting at (x, y), with dimensions (w, h)
func (c *Canvas) AbsHLine(x, y, w, h float32, fillcolor color.RGBA) {
	c.AbsLine(x, y, x+w, y, h, fillcolor)
}

// AbsGrid uses horizontal and vertical lines to make a grid
func (c *Canvas) AbsGrid(width, height, size, interval float32, fillcolor color.RGBA) {
	var x, y float32
	for y = 0; y <= height; y += height / interval {
		c.AbsHLine(0, y, width, size, fillcolor)
	}
	for x = 0; x <= width; x += width / interval {
		c.AbsVLine(x, 0, size, height, fillcolor)
	}
}

// AbsCenterImage places a named image centered at (x, y)
// using the specified dimensions (w, h), and hen scaled
func (c *Canvas) AbsCenterImage(name string, x, y float32, w, h int, scale float32) {
	r, err := os.Open(name)
	if err != nil {
		return
	}
	im, _, err := image.Decode(r)
	if err != nil {
		return
	}
	c.AbsImg(im, x, y, w, h, scale)
}

// AbsImg places a image.Image centered at (x, y)
// using the specified dimensions (w, h), and then scaled
func (c *Canvas) AbsImg(im image.Image, x, y float32, w, h int, scale float32) {
	// compute scaled image dimensions
	// if w and h are zero, use the natural dimensions
	sc := scale / 100
	imw := float32(w) * sc
	imh := float32(h) * sc
	if w == 0 && h == 0 {
		b := im.Bounds()
		imw = float32(b.Max.X) * sc
		imh = float32(b.Max.Y) * sc
	}
	// center the image
	x = x - (imw / 2)
	y = y - (imh / 2)
	ops := c.Context.Ops
	paint.NewImageOp(im).Add(ops)
	paint.PaintOp{Rect: f32.Rect(x, y, x+imw, y+imh)}.Add(ops)
}

// AbsPolygon makes a closed, filled polygon with vertices in x and y
func (c *Canvas) AbsPolygon(x, y []float32, fillcolor color.RGBA) {
	if len(x) != len(y) {
		return
	}
	path := new(clip.Path)
	ops := c.Context.Ops
	r := f32.Rect(0, 0, c.Width, c.Height)

	defer op.Push(c.Context.Ops).Pop()
	path.Begin(ops)
	path.Move(f32.Point{X: x[0], Y: y[0]})

	l := len(x)
	point := f32.Point{}
	for i := 1; i < l; i++ {
		path.Line(f32.Point{X: x[i] - x[i-1], Y: y[i] - y[i-1]})
	}
	path.Line(f32.Point{X: x[0] - x[l-1], Y: y[0] - y[l-1]})
	path.Line(point)
	path.End().Add(ops)
	paint.ColorOp{Color: fillcolor}.Add(ops)
	paint.PaintOp{Rect: r}.Add(ops)
}

// AbsLine makes a line from (x0,y0) to (x1, y1) using absolute coordinates
func (c *Canvas) AbsLine(x0, y0, x1, y1, size float32, fillcolor color.RGBA) {
	lv{f32.Point{X: x0, Y: y0}, f32.Point{X: x1, Y: y1}}.Stroke(fillcolor, size, &c.Context)
}

// oAbsLine (deprecated) makes a line from (x0,y0) to (x1, y1) using absolute coordinates
// lines are formed with polygons
func (c *Canvas) oAbsLine(x0, y0, x1, y1, lw float32, fillcolor color.RGBA) {
	l2 := lw / 2

	// vertical line
	if x0 == x1 {
		c.AbsRect(x0, y0, lw, y1-y0, fillcolor)
		return
	}
	// horizontal line
	if y0 == y1 {
		c.AbsRect(x0, y0-l2, x1-x0, lw, fillcolor)
		return
	}
	// init for positive sloped line
	p0 := f32.Point{X: x0 - l2, Y: y0 - l2}
	p1 := f32.Point{X: x0 + l2, Y: y0 + l2}
	p2 := f32.Point{X: x1 + l2, Y: y1 + l2}
	p3 := f32.Point{X: x1 - l2, Y: y1 - l2}

	// adjust for negative slope
	if y0 < y1 {
		p0.Y = y0 + l2
		p1.Y = y0 - l2
		p2.Y = y1 - l2
		p3.Y = y1 + l2
	}
	c.quadline(p0, p1, p2, p3, fillcolor)
}

// quadline (deprecated) makes a four-sided polygon with vertices at (p0, p1, p2, p3) forming a line
func (c *Canvas) quadline(p0, p1, p2, p3 f32.Point, fillcolor color.RGBA) {
	path := new(clip.Path)
	ops := c.Context.Ops
	r := f32.Rect(0, 0, c.Width, c.Height)

	defer op.Push(c.Context.Ops).Pop()
	path.Begin(ops)
	path.Move(p0)
	path.Line(f32.Point{X: p1.X - p0.X, Y: p1.Y - p0.Y})
	path.Line(f32.Point{X: p2.X - p1.X, Y: p2.Y - p1.Y})
	path.Line(f32.Point{X: p3.X - p2.X, Y: p3.Y - p2.Y})
	path.Line(f32.Point{X: p0.X - p3.X, Y: p0.Y - p3.Y})
	path.End().Add(ops)
	paint.ColorOp{Color: fillcolor}.Add(ops)
	paint.PaintOp{Rect: r}.Add(ops)
}

// AbsQuadBezier makes a quadratic curve
// starting at (x, y), control point at (cx, cy), end point (ex, ey)
func (c *Canvas) AbsQuadBezier(x, y, cx, cy, ex, ey, size float32, fillcolor color.RGBA) {
	path := new(clip.Path)
	ops := c.Context.Ops
	r := f32.Rect(0, 0, c.Width, c.Height)
	// control and endpoints are relative to the starting point
	ctrl := f32.Point{X: cx - x, Y: cy - y}
	to := f32.Point{X: ex - x, Y: ey - y}

	defer op.Push(c.Context.Ops).Pop()
	path.Begin(ops)
	path.Move(f32.Point{X: x, Y: y})
	path.Quad(ctrl, to)
	path.End().Add(ops)
	paint.ColorOp{Color: fillcolor}.Add(ops)
	paint.PaintOp{Rect: r}.Add(ops)
}

// AbsCubicBezier makes a cubic bezier curve
func (c *Canvas) AbsCubicBezier(x, y, cx1, cy1, cx2, cy2, ex, ey, size float32, fillcolor color.RGBA) {
	path := new(clip.Path)
	ops := c.Context.Ops
	r := f32.Rect(0, 0, c.Width, c.Height)
	// control and end points are relative to the starting point
	sp := f32.Point{X: x, Y: y}
	cp0 := f32.Point{X: cx1 - x, Y: cy1 - y}
	cp1 := f32.Point{X: cx2 - x, Y: cy2 - y}
	ep := f32.Point{X: ex - x, Y: ey - y}

	defer op.Push(c.Context.Ops).Pop()
	path.Begin(ops)
	path.Move(sp)
	path.Cube(cp0, cp1, ep)
	path.End().Add(ops)
	paint.ColorOp{Color: fillcolor}.Add(ops)
	paint.PaintOp{Rect: r}.Add(ops)
}

// AbsCircle makes a circle centered at (x, y), radius r
func (c *Canvas) AbsCircle(x, y, radius float32, fillcolor color.RGBA) {
	path := new(clip.Path)
	ops := c.Context.Ops
	const k = 0.551915024494 // http://spencermortensen.com/articles/bezier-circle/
	r := f32.Rect(0, 0, c.Width, c.Height)

	defer op.Push(c.Context.Ops).Pop()
	path.Begin(ops)
	path.Move(f32.Point{X: x + radius, Y: y})
	path.Cube(f32.Point{X: 0, Y: radius * k}, f32.Point{X: -radius + radius*k, Y: radius}, f32.Point{X: -radius, Y: radius})    // SE
	path.Cube(f32.Point{X: -radius * k, Y: 0}, f32.Point{X: -radius, Y: -radius + radius*k}, f32.Point{X: -radius, Y: -radius}) // SW
	path.Cube(f32.Point{X: 0, Y: -radius * k}, f32.Point{X: radius - radius*k, Y: -radius}, f32.Point{X: radius, Y: -radius})   // NW
	path.Cube(f32.Point{X: radius * k, Y: 0}, f32.Point{X: radius, Y: radius - radius*k}, f32.Point{X: radius, Y: radius})      // NE
	path.End().Add(ops)
	paint.ColorOp{Color: fillcolor}.Add(ops)
	paint.PaintOp{Rect: r}.Add(ops)
}

// AbsEllipse makes a ellipse centered at (x, y) radii (w, h)
func (c *Canvas) AbsEllipse(x, y, w, h float32, fillcolor color.RGBA) {
	path := new(clip.Path)
	ops := c.Context.Ops
	const k = 0.551915024494 // http://spencermortensen.com/articles/bezier-circle/
	r := f32.Rect(0, 0, c.Width, c.Height)
	defer op.Push(c.Context.Ops).Pop()
	path.Begin(ops)
	path.Move(f32.Point{X: x + w, Y: y})
	path.Cube(f32.Point{X: 0, Y: h * k}, f32.Point{X: -w + w*k, Y: h}, f32.Point{X: -w, Y: h})    // SE
	path.Cube(f32.Point{X: -w * k, Y: 0}, f32.Point{X: -w, Y: -h + h*k}, f32.Point{X: -w, Y: -h}) // SW
	path.Cube(f32.Point{X: 0, Y: -h * k}, f32.Point{X: w - w*k, Y: -h}, f32.Point{X: w, Y: -h})   // NW
	path.Cube(f32.Point{X: w * k, Y: 0}, f32.Point{X: w, Y: h - h*k}, f32.Point{X: w, Y: h})      // NE
	path.End().Add(ops)
	paint.ColorOp{Color: fillcolor}.Add(ops)
	paint.PaintOp{Rect: r}.Add(ops)
}

func anglebetweenpoints(p0x, p0y, p1x, p1y float64) float64 {
	return math.Atan2(p1y-p0y, p1x-p0x)
}

func polar(r, theta float64) (float64, float64) {
	return (r * math.Cos(theta)), (r * math.Sin(theta))
}

func radians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func degrees(radians float64) float64 {
	deg := radians * (180 / math.Pi)
	return deg
}

func controls(Ax, Ay, R float64) (float64, float64, float64, float64) {
	Aprimex := ((4 * R) - Ax) / 3
	Aprimey := (R - Ax) * ((3 * R) - Ax) / (3 * Ay)
	return Aprimex, Aprimey, Aprimex, -Aprimey
}

// AbsArc makes an arc centered at (x, y), through angles a1 and a2
func (c *Canvas) AbsArc(x, y, radius, a1, a2 float64, fillcolor color.RGBA) {

	c.AbsCircle(float32(x), float32(y), float32(radius), color.RGBA{0, 0, 0, 100})

	p0x, p0y := polar(radius, a1)
	p1x, p1y := polar(radius, a2)
	theta := anglebetweenpoints(p0x, p0y, p1x, p1y)

	c.AbsTextMid(500, 100, 15, fmt.Sprintf("a1 = %.1f a2 = %.1f theta = %.1f", degrees(a1), degrees(a2), degrees(theta)), color.RGBA{0, 0, 0, 255})
	fmt.Fprintf(os.Stderr, "begin: a1=%.1f, a2=%.1f\n", a1, a2)

	x0 := (radius * math.Cos(theta/2))
	y0 := (radius * math.Sin(theta/2))

	x1, y1, x2, y2 := controls(x0, y0, radius)

	x3 := x0
	y3 := -y0

	c.coord(x, y, 15, "center", color.RGBA{0, 0, 0, 255})
	c.coord(x+p0x, y+p0y, 15, "start", color.RGBA{128, 0, 0, 128})
	c.coord(x+p1x, y+p1y, 15, "end", color.RGBA{0, 128, 0, 128})

	c.coord(x+x0, y+y0, 10, "X0", color.RGBA{0, 0, 0, 255})
	c.coord(x+x1, y+y1, 15, "C1", color.RGBA{0, 0, 0, 255})
	c.coord(x+x2, y+y2, 15, "C2", color.RGBA{0, 0, 0, 255})
	c.coord(x+x3, y+y3, 10, "X3", color.RGBA{0, 0, 0, 255})

	c.AbsPolygon(
		[]float32{float32(x), float32(x + x0), float32(x + x0)},
		[]float32{float32(y), float32(y + y0), float32(y + y3)},
		color.RGBA{128, 0, 0, 100})

	c.AbsCubicBezier(
		float32(x+x0), float32(y+y0),
		float32(x+x1), float32(y+y1),
		float32(x+x2), float32(y+y2),
		float32(x+x3), float32(y+y3),
		float32(radius), color.RGBA{128, 0, 0, 100})

}

func (c *Canvas) coord(x, y, size float64, s string, fillcolor color.RGBA) {
	px := float32(x)
	py := float32(y)
	ls := float32(size)

	c.AbsCircle(px, py, ls, fillcolor)
	c.AbsTextMid(px, py+ls*2, ls, s, fillcolor)
	c.AbsTextMid(px, py-ls, ls, fmt.Sprintf("(%.1f, %.1f)", px, py), fillcolor)
}

// Line implementation from github.com/wrnrlr/shape

type lv []f32.Point

const (
	rad45  = float32(45 * math.Pi / 180)
	rad135 = float32(135 * math.Pi / 180)
	rad315 = float32(315 * math.Pi / 180)
	rad225 = float32(225 * math.Pi / 180)
	rad90  = float32(90 * math.Pi / 180)
	rad180 = float32(180 * math.Pi / 180)
)

func (l lv) Stroke(c color.RGBA, width float32, gtx *layout.Context) (box f32.Rectangle) {
	if len(l) < 2 {
		return box
	}
	defer op.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: c}.Add(gtx.Ops)
	var path clip.Path
	path.Begin(gtx.Ops)
	distance := width
	var angles []float32
	var offsetPoints, originalPoints, deltaPoints []f32.Point
	var tilt float32
	var prevDelta f32.Point
	for i, point := range l {
		if i == 0 {
			nextPoint := l[i+1]
			tilt = angle(point, nextPoint) + rad225
		} else if i == len(l)-1 {
			prevPoint := l[i-1]
			tilt = angle(prevPoint, point) + rad315
		}
		angles = append(angles, tilt)
		originalPoints = append(originalPoints, point)
		point = offsetPoint(point, distance, tilt)
		offsetPoints = append(offsetPoints, point)
		newPoint := point.Sub(prevDelta)
		deltaPoints = append(deltaPoints, newPoint)
		prevDelta = point
		if i == 0 {
			path.Move(newPoint)
		} else {
			path.Line(newPoint)
		}
	}
	for i := len(l) - 1; i >= 0; i-- {
		point := l[i]
		if i == 0 {
			nextPoint := l[i+1]
			tilt = angle(point, nextPoint) + rad135
		} else if i == len(l)-1 {
			prevPoint := l[i-1]
			tilt = angle(prevPoint, point) + rad45
		}
		angles = append(angles, tilt)
		originalPoints = append(originalPoints, point)
		point = offsetPoint(point, distance, tilt)
		offsetPoints = append(offsetPoints, point)
		newPoint := point.Sub(prevDelta)
		deltaPoints = append(deltaPoints, newPoint)
		prevDelta = point
		path.Line(newPoint)
	}
	point := l[0]
	nextPoint := l[1]
	tilt = angle(point, nextPoint) + rad225
	angles = append(angles, tilt)
	originalPoints = append(originalPoints, point)
	point = offsetPoint(point, distance, tilt)
	offsetPoints = append(offsetPoints, point)
	point = point.Sub(prevDelta)
	path.Line(point)
	deltaPoints = append(deltaPoints, point)
	// fmt.Printf("Original Points: %v\n", originalPoints)
	// printDegrees(angles)
	// fmt.Printf("Offset Points:   %v\n", offsetPoints)
	for _, p := range offsetPoints {
		box.Min.X = f32Min(box.Min.X, p.X)
		box.Min.Y = f32Min(box.Min.Y, p.Y)
		box.Max.X = f32Max(box.Max.X, p.X)
		box.Max.Y = f32Max(box.Max.Y, p.Y)
	}
	//fmt.Printf("Min and Max:   %v\n", box)
	//fmt.Printf("Delta Points:    %v\n", deltaPoints)
	path.End().Add(gtx.Ops)
	//paint.PaintOp{f32.Rectangle{Max:f32.Point{w,h}}}.Add(gtx.Ops)
	paint.PaintOp{box}.Add(gtx.Ops)
	return box
}

func angle(p1, p2 f32.Point) float32 {
	return float32(math.Atan2(float64(p2.Y-p1.Y), float64(p2.X-p1.X)))
}

func offsetPoint(point f32.Point, distance, angle float32) f32.Point {
	//fmt.Printf("Point X: %f, Y: %f, Angle: %f\n", point.X, point.Y, angle)
	x := point.X + distance*cos(angle)
	y := point.Y + distance*sin(angle)
	//fmt.Printf("Point X: %f, Y: %f \n", x, y)
	return f32.Point{X: x, Y: y}
}

func cos(v float32) float32 {
	return float32(math.Cos(float64(v)))
}

func sin(v float32) float32 {
	return float32(math.Sin(float64(v)))
}

func atan2(y, x float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

func f32Min(x, y float32) float32 {
	return float32(math.Min(float64(x), float64(y)))
}

func f32Max(x, y float32) float32 {
	return float32(math.Max(float64(x), float64(y)))
}

func printDegrees(radials []float32) {
	var degrees []float32
	for _, a := range radials {
		degrees = append(degrees, f32mod(a*180/math.Pi, 360))
	}
	fmt.Printf("Angles: %v\n", degrees)
}

func f32mod(x, y float32) float32 {
	return float32(math.Mod(float64(x), float64(y)))
}

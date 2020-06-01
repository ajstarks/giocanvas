// Package giocanvas is a 2D canvas API built on gio
package giocanvas

import (
	"image"
	"image/color"
	_ "image/gif" // needed by image
	_ "image/jpeg"
	_ "image/png"
	"os"

	"gioui.org/f32"
	"gioui.org/io/event"
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
func NewCanvas(width, height float32, cfg system.Config, q event.Queue, size image.Point) *Canvas {
	canvas := new(Canvas)
	canvas.Width = width
	canvas.Height = height
	canvas.TextColor = color.RGBA{0, 0, 0, 255}
	size.X = int(width)
	size.Y = int(height)
	canvas.Context = layout.NewContext(new(op.Ops), q, cfg, size)
	return canvas
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

// Line makes a stroked line from (x0, y0) to (x1, y1) using percentage-based measures
func (c *Canvas) Line(x0, y0, x1, y1, size float32, color color.RGBA) {
	x0, y0 = dimen(x0, y0, c.Width, c.Height)
	x1, y1 = dimen(x1, y1, c.Width, c.Height)
	size = pct(size, c.Width)
	c.AbsLine(x0, y0, x1, y1, size, color)
}

// Polygon makes a filled polygon with vertices in x and y, using percentage-based measures
func (c *Canvas) Polygon(x, y []float32, color color.RGBA) {
	if len(x) != len(y) || len(x) < 3 {
		return
	}
	nx := make([]float32, len(x))
	ny := make([]float32, len(y))
	for i := 0; i < len(x); i++ {
		nx[i], ny[i] = dimen(x[i], y[i], c.Width, c.Height)
	}
	c.AbsPolygon(nx, ny, color)
}

// Curve makes a quadric Bezier curve, using percentage-based measures
// starting at (x, y), control point at (cx, cy), end point (ex, ey)
func (c *Canvas) Curve(x, y, cx, cy, ex, ey float32, color color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	cx, cy = dimen(cx, cy, c.Width, c.Height)
	ex, ey = dimen(ex, ey, c.Width, c.Height)
	c.AbsQuadBezier(x, y, cx, cy, ex, ey, 0, color)
}

// CubeCurve makes a cubic Bezier curve, using percentage-based measures
// starting at (x, y), control points at (cx1, cy1), (cx2, cy2), end point (ex, ey)
func (c *Canvas) CubeCurve(x, y, cx1, cy1, cx2, cy2, ex, ey float32, color color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	cx1, cy1 = dimen(cx1, cy1, c.Width, c.Height)
	cx2, cy2 = dimen(cx2, cy2, c.Width, c.Height)
	ex, ey = dimen(ex, ey, c.Width, c.Height)
	c.AbsCubicBezier(x, y, cx1, cy1, cx2, cy2, ex, ey, 0, color)
}

// Circle makes a filled circle, using percentage-based measures
// center is (x,y), radius r
func (c *Canvas) Circle(x, y, r float32, color color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	r = pct(r, c.Width)
	c.AbsCircle(x, y, r, color)
}

// Ellipse makes a filled circle, using percentage-based measures
// center is (x,y), radii (w, h)
func (c *Canvas) Ellipse(x, y, w, h float32, color color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	w = pct(w, c.Width)
	h = pct(h, c.Height)
	c.AbsEllipse(x, y, w, h, color)
}

// Arc makes a filled arc, using percentage-based measures
// center is (x, y) the arc begins at angle a1, and ends at a2
// TODO: placeholder only
func (c *Canvas) Arc(x, y, a1, a2 float32, color color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	a1 = pct(a1, c.Width)
	a2 = pct(a2, c.Width)
	c.AbsArc(x, y, a1, a2, color)
}

// Text places text using percentage-based measures
func (c *Canvas) Text(x, y, size float32, s string, color color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	size = pct(size, c.Width)
	c.textops(x, y, size, text.Start, s, color)
}

// TextEnd places text using percentage-based measures
func (c *Canvas) TextEnd(x, y, size float32, s string, color color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	size = pct(size, c.Width)
	c.textops(x, y, size, text.End, s, color)
}

// TextMid places text using percentage-based measures
func (c *Canvas) TextMid(x, y, size float32, s string, color color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	size = pct(size, c.Width)
	c.textops(x, y, size, text.Middle, s, color)
}

// Rect makes a rectangle with upper left corner at (x,y), with sized at (w,h)
func (c *Canvas) Rect(x, y, w, h float32, color color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	w = pct(w, c.Width)
	h = pct(h, c.Height)
	c.AbsRect(x, y, w, h, color)
}

// Square makes a square shape, centered at (x, y), accounts for screen aspect
func (c *Canvas) Square(x, y, w float32, color color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	w = pct(w, c.Height)
	h := pct(100, w)
	c.AbsCenterRect(x, y, w, h, color)
}

// CenterRect makes a rectangle with upper left corner at (x,y), with sized at (w,h)
func (c *Canvas) CenterRect(x, y, w, h float32, color color.RGBA) {
	x, y = dimen(x, y, c.Width, c.Height)
	w = pct(w, c.Width)
	h = pct(h, c.Height)
	c.AbsCenterRect(x, y, w, h, color)
}

// VLine makes a vertical line beginning at (x,y) with dimension (w, h)
// half of the width is left of x, the other half is the to right of x
func (c *Canvas) VLine(x, y, lineheight, size float32, color color.RGBA) {
	c.Rect(x-(size/2), y+lineheight, size, lineheight, color)
}

// HLine makes a horizontal line starting at (x, y), with dimensions (w, h)
// half of the height is above y, the other below
func (c *Canvas) HLine(x, y, linewidth, size float32, color color.RGBA) {
	c.Rect(x, y+(size/2), linewidth, size, color)
}

// Image places a scaled image centered at (x,y)
func (c *Canvas) Image(name string, x, y float32, w, h int, scale float32) {
	c.CenterImage(name, x, y, w, h, scale)
}

// CenterImage places a scaled image centered at (x,y)
func (c *Canvas) CenterImage(name string, x, y float32, w, h int, scale float32) {
	x, y = dimen(x, y, c.Width, c.Height)
	c.AbsCenterImage(name, x, y, w, h, scale)
}

// Grid makes vertical and horizontal grid lines
func (c *Canvas) Grid(x, y, w, h, size, interval float32, linecolor color.RGBA) {
	for xp := x; xp <= x+w; xp += interval {
		c.Rect(xp, y+h, size, h, linecolor)
	}
	for yp := y; yp <= y+h; yp += interval {
		c.Rect(x, yp, w, size, linecolor)
	}
}

// MapRange maps a value between low1 and high1, return the corresponding value between low2 and high2
func MapRange(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

// textops places text
func (c *Canvas) textops(x, y, size float32, alignment text.Alignment, s string, color color.RGBA) {
	offset := x
	th := material.NewTheme()
	switch alignment {
	case text.End:
		offset = x - c.Width
	case text.Middle:
		offset = x - c.Width/2
	}
	var stack op.StackOp
	stack.Push(c.Context.Ops)
	op.TransformOp{}.Offset(f32.Point{X: offset, Y: y - size}).Add(c.Context.Ops) // shift to use baseline
	l := material.Label(th, unit.Px(size), s)
	l.Color = color
	l.Alignment = alignment
	l.Layout(c.Context)
	stack.Pop()
}

// AbsText places text at (x,y)
func (c *Canvas) AbsText(x, y, size float32, s string, color color.RGBA) {
	c.textops(x, y, size, text.Start, s, color)
}

// AbsTextMid places text centered at (x,y)
func (c *Canvas) AbsTextMid(x, y, size float32, s string, color color.RGBA) {
	c.textops(x, y, size, text.Middle, s, color)
}

// AbsTextEnd places text aligned to the end
func (c *Canvas) AbsTextEnd(x, y, size float32, s string, color color.RGBA) {
	c.textops(x, y, size, text.End, s, color)
}

// AbsRect makes a filled Rectangle; left corner at (x, y), with dimensions (w,h)
func (c *Canvas) AbsRect(x, y, w, h float32, color color.RGBA) {
	ops := c.Context.Ops
	r := f32.Rect(x, y+h, x+w, y)
	paint.ColorOp{Color: color}.Add(ops)
	paint.PaintOp{Rect: r}.Add(ops)
}

// AbsCenterRect makes a filled rectangle centered at (x, y), with dimensions (w,h)
func (c *Canvas) AbsCenterRect(x, y, w, h float32, color color.RGBA) {
	ops := c.Context.Ops
	r := f32.Rect(x-(w/2), y+(h/2), x+(w/2), y-(h/2))
	paint.ColorOp{Color: color}.Add(ops)
	paint.PaintOp{Rect: r}.Add(ops)
}

// AbsVLine makes a vertical line beginning at (x,y) with dimension (w, h)
// half of the width is left of x, the other half is the to right of x
func (c *Canvas) AbsVLine(x, y, w, h float32, color color.RGBA) {
	c.AbsRect(x-(w/2), y, w, h, color)
}

// AbsHLine makes a horizontal line starting at (x, y), with dimensions (w, h)
// half of the height is above y, the other below
func (c *Canvas) AbsHLine(x, y, w, h float32, color color.RGBA) {
	c.AbsRect(x, y-(h/2), w, h, color)
}

// AbsGrid uses horizontal and vertical lines to make a grid
func (c *Canvas) AbsGrid(width, height, size, interval float32, color color.RGBA) {
	var x, y float32
	for y = 0; y <= height; y += height / interval {
		c.AbsHLine(0, y, width, size, color)
	}
	for x = 0; x <= width; x += width / interval {
		c.AbsVLine(x, 0, size, height, color)
	}
}

// AbsCenterImage places a named image centered at (x, y)
// scaled using the specified dimensions (w, h)
func (c *Canvas) AbsCenterImage(name string, x, y float32, w, h int, scale float32) {
	r, err := os.Open(name)
	if err != nil {
		return
	}
	im, _, err := image.Decode(r)
	if err != nil {
		return
	}
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
func (c *Canvas) AbsPolygon(x, y []float32, color color.RGBA) {
	if len(x) != len(y) {
		return
	}
	path := new(clip.Path)
	ops := c.Context.Ops
	r := f32.Rect(0, 0, c.Width, c.Height)
	var stack op.StackOp
	defer stack.Pop()
	stack.Push(c.Context.Ops)

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
	paint.ColorOp{Color: color}.Add(ops)
	paint.PaintOp{Rect: r}.Add(ops)
}

// quadline makes a four-sided polygon with vertices at (p0, p1, p2, p3) forming a line
func (c *Canvas) quadline(p0, p1, p2, p3 f32.Point, color color.RGBA) {
	path := new(clip.Path)
	ops := c.Context.Ops
	r := f32.Rect(0, 0, c.Width, c.Height)
	var stack op.StackOp
	defer stack.Pop()
	stack.Push(c.Context.Ops)
	path.Begin(ops)
	path.Move(p0)
	path.Line(f32.Point{X: p1.X - p0.X, Y: p1.Y - p0.Y})
	path.Line(f32.Point{X: p2.X - p1.X, Y: p2.Y - p1.Y})
	path.Line(f32.Point{X: p3.X - p2.X, Y: p3.Y - p2.Y})
	path.Line(f32.Point{X: p0.X - p3.X, Y: p0.Y - p3.Y})
	path.End().Add(ops)
	paint.ColorOp{Color: color}.Add(ops)
	paint.PaintOp{Rect: r}.Add(ops)
}

// AbsLine makes a line from (x0,y0) to (x1, y1) using absolute coordinates
func (c *Canvas) AbsLine(x0, y0, x1, y1, lw float32, color color.RGBA) {
	l2 := lw / 2

	// vertical line
	if x0 == x1 {
		c.AbsRect(x0, y0, lw, y1-y0, color)
		return
	}
	// horizontal line
	if y0 == y1 {
		c.AbsRect(x0, y0-l2, x1-x0, lw, color)
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
	c.quadline(p0, p1, p2, p3, color)
}

// AbsQuadBezier makes a quadratic curve
// starting at (x, y), control point at (cx, cy), end point (ex, ey)
func (c *Canvas) AbsQuadBezier(x, y, cx, cy, ex, ey, size float32, color color.RGBA) {
	path := new(clip.Path)
	ops := c.Context.Ops
	r := f32.Rect(0, 0, c.Width, c.Height)
	// control and endpoints are relative to the starting point
	ctrl := f32.Point{X: cx - x, Y: cy - y}
	to := f32.Point{X: ex - x, Y: ey - y}
	var stack op.StackOp
	defer stack.Pop()
	stack.Push(c.Context.Ops)
	path.Begin(ops)
	path.Move(f32.Point{X: x, Y: y})
	path.Quad(ctrl, to)
	path.End().Add(ops)
	paint.ColorOp{Color: color}.Add(ops)
	paint.PaintOp{Rect: r}.Add(ops)
}

// AbsCubicBezier makes a cubic bezier curve
func (c *Canvas) AbsCubicBezier(x, y, cx1, cy1, cx2, cy2, ex, ey, size float32, color color.RGBA) {
	path := new(clip.Path)
	ops := c.Context.Ops
	r := f32.Rect(0, 0, c.Width, c.Height)
	// control and end points are relative to the starting point
	sp := f32.Point{X: x, Y: y}
	cp0 := f32.Point{X: cx1 - x, Y: cy1 - y}
	cp1 := f32.Point{X: cx2 - x, Y: cy2 - y}
	ep := f32.Point{X: ex - x, Y: ey - y}
	var stack op.StackOp
	defer stack.Pop()
	stack.Push(c.Context.Ops)
	path.Begin(ops)
	path.Move(sp)
	path.Cube(cp0, cp1, ep)
	path.End().Add(ops)
	paint.ColorOp{Color: color}.Add(ops)
	paint.PaintOp{Rect: r}.Add(ops)
}

// AbsCircle makes a circle centered at (x, y), radius r
func (c *Canvas) AbsCircle(x, y, radius float32, color color.RGBA) {
	path := new(clip.Path)
	ops := c.Context.Ops
	const k = 0.551915024494 // http://spencermortensen.com/articles/bezier-circle/
	r := f32.Rect(0, 0, c.Width, c.Height)
	var stack op.StackOp
	defer stack.Pop()
	stack.Push(c.Context.Ops)
	path.Begin(ops)
	path.Move(f32.Point{X: x + radius, Y: y})
	path.Cube(f32.Point{X: 0, Y: radius * k}, f32.Point{X: -radius + radius*k, Y: radius}, f32.Point{X: -radius, Y: radius})    // SE
	path.Cube(f32.Point{X: -radius * k, Y: 0}, f32.Point{X: -radius, Y: -radius + radius*k}, f32.Point{X: -radius, Y: -radius}) // SW
	path.Cube(f32.Point{X: 0, Y: -radius * k}, f32.Point{X: radius - radius*k, Y: -radius}, f32.Point{X: radius, Y: -radius})   // NW
	path.Cube(f32.Point{X: radius * k, Y: 0}, f32.Point{X: radius, Y: radius - radius*k}, f32.Point{X: radius, Y: radius})      // NE
	path.End().Add(ops)
	paint.ColorOp{Color: color}.Add(ops)
	paint.PaintOp{Rect: r}.Add(ops)
}

// AbsEllipse makes a ellipse centered at (x, y) radii (w, h)
func (c *Canvas) AbsEllipse(x, y, w, h float32, color color.RGBA) {
	path := new(clip.Path)
	ops := c.Context.Ops
	const k = 0.551915024494 // http://spencermortensen.com/articles/bezier-circle/
	r := f32.Rect(0, 0, c.Width, c.Height)
	var stack op.StackOp
	defer stack.Pop()
	stack.Push(c.Context.Ops)
	path.Begin(ops)
	path.Move(f32.Point{X: x + w, Y: y})
	path.Cube(f32.Point{X: 0, Y: h * k}, f32.Point{X: -w + w*k, Y: h}, f32.Point{X: -w, Y: h})    // SE
	path.Cube(f32.Point{X: -w * k, Y: 0}, f32.Point{X: -w, Y: -h + h*k}, f32.Point{X: -w, Y: -h}) // SW
	path.Cube(f32.Point{X: 0, Y: -h * k}, f32.Point{X: w - w*k, Y: -h}, f32.Point{X: w, Y: -h})   // NW
	path.Cube(f32.Point{X: w * k, Y: 0}, f32.Point{X: w, Y: h - h*k}, f32.Point{X: w, Y: h})      // NE
	path.End().Add(ops)
	paint.ColorOp{Color: color}.Add(ops)
	paint.PaintOp{Rect: r}.Add(ops)
}

// AbsArc makes an arc centered at (x, y), through angles a1 and a2
// TODO: placeholder only
func (c *Canvas) AbsArc(x, y, a1, a2 float32, color color.RGBA) {
}

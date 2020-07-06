package giocanvas

import (
	"image"
	"image/color"
	_ "image/gif" // needed by image
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"

	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

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
	op.Offset(f32.Point{X: offset, Y: y - size}).Add(c.Context.Ops) // shift to use baseline
	l := material.Label(material.NewTheme(gofont.Collection()), unit.Px(size), s)
	l.Color = fillcolor
	l.Alignment = alignment
	l.Layout(c.Context)
}

// AbsTextWrap places and wraps text at (x, y), wrapped at width
func (c *Canvas) AbsTextWrap(x, y, size, width float32, s string, fillcolor color.RGBA) {
	defer op.Push(c.Context.Ops).Pop()
	op.Offset(f32.Point{X: x, Y: y - size}).Add(c.Context.Ops) // shift to use baseline
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
	lv{f32.Point{X: x0, Y: y0}, f32.Point{X: x1, Y: y1}}.stroke(fillcolor, size, &c.Context)
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

// AbsArc makes an arc centered at (x, y), through angles a1 and a2
func (c *Canvas) AbsArc(x, y, radius, a1, a2 float64, fillcolor color.RGBA) {

	c.AbsCircle(float32(x), float32(y), float32(radius), color.RGBA{0, 0, 0, 100})

	p0x, p0y := polar(radius, a1)
	p1x, p1y := polar(radius, a2)
	theta := anglebetweenpoints(p0x, p0y, p1x, p1y)

	//fmt.Fprintf(os.Stderr, "begin: a1=%.1f, a2=%.1f\n", a1, a2)

	x0 := (radius * math.Cos(theta/2))
	y0 := (radius * math.Sin(theta/2))

	x1, y1, x2, y2 := controls(x0, y0, radius)

	x3 := x0
	y3 := -y0

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

func (l lv) stroke(c color.RGBA, width float32, gtx *layout.Context) (box f32.Rectangle) {
	if len(l) < 2 {
		return box
	}
	defer op.Push(gtx.Ops).Pop()
	paint.ColorOp{Color: c}.Add(gtx.Ops)
	var path clip.Path
	path.Begin(gtx.Ops)
	distance := width
	var offsetPoints []f32.Point
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
		point = offsetPoint(point, distance, tilt)
		offsetPoints = append(offsetPoints, point)
		newPoint := point.Sub(prevDelta)
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
		point = offsetPoint(point, distance, tilt)
		offsetPoints = append(offsetPoints, point)
		newPoint := point.Sub(prevDelta)
		prevDelta = point
		path.Line(newPoint)
	}
	point := l[0]
	nextPoint := l[1]
	tilt = angle(point, nextPoint) + rad225

	point = offsetPoint(point, distance, tilt)
	offsetPoints = append(offsetPoints, point)
	point = point.Sub(prevDelta)
	path.Line(point)

	for _, p := range offsetPoints {
		box.Min.X = f32Min(box.Min.X, p.X)
		box.Min.Y = f32Min(box.Min.Y, p.Y)
		box.Max.X = f32Max(box.Max.X, p.X)
		box.Max.Y = f32Max(box.Max.Y, p.Y)
	}

	path.End().Add(gtx.Ops)
	paint.PaintOp{Rect: box}.Add(gtx.Ops)
	return box
}

func angle(p1, p2 f32.Point) float32 {
	return float32(math.Atan2(float64(p2.Y-p1.Y), float64(p2.X-p1.X)))
}

func offsetPoint(point f32.Point, distance, angle float32) f32.Point {
	x := point.X + distance*cos(angle)
	y := point.Y + distance*sin(angle)
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

func f32mod(x, y float32) float32 {
	return float32(math.Mod(float64(x), float64(y)))
}

// AbsTranslate moves current location by (x,y)
func (c *Canvas) AbsTranslate(x, y float32) op.StackOp {
	ops := c.Context.Ops
	op.InvalidateOp{}.Add(ops)
	stack := op.Push(ops)
	tr := f32.Affine2D{}
	tr = tr.Offset(f32.Pt(x, y))
	op.Affine(tr).Add(ops)
	return stack
}

// AbsRotate rotates around (x,y) using angle (radians)
func (c *Canvas) AbsRotate(x, y, angle float32) op.StackOp {
	ops := c.Context.Ops
	op.InvalidateOp{}.Add(ops)
	stack := op.Push(ops)
	tr := f32.Affine2D{}.Rotate(f32.Pt(x, y), angle)
	op.Affine(tr).Add(ops)
	return stack
}

// AbsScale scales by factor at (x,y)
func (c *Canvas) AbsScale(x, y, factor float32) op.StackOp {
	ops := c.Context.Ops
	op.InvalidateOp{}.Add(ops)
	stack := op.Push(ops)
	tr := f32.Affine2D{}.Scale(f32.Pt(x, y), f32.Pt(factor, factor))
	op.Affine(tr).Add(ops)
	return stack
}

// AbsShear shears at (x,y) using angle ax and ay
func (c *Canvas) AbsShear(x, y, ax, ay float32) op.StackOp {
	ops := c.Context.Ops
	op.InvalidateOp{}.Add(ops)
	stack := op.Push(ops)
	tr := f32.Affine2D{}.Shear(f32.Pt(x, y), ax, ay)
	op.Affine(tr).Add(ops)
	return stack
}

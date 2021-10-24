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
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

// Foundational methods, and methods using Gio standard coordinates

// textops places text
func (c *Canvas) textops(x, y, size float32, alignment text.Alignment, s string, fillcolor color.NRGBA) {
	offset := x
	switch alignment {
	case text.End:
		offset = x - c.Width
	case text.Middle:
		offset = x - c.Width/2
	}
	defer op.Save(c.Context.Ops).Load()
	op.Offset(f32.Point{X: offset, Y: y - size}).Add(c.Context.Ops) // shift to use baseline
	l := material.Label(material.NewTheme(gofont.Collection()), unit.Px(size), s)
	l.Color = fillcolor
	l.Alignment = alignment
	l.Layout(c.Context)
}

// AbsTextWrap places and wraps text at (x, y), wrapped at width
func (c *Canvas) AbsTextWrap(x, y, size, width float32, s string, fillcolor color.NRGBA) {
	defer op.Save(c.Context.Ops).Load()
	op.Offset(f32.Point{X: x, Y: y - size}).Add(c.Context.Ops) // shift to use baseline
	l := material.Label(material.NewTheme(gofont.Collection()), unit.Px(size), s)
	l.Color = fillcolor
	c.Context.Constraints.Max.X = int(width)
	l.Layout(c.Context)
	c.Context.Constraints.Max.X = int(c.Width) // restore width...
}

// AbsText places text at (x,y)
func (c *Canvas) AbsText(x, y, size float32, s string, fillcolor color.NRGBA) {
	c.textops(x, y, size, text.Start, s, fillcolor)
}

// AbsTextMid places text centered at (x,y)
func (c *Canvas) AbsTextMid(x, y, size float32, s string, fillcolor color.NRGBA) {
	c.textops(x, y, size, text.Middle, s, fillcolor)
}

// AbsTextEnd places text aligned to the end
func (c *Canvas) AbsTextEnd(x, y, size float32, s string, fillcolor color.NRGBA) {
	c.textops(x, y, size, text.End, s, fillcolor)
}

// AbsRect makes a filled Rectangle; left corner at (x, y), with dimensions (w,h)
func (c *Canvas) AbsRect(x, y, w, h float32, fillcolor color.NRGBA) {
	px := make([]float32, 4)
	py := make([]float32, 4)
	px[0], py[0] = x, y
	px[1], py[1] = x+w, y
	px[2], py[2] = x+w, y+h
	px[3], py[3] = x, y+h
	c.AbsPolygon(px, py, fillcolor)
}

// AbsCenterRect makes a filled rectangle centered at (x, y), with dimensions (w,h)
func (c *Canvas) AbsCenterRect(x, y, w, h float32, fillcolor color.NRGBA) {
	c.AbsRect(x-(w/2), y-(h/2), w, h, fillcolor)
}

// AbsVLine makes a vertical line beginning at (x,y) with dimension (w, h)
func (c *Canvas) AbsVLine(x, y, w, h float32, fillcolor color.NRGBA) {
	c.AbsLine(x, y, x, y+h, w, fillcolor)
}

// AbsHLine makes a horizontal line starting at (x, y), with dimensions (w, h)
func (c *Canvas) AbsHLine(x, y, w, h float32, fillcolor color.NRGBA) {
	c.AbsLine(x, y, x+w, y, h, fillcolor)
}

// AbsGrid uses horizontal and vertical lines to make a grid
func (c *Canvas) AbsGrid(width, height, size, interval float32, fillcolor color.NRGBA) {
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
	defer r.Close()
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
	stack := op.Save(ops)
	op.Offset(f32.Pt(x, y)).Add(ops)
	op.Affine(f32.Affine2D{}.Scale(f32.Pt(0, 0), f32.Pt(sc, sc))).Add(ops)
	paint.NewImageOp(im).Add(ops)
	paint.PaintOp{}.Add(ops)
	stack.Load()
}

// AbsPolygon makes a closed, filled polygon with vertices in x and y
func (c *Canvas) AbsPolygon(x, y []float32, fillcolor color.NRGBA) {
	if len(x) != len(y) {
		return
	}
	path := new(clip.Path)
	ops := c.Context.Ops

	defer op.Save(c.Context.Ops).Load()
	path.Begin(ops)
	path.Move(f32.Point{X: x[0], Y: y[0]})

	l := len(x)
	point := f32.Point{}
	for i := 1; i < l; i++ {
		path.Line(f32.Point{X: x[i] - x[i-1], Y: y[i] - y[i-1]})
	}
	path.Line(f32.Point{X: x[0] - x[l-1], Y: y[0] - y[l-1]})
	path.Line(point)
	path.Close()
	clip.Outline{Path: path.End()}.Op().Add(ops)
	paint.ColorOp{Color: fillcolor}.Add(ops)
	paint.PaintOp{}.Add(ops)
}

// AbsLine makes a line from (x0,y0) to (x1, y1) using absolute coordinates
func (c *Canvas) AbsLine(x0, y0, x1, y1, size float32, fillcolor color.NRGBA) {
	path := new(clip.Path)
	ops := c.Context.Ops
	defer op.Save(c.Context.Ops).Load()
	path.Begin(ops)
	path.MoveTo(f32.Point{X: x0, Y: y0})
	path.LineTo(f32.Point{X: x1, Y: y1})
	clip.Stroke{Path: path.End(), Width: size}.Op().Add(ops)
	paint.Fill(ops, fillcolor)
}

// AbsQuadBezier makes a filled quadratic curve
// starting at (x, y), control point at (cx, cy), end point (ex, ey)
func (c *Canvas) AbsQuadBezier(x, y, cx, cy, ex, ey, size float32, fillcolor color.NRGBA) {
	path := new(clip.Path)
	ops := c.Context.Ops
	// control and endpoints are relative to the starting point
	ctrl := f32.Point{X: cx - x, Y: cy - y}
	to := f32.Point{X: ex - x, Y: ey - y}

	defer op.Save(c.Context.Ops).Load()
	path.Begin(ops)
	path.Move(f32.Point{X: x, Y: y})
	path.Quad(ctrl, to)
	path.Close()
	clip.Outline{Path: path.End()}.Op().Add(ops)
	paint.ColorOp{Color: fillcolor}.Add(ops)
	paint.PaintOp{}.Add(ops)
}

// AbsStrokedQuadBezier makes a stroked quadratic curve
// starting at (x, y), control point at (cx, cy), end point (ex, ey)
func (c *Canvas) AbsStrokedQuadBezier(x, y, cx, cy, ex, ey, size float32, strokecolor color.NRGBA) {
	path := new(clip.Path)
	ops := c.Context.Ops
	// control and endpoints are relative to the starting point
	ctrl := f32.Point{X: cx - x, Y: cy - y}
	to := f32.Point{X: ex - x, Y: ey - y}

	defer op.Save(c.Context.Ops).Load()
	path.Begin(ops)
	path.Move(f32.Point{X: x, Y: y})
	path.Quad(ctrl, to)
	clip.Stroke{Path: path.End(), Width: size}.Op().Add(ops)
	paint.Fill(ops, strokecolor)
}

// AbsCubicBezier makes a filled cubic bezier curve
func (c *Canvas) AbsCubicBezier(x, y, cx1, cy1, cx2, cy2, ex, ey, size float32, fillcolor color.NRGBA) {
	path := new(clip.Path)
	ops := c.Context.Ops
	// control and end points are relative to the starting point
	sp := f32.Point{X: x, Y: y}
	cp0 := f32.Point{X: cx1 - x, Y: cy1 - y}
	cp1 := f32.Point{X: cx2 - x, Y: cy2 - y}
	ep := f32.Point{X: ex - x, Y: ey - y}

	defer op.Save(c.Context.Ops).Load()
	path.Begin(ops)
	path.Move(sp)
	path.Cube(cp0, cp1, ep)
	path.Close()
	clip.Outline{Path: path.End()}.Op().Add(ops)
	paint.ColorOp{Color: fillcolor}.Add(ops)
	paint.PaintOp{}.Add(ops)
}

// AbsStrokedCubicBezier makes a stroked cubic bezier curve
func (c *Canvas) AbsStrokedCubicBezier(x, y, cx1, cy1, cx2, cy2, ex, ey, size float32, strokecolor color.NRGBA) {
	path := new(clip.Path)
	ops := c.Context.Ops
	// control and end points are relative to the starting point
	sp := f32.Point{X: x, Y: y}
	cp0 := f32.Point{X: cx1 - x, Y: cy1 - y}
	cp1 := f32.Point{X: cx2 - x, Y: cy2 - y}
	ep := f32.Point{X: ex - x, Y: ey - y}

	defer op.Save(c.Context.Ops).Load()
	path.Begin(ops)
	path.Move(sp)
	path.Cube(cp0, cp1, ep)
	clip.Stroke{Path: path.End(), Width: size}.Op().Add(ops)
	paint.Fill(ops, strokecolor)
}

// AbsCircle makes a circle centered at (x, y), radius r
func (c *Canvas) AbsCircle(x, y, radius float32, fillcolor color.NRGBA) {
	path := new(clip.Path)
	ops := c.Context.Ops
	const k = 0.551915024494 // http://spencermortensen.com/articles/bezier-circle/

	defer op.Save(c.Context.Ops).Load()
	path.Begin(ops)
	path.Move(f32.Point{X: x + radius, Y: y})
	path.Cube(f32.Point{X: 0, Y: radius * k}, f32.Point{X: -radius + radius*k, Y: radius}, f32.Point{X: -radius, Y: radius})    // SE
	path.Cube(f32.Point{X: -radius * k, Y: 0}, f32.Point{X: -radius, Y: -radius + radius*k}, f32.Point{X: -radius, Y: -radius}) // SW
	path.Cube(f32.Point{X: 0, Y: -radius * k}, f32.Point{X: radius - radius*k, Y: -radius}, f32.Point{X: radius, Y: -radius})   // NW
	path.Cube(f32.Point{X: radius * k, Y: 0}, f32.Point{X: radius, Y: radius - radius*k}, f32.Point{X: radius, Y: radius})      // NE
	path.Close()
	clip.Outline{Path: path.End()}.Op().Add(ops)
	paint.ColorOp{Color: fillcolor}.Add(ops)
	paint.PaintOp{}.Add(ops)
}

// AbsEllipse makes a ellipse centered at (x, y) radii (w, h)
func (c *Canvas) AbsEllipse(x, y, w, h float32, fillcolor color.NRGBA) {
	path := new(clip.Path)
	ops := c.Context.Ops
	const k = 0.551915024494 // http://spencermortensen.com/articles/bezier-circle/
	defer op.Save(c.Context.Ops).Load()
	path.Begin(ops)
	path.Move(f32.Point{X: x + w, Y: y})
	path.Cube(f32.Point{X: 0, Y: h * k}, f32.Point{X: -w + w*k, Y: h}, f32.Point{X: -w, Y: h})    // SE
	path.Cube(f32.Point{X: -w * k, Y: 0}, f32.Point{X: -w, Y: -h + h*k}, f32.Point{X: -w, Y: -h}) // SW
	path.Cube(f32.Point{X: 0, Y: -h * k}, f32.Point{X: w - w*k, Y: -h}, f32.Point{X: w, Y: -h})   // NW
	path.Cube(f32.Point{X: w * k, Y: 0}, f32.Point{X: w, Y: h - h*k}, f32.Point{X: w, Y: h})      // NE
	clip.Outline{Path: path.End()}.Op().Add(ops)
	paint.ColorOp{Color: fillcolor}.Add(ops)
	paint.PaintOp{}.Add(ops)
}

// AbsArc makes circular arc centered at (x, y), through angles start and end;
// the angles are measured in radians and increase counter-clockwise.
// N.B: derived from the clipLoader function in widget/material/loader.go
func (c *Canvas) AbsArc(x, y, radius float32, start, end float64, fillcolor color.NRGBA) {
	ops := c.Context.Ops
	sine, cose := math.Sincos(start)
	defer op.Save(ops).Load()
	path := new(clip.Path)
	path.Begin(ops)
	path.Move(f32.Pt(x, y))                                 // move to the center
	pen := f32.Pt(float32(cose), float32(sine)).Mul(radius) // starting point
	path.Line(pen)

	// The clip path uses quadratic beziér curves to approximate
	// a circle arc. Minimize the error by capping the length of
	// each curve segment.
	const maxArcLen = 20.0
	arcPerRadian := float64(radius) * math.Pi
	anglePerSegment := maxArcLen / arcPerRadian
	for angle := start; angle < end; {
		angle += anglePerSegment
		if angle > end {
			angle = end
		}
		sins, coss := sine, cose
		sine, cose = math.Sincos(angle)

		// https://pomax.github.io/bezierinfo/#circles
		div := 1.0 / (coss*sine - cose*sins)
		ctrlPt := f32.Point{X: float32((sine - sins) * div), Y: -float32((cose - coss) * div)}.Mul(radius)
		endPt := f32.Pt(float32(cose), float32(sine)).Mul(radius)
		path.Quad(ctrlPt.Sub(pen), endPt.Sub(pen))
		pen = endPt
	}
	path.Close()
	clip.Outline{Path: path.End()}.Op().Add(ops)
	paint.ColorOp{Color: fillcolor}.Add(ops)
	paint.PaintOp{}.Add(ops)
}

// AbsTranslate moves current location by (x,y)
func (c *Canvas) AbsTranslate(x, y float32) op.TransformStack {
	ops := c.Context.Ops
	op.InvalidateOp{}.Add(ops)
	stack := op.Offset(f32.Pt(0, 0)).Push(ops)
	tr := f32.Affine2D{}
	tr = tr.Offset(f32.Pt(x, y))
	op.Affine(tr).Add(ops)
	return stack
}

// AbsRotate rotates around (x,y) using angle (radians)
func (c *Canvas) AbsRotate(x, y, angle float32) op.TransformStack {
	ops := c.Context.Ops
	op.InvalidateOp{}.Add(ops)
	stack := op.Offset(f32.Pt(0, 0)).Push(ops)
	tr := f32.Affine2D{}.Rotate(f32.Pt(x, y), angle)
	op.Affine(tr).Add(ops)
	return stack
}

// AbsScale scales by factor at (x,y)
func (c *Canvas) AbsScale(x, y, factor float32) op.TransformStack {
	ops := c.Context.Ops
	op.InvalidateOp{}.Add(ops)
	stack := op.Offset(f32.Pt(0, 0)).Push(ops)
	tr := f32.Affine2D{}.Scale(f32.Pt(x, y), f32.Pt(factor, factor))
	op.Affine(tr).Add(ops)
	return stack
}

// AbsShear shears at (x,y) using angle ax and ay
func (c *Canvas) AbsShear(x, y, ax, ay float32) op.TransformStack {
	ops := c.Context.Ops
	op.InvalidateOp{}.Add(ops)
	stack := op.Offset(f32.Pt(0, 0)).Push(ops)
	tr := f32.Affine2D{}.Shear(f32.Pt(x, y), ax, ay)
	op.Affine(tr).Add(ops)
	return stack
}

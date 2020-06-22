package main

import (
	"flag"
	"image/color"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func translate(ops *op.Ops, x, y float32) op.StackOp {
	op.InvalidateOp{}.Add(ops)
	stack := op.Push(ops)
	tr := f32.Affine2D{}
	tr = tr.Offset(f32.Pt(x, y))
	op.Affine(tr).Add(ops)
	return stack
}

func rotate(ops *op.Ops, x, y, angle float32) op.StackOp {
	op.InvalidateOp{}.Add(ops)
	stack := op.Push(ops)
	tr := f32.Affine2D{}.Rotate(f32.Pt(x, y), angle)
	op.Affine(tr).Add(ops)
	return stack
}

func scale(ops *op.Ops, x, y, factor float32) op.StackOp {
	op.InvalidateOp{}.Add(ops)
	stack := op.Push(ops)
	tr := f32.Affine2D{}.Scale(f32.Pt(x, y), f32.Pt(factor, factor))
	op.Affine(tr).Add(ops)
	return stack
}

func shear(ops *op.Ops, x, y, ax, ay float32) op.StackOp {
	op.InvalidateOp{}.Add(ops)
	stack := op.Push(ops)
	tr := f32.Affine2D{}.Shear(f32.Pt(x, y), ax, ay)
	op.Affine(tr).Add(ops)
	return stack
}

func endtransform(stack op.StackOp) {
	stack.Pop()
}

func transforms(title string, w, h int) {
	width, height := float32(w), float32(h)
	win := app.NewWindow(app.Title(title), app.Size(unit.Px(width), unit.Px(height)))

	midx := width / 2
	rectw := width * 0.4
	recth := rectw / 4
	ts := width * 0.05
	black := color.RGBA{0, 0, 0, 255}
	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, e)
			ctx := canvas.Context.Ops

			canvas.AbsCenterRect(midx, 100, rectw, recth, color.RGBA{128, 128, 128, 128})
			canvas.AbsTextMid(midx, 110, ts, "Reference", black)

			stack := shear(ctx, midx, midx, math.Pi/4, 0)
			canvas.AbsCenterRect(midx, midx, rectw, recth, color.RGBA{128, 0, 0, 128})
			canvas.AbsTextMid(midx, 510, ts, "shear", black)
			endtransform(stack)

			stack = translate(ctx, 200, 150)
			canvas.AbsCenterRect(midx, midx, rectw, recth, color.RGBA{0, 128, 0, 128})
			canvas.AbsTextMid(midx, 510, ts, "translate", black)
			endtransform(stack)

			stack = scale(ctx, midx, 300, 2)
			canvas.AbsCenterRect(midx, 300, rectw, recth, color.RGBA{0, 0, 128, 128})
			canvas.AbsTextMid(midx, 310, ts, "scale", black)
			endtransform(stack)

			stack = rotate(ctx, midx, 800, math.Pi/4)
			canvas.AbsCenterRect(midx, 800, rectw, recth, color.RGBA{255, 50, 0, 200})
			canvas.AbsTextMid(midx, 810, ts, "rotate", black)
			endtransform(stack)

			e.Frame(ctx)
		case key.Event:
			switch e.Name {
			case "Q", key.NameEscape:
				os.Exit(0)
			}
		}
	}
}

func main() {
	var w, h int
	flag.IntVar(&w, "width", 1000, "canvas width")
	flag.IntVar(&h, "height", 1000, "canvas height")
	flag.Parse()
	go transforms("transforms", w, h)
	app.Main()
}

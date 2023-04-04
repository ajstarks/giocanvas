// bezsketch
package main

import (
	"flag"
	"fmt"
	"image/color"
	"io"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

func main() {
	var cw, ch int
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.Parse()

	width := float32(cw)
	height := float32(ch)

	go func() {
		w := app.NewWindow(app.Title("bezsketch"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := bezsketch(w, width, height); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

func curve(ops *op.Ops, x, y, cx, cy, ex, ey, size float32, strokecolor color.NRGBA) {
	path := new(clip.Path)
	// control and endpoints are relative to the starting point
	ctrl := f32.Point{X: cx - x, Y: cy - y}
	to := f32.Point{X: ex - x, Y: ey - y}

	path.Begin(ops)
	path.Move(f32.Point{X: x, Y: y})
	path.Quad(ctrl, to)
	stack := clip.Stroke{Path: path.End(), Width: size}.Op().Push(ops)
	paint.Fill(ops, strokecolor)
	stack.Pop()

}

func circle(ops *op.Ops, x, y, radius float32, fillcolor color.NRGBA) {
	const k = 0.551915024494 // http://spencermortensen.com/articles/bezier-circle/
	path := new(clip.Path)
	path.Begin(ops)
	path.Move(f32.Point{X: x + radius, Y: y})
	path.Cube(f32.Point{X: 0, Y: radius * k}, f32.Point{X: -radius + radius*k, Y: radius}, f32.Point{X: -radius, Y: radius})    // SE
	path.Cube(f32.Point{X: -radius * k, Y: 0}, f32.Point{X: -radius, Y: -radius + radius*k}, f32.Point{X: -radius, Y: -radius}) // SW
	path.Cube(f32.Point{X: 0, Y: -radius * k}, f32.Point{X: radius - radius*k, Y: -radius}, f32.Point{X: radius, Y: -radius})   // NW
	path.Cube(f32.Point{X: radius * k, Y: 0}, f32.Point{X: radius, Y: radius - radius*k}, f32.Point{X: radius, Y: radius})      // NE
	path.Close()
	stack := clip.Outline{Path: path.End()}.Op().Push(ops)
	paint.ColorOp{Color: fillcolor}.Add(ops)
	paint.PaintOp{}.Add(ops)
	stack.Pop()
}

var pressed bool
var mouseX, mouseY float32
var bx, by, cx, cy, ex, ey float32

func mousePos(q event.Queue) {
	for _, ev := range q.Events(pressed) {
		if p, ok := ev.(pointer.Event); ok {
			switch p.Type {
			case pointer.Drag:
				mouseX = p.Position.X
				mouseY = p.Position.Y
			case pointer.Press:
				switch p.Buttons {
				case pointer.ButtonPrimary:
					bx = p.Position.X
					by = p.Position.Y
				case pointer.ButtonSecondary:
					ex = p.Position.X
					ey = p.Position.Y
				}
				pressed = true
			}
		}
	}
}

func bezsketch(w *app.Window, width, height float32) error {
	beginColor := color.NRGBA{R: 0, G: 0, B: 255, A: 255}
	endColor := color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	bx, by = 250.0, 500.0
	cx, cy = 100.0, 100.0
	ex, ey = 750.0, 500.0
	curveColor := color.NRGBA{R: 128, G: 128, B: 128, A: 128}
	var ops op.Ops
	for {
		ev := <-w.Events()
		switch e := ev.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			ops.Reset()
			pointer.InputOp{Tag: pressed, Grab: true, Types: pointer.Press | pointer.Drag}.Add(&ops)
			mousePos(e.Queue)
			cx, cy = mouseX, mouseY
			fmt.Printf("curve %v %v %v %v %v %v\n", bx, by, cx, cy, ex, ey)
			circle(&ops, bx, by, 10, beginColor)
			circle(&ops, ex, ey, 10, endColor)
			curve(&ops, bx, by, cx, cy, ex, ey, 10, curveColor)
			e.Frame(&ops)
		}
	}
}

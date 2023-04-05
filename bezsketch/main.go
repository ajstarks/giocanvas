// bezsketch: sketch quadradic Bezier curves
package main

import (
	"flag"
	"image/color"
	"io"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
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

var pressed bool
var mouseX, mouseY float32
var bx, by, ex, ey, cx, cy float32
var black = color.NRGBA{R: 0, G: 0, B: 0, A: 255}

// pctcoord converts device coordinates to canvas percents
func pctcoord(x, y, width, height float32) (float32, float32) {
	return 100 * (x / width), 100 - (100 * (y / height))
}

// ftoa converts float to string, with leading space
func ftoa(x float32) string {
	return " " + strconv.FormatFloat(float64(x), 'f', 1, 32)
}

// textcoord displays a labelled coordinate
func textcoord(canvas *giocanvas.Canvas, x, y, size float32, color color.NRGBA) {
	canvas.Circle(x, y, size/2, color)
	canvas.TextMid(x, y+size, size, ftoa(x)+","+ftoa(y), black)
}

// pctPointerPos records and processes the pointer events in percent coordinates
func pctPointerPos(q event.Queue, width, height float32) {
	for _, ev := range q.Events(pressed) {
		if p, ok := ev.(pointer.Event); ok {
			switch p.Type {
			case pointer.Move:
				mouseX, mouseY = pctcoord(p.Position.X, p.Position.Y, width, height)
			case pointer.Press:
				switch p.Buttons {
				case pointer.ButtonPrimary:
					bx, by = pctcoord(p.Position.X, p.Position.Y, width, height)
				case pointer.ButtonSecondary:
					ex, ey = pctcoord(p.Position.X, p.Position.Y, width, height)
				case pointer.ButtonTertiary:
					io.WriteString(os.Stdout, "curve"+ftoa(bx)+ftoa(by)+ftoa(cx)+ftoa(cy)+ftoa(ex)+ftoa(ey)+"\n")
				}
				pressed = true
			}
		}
	}
}

// bezsketch sketches quadratic bezier curves: left pointer press defines the begin point,
// right pointer press defines the end point, middle pointer press shows the curve spec,
// pointer move defines the control point
func bezsketch(w *app.Window, width, height float32) error {
	beginColor := color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	endColor := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	curveColor := color.NRGBA{R: 128, G: 128, B: 128, A: 128}

	bx, by = 25.0, 50.0
	ex, ey = 75.0, 50.0
	cx, cy = 10, 10
	var ts float32 = 2.5

	// event loop
	for {
		ev := <-w.Events()
		switch e := ev.(type) {
		// return an error on close
		case system.DestroyEvent:
			return e.Err

		// for each frame:
		// register press and move events, draw coordinates, and the curve,
		// track the pointer position for the control point, show curve spec on middle click
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, system.FrameEvent{})
			pointer.InputOp{Tag: pressed, Grab: false, Types: pointer.Press | pointer.Move}.Add(canvas.Context.Ops)
			textcoord(canvas, bx, by, ts, beginColor)
			textcoord(canvas, ex, ey, ts, endColor)
			textcoord(canvas, cx, cy, ts, curveColor)
			canvas.QuadStrokedCurve(bx, by, cx, cy, ex, ey, 0.75, curveColor)
			pctPointerPos(e.Queue, width, height)
			cx, cy = mouseX, mouseY
			e.Frame(canvas.Context.Ops)
		}
	}
}

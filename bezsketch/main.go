// bezsketch
package main

import (
	"flag"
	"fmt"
	"image/color"
	"io"
	"os"

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
var bx, by, ex, ey float32

func pctcoord(x, y, width, height float32) (float32, float32) {
	return 100 * (x / width), 100 - (100 * (y / height))
}

func textcoord(canvas *giocanvas.Canvas, x, y, size float32, color color.NRGBA) {
	canvas.Circle(x, y, size/2, color)
	canvas.TextMid(x, y+size, size, fmt.Sprintf("(%.1f, %.1f)", x, y), giocanvas.ColorLookup("black"))
}

func pctmousePos(q event.Queue, width, height float32) {
	for _, ev := range q.Events(pressed) {
		if p, ok := ev.(pointer.Event); ok {
			switch p.Type {
			case pointer.Drag:
				mouseX, mouseY = pctcoord(p.Position.X, p.Position.Y, width, height)
			case pointer.Press:
				switch p.Buttons {
				case pointer.ButtonPrimary:
					bx, by = pctcoord(p.Position.X, p.Position.Y, width, height)
				case pointer.ButtonSecondary:
					ex, ey = pctcoord(p.Position.X, p.Position.Y, width, height)
				}
				pressed = true
			}
		}
	}
}

func bezsketch(w *app.Window, width, height float32) error {
	beginColor := color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	endColor := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	curveColor := color.NRGBA{R: 128, G: 128, B: 128, A: 128}

	bx, by = 25.0, 50.0
	ex, ey = 75.0, 50.0
	var cx, cy float32 = 10, 10
	for {
		ev := <-w.Events()
		switch e := ev.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, system.FrameEvent{})
			pointer.InputOp{Tag: pressed, Grab: false, Types: pointer.Press | pointer.Drag}.Add(canvas.Context.Ops)
			textcoord(canvas, bx, by, 2, beginColor)
			textcoord(canvas, ex, ey, 2, endColor)
			canvas.QuadStrokedCurve(bx, by, cx, cy, ex, ey, 1, curveColor)
			fmt.Printf("curve %.1f %.1f %.1f %.1f %.1f %.1f\n", bx, by, cx, cy, ex, ey)
			pctmousePos(e.Queue, width, height)
			cx, cy = mouseX, mouseY
			e.Frame(canvas.Context.Ops)
		}
	}
}

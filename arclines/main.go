package main

import (
	"image/color"
	"math"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/op"

	"github.com/ajstarks/giocanvas"
)

const (
	// basic positions in radians
	right  = 0
	top    = math.Pi / 2
	left   = math.Pi
	bottom = math.Pi * 1.5
	twoPi  = math.Pi * 2

	// arc coordinates
	x1  = 25.0
	mid = 50.0
	x3  = 100 - x1

	opFramerate = time.Second / 50.0 // 40ms
)

func main() {
	go draw()
	app.Main()
}

func draw() {
	w := &app.Window{}
	w.Option(app.Title("hello"))
	c1 := color.NRGBA{128, 0, 0, 255}
	c2 := color.NRGBA{0, 0, 128, 255}
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			// Calculate delta time
			currentStep := 0.1

			// canvas
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})
			// context must be set before drawing to avoid bug?
			gtx := app.NewContext(canvas.Context.Ops, e)

			// Draw the arcs
			drawArcAntiClockwise(canvas, currentStep, c1)
			drawArcClockwise(canvas, currentStep, c2)

			// Redraw the canvas after opFramerate duration
			inv := op.InvalidateCmd{At: gtx.Now.Add(opFramerate)}
			gtx.Execute(inv)

			e.Frame(canvas.Context.Ops)
		}
	}
}

var (
	angle1 = twoPi
	angle3 = twoPi
)

func drawArcAntiClockwise(canvas *giocanvas.Canvas, currentStep float64, color color.NRGBA) {
	angle1 += currentStep
	// avoid angle1 getting too big for a float64 on long running animation
	angle1 = math.Mod(angle1, twoPi)

	base := top
	a1 := base
	a2 := base + angle1
	canvas.ArcLine(x1, mid, 5, a1, a2, 0.5, color)
	canvas.ArcLine(x1, mid, 10, a1, a2, 0.5, color)
	canvas.ArcLine(x1, mid, 20, a1, a2, 0.5, color)
}

func drawArcClockwise(canvas *giocanvas.Canvas, currentStep float64, color color.NRGBA) {
	angle3 += currentStep
	// avoid angle3 getting too big for a float64 on long running animation
	angle3 = math.Mod(angle3, twoPi)

	base := top
	a1 := base - angle3
	a2 := base
	canvas.ArcLine(x3, mid, 5, a1, a2, 0.5, color)
	canvas.ArcLine(x3, mid, 10, a1, a2, 0.5, color)
	canvas.ArcLine(x3, mid, 20, a1, a2, 0.5, color)
}

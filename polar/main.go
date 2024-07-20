// polar demos polar coordinates
package main

import (
	"flag"
	"image/color"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func polar(title string, width, height int) {
	// init app window
	w := &app.Window{}
	w.Option(app.Title(title), app.Size(unit.Dp(float32(width)), unit.Dp(float32(height))))
	// set colors
	topcolor := color.NRGBA{255, 0, 0, 100}
	botcolor := color.NRGBA{0, 0, 255, 100}
	bgcolor := color.NRGBA{0, 0, 0, 255}
	// app loop
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})
			canvas.Background(bgcolor)
			var theta, radius float32
			for radius = 2; radius < 50; radius += 2 {
				for theta = 180; theta <= 360; theta += 15 { // degrees
					x, y := canvas.PolarDegrees(50, 50, radius, theta)
					canvas.Circle(x, y, radius/12, topcolor)
				}
				for theta = math.Pi / 16; theta < math.Pi; theta += math.Pi / 16 { // radians
					x, y := canvas.Polar(50, 50, radius, theta)
					canvas.Circle(x, y, radius/12, botcolor)
				}
			}
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {
	cw := flag.Int("width", 1000, "canvas width")
	ch := flag.Int("height", 1000, "canvas height")
	flag.Parse()
	go polar("polar coordinates", *cw, *ch)
	app.Main()
}

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

func polar(title string, width, height float32) {
	w := &app.Window{}
	w.Option(app.Title(title), app.Size(unit.Dp(width), unit.Dp(height)))
	topcolor := color.NRGBA{255, 0, 0, 100}
	botcolor := color.NRGBA{0, 0, 255, 100}
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})
			canvas.Background(color.NRGBA{0, 0, 0, 255})
			var theta, radius float32
			for radius = 2; radius < 50; radius += 2 {
				for theta = 180; theta <= 360; theta += 15 {
					x, y := canvas.PolarDegrees(50, 50, radius, theta)
					canvas.Circle(x, y, radius/12, topcolor)
				}
				for theta = math.Pi / 16; theta < math.Pi; theta += math.Pi / 16 {
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
	go polar("polar coordinates", float32(*cw), float32(*ch))
	app.Main()
}

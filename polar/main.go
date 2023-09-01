// polar demos polar coordinates
package main

import (
	"flag"
	"image/color"
	"io"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func circles(w *app.Window, width, height float32) error {
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), system.FrameEvent{})
			canvas.Background(color.NRGBA{0, 0, 0, 255})
			var theta, radius float32
			for radius = 5; radius < 50; radius += 5 {
				for theta = 180; theta <= 360; theta += 15 {
					x, y := canvas.PolarDegrees(50, 50, radius, theta)
					canvas.Circle(x, y, radius/12, color.NRGBA{128, 0, 0, 120})
				}
				for theta = math.Pi / 16; theta < math.Pi; theta += math.Pi / 16 {
					x, y := canvas.Polar(50, 50, radius, theta)
					canvas.Circle(x, y, radius/12, color.NRGBA{0, 0, 128, 120})
				}
			}
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {
	var cw, ch int
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.Parse()
	width := float32(cw)
	height := float32(ch)

	go func() {
		w := app.NewWindow(app.Title("polar"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := circles(w, width, height); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()

}

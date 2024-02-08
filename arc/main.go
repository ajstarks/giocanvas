// test arcs
package main

import (
	"flag"
	"fmt"
	"image/color"
	"io"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

const pi = 3.14159265358979323846264338327950288419716939937510582097494459

func radians(degrees float64) float64 {
	return degrees * (pi / 180)
}

func degrees(radians float64) float64 {
	return radians * (180 / pi)
}

func arc(w *app.Window, width, height float32) error {

	var x, y, step float32
	size := width / 10
	step = width * 0.11
	ts := width * .01800
	radius := width / 20
	dotsize := width * 0.005
	bgcolor := color.NRGBA{0, 0, 0, 50}
	arcolor := color.NRGBA{0, 0, 128, 100}
	centercolor := color.NRGBA{128, 128, 128, 128}
	black := color.NRGBA{0, 0, 0, 255}

	for {
		e := w.NextEvent()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})
			angle := 45.0 // Pi / 4
			y = height * 0.2
			canvas.CText(50, 50, 10, "Arcs", black)
			for x = size; x < width-size; x += step {
				canvas.AbsCircle(x, y, radius, bgcolor)
				canvas.AbsCircle(x, y, dotsize, centercolor)
				canvas.AbsArc(x, y, radius, 0, radians(angle), arcolor)
				canvas.AbsTextMid(x, y+(width*0.10), ts, fmt.Sprintf("%.1f°", angle), black)
				canvas.AbsTextMid(x, y-(width*0.08), ts, fmt.Sprintf("%.4f rad", radians(angle)), black)
				angle += 45 // Pi / 4
			}

			y = 20
			angle = pi * 2
			for x = 10; x < 90; x += 11 {
				canvas.Circle(x, y, 5, bgcolor)
				canvas.Circle(x, y, 0.5, centercolor)
				canvas.Arc(x, y, 5, 0, angle, arcolor)
				canvas.ArcLine(x, y, 5, 0, angle, 0.2, color.NRGBA{128, 0, 0, 255})
				canvas.CText(x, y-10, 1.8, fmt.Sprintf("%.1f°", degrees(angle)), black)
				canvas.CText(x, y+8, 1.8, fmt.Sprintf("%.4f rad", angle), black)
				angle -= pi / 4
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
		w := app.NewWindow(app.Title("Arc"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := arc(w, width, height); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()

}

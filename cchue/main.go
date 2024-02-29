// cchue - concentric circles
package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

// hsv makes a color given hue, saturation, value
func hsv(hue, sat, value int) color.NRGBA {
	return giocanvas.ColorLookup(fmt.Sprintf("hsv(%d,%d,%d)", hue, sat, value))
}

// planet makes circles around a central point
func planet(canvas *giocanvas.Canvas, x, y, size, radius, a1, a2, steps float32, color color.NRGBA) {
	for t := a1; t < a2; t += steps {
		px, py := canvas.PolarDegrees(x, y, radius, t)
		canvas.Circle(px, py, size, color)
	}
}

// cchue makes a concentric circle pattern with varying hue
func cchue(canvas *giocanvas.Canvas, r, step float32, starthue int, bgcolor color.NRGBA) {
	var cstep, c, halfstep, csize float32
	cstep = 0.5
	c = 0.5
	halfstep = step / 2
	csize = r * 0.75
	hue := starthue

	canvas.Background(bgcolor)
	canvas.Circle(50, 50, csize, hsv(hue, 100, 100))
	planet(canvas, 50, 50, c, r, 0, 360, step, hsv(hue, 100, 100))
	r += 2
	c += cstep
	hue += 7
	planet(canvas, 50, 50, c, r, halfstep, 360, step, hsv(hue, 100, 100))
	r += 3
	c += cstep
	hue += 7
	planet(canvas, 50, 50, c, r, 0, 360, step, hsv(hue, 100, 100))
	r += 4
	c += cstep
	hue += 7
	planet(canvas, 50, 50, c, r, halfstep, 360, step, hsv(hue, 100, 100))
	r += 5
	c += cstep
	hue += 7
	planet(canvas, 50, 50, c, r, 0, 360, step, hsv(hue, 100, 100))
	r += 6
	c += cstep
	hue += 7
	planet(canvas, 50, 50, c, r, halfstep, 360, step, hsv(hue, 100, 100))
	r += 8
	hue += 7
	planet(canvas, 50, 50, 3.5, r, 0, 360, step, hsv(hue, 100, 100))

	var t float32
	for t = 0.0; t <= 360; t += step {
		px, py := canvas.PolarDegrees(50, 50, r, t)
		planet(canvas, px, py, 0.5, 5, 0, 360, 30, hsv(starthue, 100, 100))
	}
}

func cc(w *app.Window) error {
	for {
		e := w.NextEvent()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})
			cchue(canvas, 10, 20, 0, color.NRGBA{0, 0, 0, 255})
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {
	var cw, ch int
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.Parse()

	width, height := float32(cw), float32(ch)

	go func() {
		w := app.NewWindow(app.Title("cchue"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := cc(w); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create the window: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

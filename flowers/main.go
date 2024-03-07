// flowers demonstrated transforms with "flowers"
package main

import (
	"flag"
	"image/color"
	"io"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func petals(canvas *giocanvas.Canvas, x, y, w, h float32, fill color.NRGBA) {
	var a float32
	for a = 0; a < math.Pi*2; a += math.Pi / 9 {
		stack := canvas.Rotate(x, y, a)
		canvas.Ellipse(x, y, w, h, fill)
		giocanvas.EndTransform(stack)
	}
}

func flower(w *app.Window) error {
	red := color.NRGBA{128, 0, 0, 100}
	blue := color.NRGBA{0, 0, 128, 100}
	green := color.NRGBA{0, 128, 0, 100}
	orange := giocanvas.ColorLookup("orange")
	bgcolor := giocanvas.ColorLookup("linen")

	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})
			canvas.Background(bgcolor)
			petals(canvas, 10, 90, 5, 1, red)
			petals(canvas, 25, 75, 10, 1.5, green)
			petals(canvas, 50, 50, 15, 3.0, blue)
			petals(canvas, 80, 20, 20, 4.5, orange)
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
		w := &app.Window{}
		w.Option(app.Title("flowers"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := flower(w); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

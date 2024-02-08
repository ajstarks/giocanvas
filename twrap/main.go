// text wrapping
package main

import (
	"flag"
	"image/color"
	"io"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

// pct returns the percentage of its input
func pct(p float32, m float32) float32 {
	return ((p / 100.0) * m)
}

func twrap(w *app.Window, width, height float32) error {
	var ts float32 = 2.8
	subsize := ts * 0.6
	gray := color.NRGBA{100, 100, 100, 255}
	red := color.NRGBA{128, 0, 0, 255}
	green := color.NRGBA{0, 128, 0, 255}
	blue := color.NRGBA{0, 0, 128, 255}
	s := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."

	for {
		e := w.NextEvent()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})

			var left float32 = 15
			var y1, y2, y3, ys float32

			y1 = 95
			y2 = y1 - 20
			y3 = y2 - 25
			ys = ts * 2

			canvas.Text(left, y1, subsize, "TextWrap(x, y, size, s, 60, red)", gray)
			canvas.Text(left, y2, subsize, "TextWrap(x, y, size, s, 40, red)", gray)
			canvas.Text(left, y3, subsize, "TextWrap(x, y, size, s, 20, red)", gray)

			red.A, green.A, blue.A = 50, 50, 50
			canvas.CornerRect(left, y1-ts, 60, ts*4, red)
			canvas.CornerRect(left, y2-ts, 40, ts*6, green)
			canvas.CornerRect(left, y3-ts, 20, ts*13, blue)
			red.A, green.A, blue.A = 255, 255, 255

			canvas.TextWrap(left, y1-ys, ts, 60, s, red)
			canvas.TextWrap(left, y2-ys, ts, 40, s, green)
			canvas.TextWrap(left, y3-ys, ts, 20, s, blue)
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
		w := app.NewWindow(app.Title("text wrapping"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := twrap(w, width, height); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()

}

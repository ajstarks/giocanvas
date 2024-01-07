// component charts
package main

import (
	"flag"
	"image/color"
	"io"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
	"github.com/ajstarks/giocanvas/chart"
)

func comp(canvas *giocanvas.Canvas) error {
	sr, err := os.Open("sine.d")
	if err != nil {
		return err
	}
	cr, err := os.Open("cosine.d")
	if err != nil {
		return err
	}
	sine, err := chart.DataRead(sr)
	if err != nil {
		return err
	}
	cosine, err := chart.DataRead(cr)
	if err != nil {
		return err
	}
	cosine.Zerobased = false
	sine.Zerobased = false
	cosine.Frame(canvas, 5)
	sine.Label(canvas, 1.5, 10)
	cosine.YAxis(canvas, 1.2, -1.0, 1.0, 1.0, "%0.2f", true)
	cosine.Color = color.NRGBA{0, 128, 0, 255}
	sine.Color = color.NRGBA{128, 0, 0, 255}
	cosine.Scatter(canvas, 0.5)
	sine.Scatter(canvas, 0.5)

	sine.Left = 10
	sine.Right = sine.Left + 40
	sine.Top, cosine.Top = 30, 30
	sine.Bottom, cosine.Bottom = 10, 10

	sine.CTitle(canvas, 2, 2)
	sine.Frame(canvas, 10)
	sine.Scatter(canvas, 0.25)

	offset := 45.0
	cosine.Left = sine.Left + offset
	cosine.Right = sine.Right + offset

	cosine.CTitle(canvas, 2, 2)
	cosine.Frame(canvas, 10)
	cosine.Scatter(canvas, 0.25)

	return nil
}

func sincos(w *app.Window, width, height float32) error {
	for {
		e := w.NextEvent()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), e)
			comp(canvas)
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
		w := app.NewWindow(app.Title("sine+cosine"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := sincos(w, width, height); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()

}

// component charts
package main

import (
	"flag"
	"image/color"
	"os"

	"gioui.org/app"
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
	sine.Label(canvas, 1.5, 10, "", "gray")
	cosine.YAxis(canvas, 1.2, -1.0, 1.0, 0.5, "%0.2f", true)
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

func sincos(width, height float32) error {
	w := new(app.Window)
	w.Option(app.Title("sine+cosine"), app.Size(unit.Dp(width), unit.Dp(height)))
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), e)
			comp(canvas)
			e.Frame(canvas.Context.Ops)
		case app.DestroyEvent:
			os.Exit(0)
		}
	}
}

func main() {
	var cw, ch int
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.Parse()
	go sincos(float32(cw), float32(ch))
	app.Main()
}

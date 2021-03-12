// component charts
package main

import (
	"flag"
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
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

func work(title string, width, height float32) {
	win := app.NewWindow(app.Title(title), app.Size(unit.Px(width), unit.Px(height)))
	for e := range win.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			os.Exit(0)
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, e)
			// your code here
			comp(canvas)

			e.Frame(canvas.Context.Ops)
		case key.Event:
			switch e.Name {
			case "Q", key.NameEscape:
				os.Exit(0)
			}
		}
	}
}

func main() {
	var w, h int
	flag.IntVar(&w, "width", 1000, "canvas width")
	flag.IntVar(&h, "height", 1000, "canvas height")
	flag.Parse()
	go work("sine+cosine", float32(w), float32(h))
	app.Main()
}

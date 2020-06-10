package main

import (
	"flag"
	"image/color"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

// pct returns the percentage of its input
func pct(p float32, m float32) float32 {
	return ((p / 100.0) * m)
}

func work(title string, w, h int) {
	width, height := float32(w), float32(h)
	win := app.NewWindow(app.Title(title), app.Size(unit.Px(width), unit.Px(height)))

	for e := range win.Events() {
		if e, ok := e.(system.FrameEvent); ok {
			canvas := giocanvas.NewCanvas(width, height, e.Config, e.Queue, e.Size)
			s := "Now is the time for all good men to come to the aid of the party, and the quick brown fox jumped over the lazy dog."
			canvas.TextWrap(20, 85, 3, 60, s, color.RGBA{128, 0, 0, 255})
			canvas.TextWrap(20, 65, 3, 40, s, color.RGBA{128, 0, 0, 255})
			canvas.TextWrap(20, 35, 3, 20, s, color.RGBA{128, 0, 0, 255})
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {
	var w, h int
	flag.IntVar(&w, "width", 1000, "canvas width")
	flag.IntVar(&h, "height", 1000, "canvas height")
	flag.Parse()
	go work("work", w, h)
	app.Main()
}

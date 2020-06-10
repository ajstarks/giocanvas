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
	var ts float32 = 3
	var lw float32 = 0.1
	gray := color.RGBA{128, 128, 128, 255}
	red := color.RGBA{128, 0, 0, 255}
	green := color.RGBA{0, 128, 0, 255}
	blue := color.RGBA{0, 0, 128, 255}
	for e := range win.Events() {
		if e, ok := e.(system.FrameEvent); ok {
			canvas := giocanvas.NewCanvas(width, height, e.Config, e.Queue, e.Size)
			s := "Now is the time for all good men to come to the aid of the party, and the quick brown fox jumped over the lazy dog."
			var left float32 = 20
			canvas.Line(left, 0, left, 100, lw*2, gray)
			canvas.Line(left+60, 85+ts, left+60, 0, lw, red)
			canvas.Line(left+40, 65+ts, left+40, 0, lw, green)
			canvas.Line(left+20, 35+ts, left+20, 0, lw, blue)
			canvas.TextWrap(left, 85, ts, 60, s, red)
			canvas.TextWrap(left, 65, ts, 40, s, green)
			canvas.TextWrap(left, 35, ts, 20, s, blue)
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

package main

import (
	"flag"
	"image/color"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func hello(title string, w, h int) {
	width, height := float32(w), float32(h)
	win := app.NewWindow(app.Title(title), app.Size(unit.Px(width), unit.Px(height)))
	for e := range win.Events() {
		if e, ok := e.(system.FrameEvent); ok {
			canvas := giocanvas.NewCanvas(width, height, e)
			canvas.CenterRect(50, 50, 100, 100, color.RGBA{0, 0, 0, 255})
			canvas.Circle(50, 0, 50, color.RGBA{0, 0, 255, 255})
			canvas.TextMid(50, 20, 10, "hello, world", color.RGBA{255, 255, 255, 0})
			canvas.CenterImage("earth.jpg", 50, 70, 1000, 1000, 30)
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {
	var w, h int
	flag.IntVar(&w, "width", 1000, "canvas width")
	flag.IntVar(&h, "height", 1000, "canvas height")
	flag.Parse()
	go hello("hello", w, h)
	app.Main()
}

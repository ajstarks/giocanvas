package main

import (
	"flag"
	"image/color"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func sunearth(s string, w, h int) {
	width := float32(w)
	height := float32(h)
	win := app.NewWindow(app.Title(s), app.Size(unit.Px(width), unit.Dp(height)))

	yellow := color.RGBA{255, 248, 231, 255}
	blue := color.RGBA{44, 77, 232, 255}
	black := color.RGBA{0, 0, 0, 255}

	var earthsize float32 = 0.8
	sunsize := earthsize * 109

	for e := range win.Events() {
		if e, ok := e.(system.FrameEvent); ok {
			canvas := giocanvas.NewCanvas(width, height, e.Config, e.Queue, e.Size)

			canvas.CenterRect(50, 50, 100, 100, black)
			canvas.Circle(100, 0, sunsize, yellow)
			canvas.Circle(30, 90, earthsize, blue)

			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {
	var w, h int
	flag.IntVar(&w, "width", 900, "canvas width")
	flag.IntVar(&h, "height", 1200, "canvas height")
	flag.Parse()
	go sunearth("sun+earth", w, h)
	app.Main()
}

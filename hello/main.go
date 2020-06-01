package main

import (
	"flag"
	"image/color"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func main() {
	var w, h int
	flag.IntVar(&w, "width", 1000, "canvas width")
	flag.IntVar(&h, "height", 1000, "canvas height")
	flag.Parse()
	width := float32(w)
	height := float32(h)
	size := app.Size(unit.Dp(width), unit.Dp(height))
	title := app.Title("hello")
	gofont.Register()
	go func() {
		w := app.NewWindow(title, size)
		for e := range w.Events() {
			if e, ok := e.(system.FrameEvent); ok {
				canvas := giocanvas.NewCanvas(width, height, e.Config, e.Queue, e.Size)
				canvas.CenterRect(50, 50, 100, 100, color.RGBA{0, 0, 0, 255})
				canvas.Circle(50, 0, 50, color.RGBA{0, 0, 255, 255})
				canvas.TextMid(50, 20, 10, "hello, world", color.RGBA{255, 255, 255, 0})
				canvas.CenterImage("earth.jpg", 50, 70, 1000, 1000, 30)
				e.Frame(canvas.Context.Ops)
			}
		}
	}()
	app.Main()
}

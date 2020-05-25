package main

import (
	"flag"
	"image/color"

	"gioui.org/app"
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

	go func() {
		w := app.NewWindow(title, size)
		canvas := giocanvas.NewCanvas(width, height)
		for e := range w.Events() {
			if e, ok := e.(system.FrameEvent); ok {
				canvas.Context.Reset(e.Queue, e.Config, e.Size)
				canvas.Rect(0, 0, width, height, color.RGBA{0, 0, 0, 255})
				canvas.TextMid(width/2, height*.15, width/10, "hello, world", color.RGBA{255, 255, 255, 0})
				canvas.CenterImage("earth.jpg", width/2, height*0.60, 1000, 1000, 50)
				e.Frame(canvas.Context.Ops)
			}
		}
	}()
	app.Main()
}

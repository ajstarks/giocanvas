// eclipse illustrates the eclipse
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
	flag.IntVar(&w, "width", 1200, "canvas width")
	flag.IntVar(&h, "height", 900, "canvas height")
	flag.Parse()
	width := float32(w)
	height := float32(h)
	size := app.Size(unit.Dp(width), unit.Dp(height))
	title := app.Title("Eclipse")
	black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}
	go func() {
		w := app.NewWindow(title, size)
		for e := range w.Events() {
			if e, ok := e.(system.FrameEvent); ok {
				canvas := giocanvas.NewCanvas(width, height, e.Config, e.Queue, e.Size)
				canvas.CenterRect(50, 50, 100, 100, black)
				var r float32 = 5.0
				var y float32 = 50.0
				var x float32 = 10.0
				for x = 10.0; x < 100.0; x += 15 {
					canvas.Circle(x, 50, r+0.5, white)
					canvas.Circle(x, y, r, black)
					y -= 2
				}
				e.Frame(canvas.Context.Ops)
			}
		}
	}()
	app.Main()

}

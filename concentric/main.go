// concentric circles
package main

import (
	"flag"
	"image/color"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func concentric(s string, w, h int) {
	width, height := float32(w), float32(h)
	win := app.NewWindow(app.Title(s), app.Size(unit.Px(width), unit.Px(height)))
	r := float32(30)
	g := uint8(5)
	for e := range win.Events() {
		if e, ok := e.(system.FrameEvent); ok {
			canvas := giocanvas.NewCanvas(width, height, e.Config, e.Queue, e.Size)
			canvas.CenterRect(50, 50, 100, 100, color.RGBA{0, 0, 0, 255})
			for i := 0; i < 6; i++ {
				canvas.Circle(50, 50, r, color.RGBA{128, g, g, 255})
				r -= 5
				g += 25
			}
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {
	var w, h int
	flag.IntVar(&w, "width", 1000, "canvas width")
	flag.IntVar(&h, "height", 1000, "canvas height")
	flag.Parse()
	go concentric("concentric", w, h)
	app.Main()
}

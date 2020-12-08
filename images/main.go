// hello is the giocanvas hello, world
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
)

func images(title string, width, height float32) {
	defer os.Exit(0)
	win := app.NewWindow(app.Title(title), app.Size(unit.Px(width), unit.Px(height)))
	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, e)
			var x, y, scale float32
			scale = 5.0
			canvas.Grid(0, 0, 100, 100, 0.1, 10, color.NRGBA{128, 128, 128, 255})
			for x = 10; x < 100; x += 10 {
				y = x
				canvas.CenterImage("earth.jpg", x, y, 1000, 1000, scale)
				scale += 2.5
			}
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
	go images("images", float32(w), float32(h))
	app.Main()
}

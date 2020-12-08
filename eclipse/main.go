// eclipse illustrates the eclipse
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

func eclipse(s string, width, height float32) {
	defer os.Exit(0)
	size := app.Size(unit.Px(width), unit.Px(height))
	title := app.Title(s)
	black := color.NRGBA{0, 0, 0, 255}
	white := color.NRGBA{255, 255, 255, 255}

	win := app.NewWindow(title, size)
	canvas := giocanvas.NewCanvas(width, height, system.FrameEvent{})
	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
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
	go eclipse("eclipse", float32(w), float32(h))
	app.Main()
}

// rgb shows RGB values
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

func rgb(title string, width, height float32) {
	win := app.NewWindow(app.Title(title), app.Size(unit.Px(width), unit.Px(height)))

	colortab := []string{
		"orange",
		"rgb(100)",
		"rgb(100,100)",
		"rgb(100,100,100)",
		"rgb(100,100,100,100)",
		"rgb()",
		"nonsense",
	}
	var x, y float32 = 50, 80
	canvas := giocanvas.NewCanvas(width, height, system.FrameEvent{})
	for e := range win.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			os.Exit(0)
		case system.FrameEvent:
			for _, c := range colortab {
				canvas.EText(x-10, y, 3, c, color.NRGBA{0, 0, 0, 255})
				canvas.Circle(x, y, 4, giocanvas.ColorLookup(c))
				y -= 10
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
	go rgb("rgb", float32(w), float32(h))
	app.Main()
}

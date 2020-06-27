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

func linetest(title string, w, h int) {
	defer os.Exit(0)
	width, height := float32(w), float32(h)
	win := app.NewWindow(app.Title(title), app.Size(unit.Px(width), unit.Px(height)))
	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, e)
			var x, y, lw, ls float32
			lw = 0.1
			ls = 1
			for y = 5; y <= 95; y += 5 {
				canvas.Line(50, 50, 95, y, lw, color.RGBA{128, 0, 0, 128})
				canvas.Line(50, 50, 5, y, lw, color.RGBA{0, 0, 128, 128})
				canvas.Coord(95, y, ls, "", color.RGBA{0, 0, 0, 255})
				canvas.Coord(5, y, ls, "", color.RGBA{0, 0, 0, 255})
				lw += 0.1
			}

			lw = 0.1
			for x = 5; x <= 95; x += 5 {
				canvas.Line(50, 50, x, 95, lw, color.RGBA{0, 128, 0, 128})
				canvas.Line(50, 50, x, 5, lw, color.RGBA{0, 0, 0, 128})
				canvas.Coord(x, 95, ls, "", color.RGBA{0, 0, 0, 255})
				canvas.Coord(x, 5, ls, "", color.RGBA{0, 0, 0, 255})
				lw += 0.1
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
	go linetest("linetest", w, h)
	app.Main()
}

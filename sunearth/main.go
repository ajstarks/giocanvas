// sunearth shows the relative size of the Sun and Earth
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

func sunearth(s string, width, height float32) {
	win := app.NewWindow(app.Title(s), app.Size(unit.Px(width), unit.Px(height)))

	yellow := color.NRGBA{255, 248, 231, 255}
	blue := color.NRGBA{44, 77, 232, 255}
	black := color.NRGBA{0, 0, 0, 255}

	var earthsize float32 = 0.8
	sunsize := earthsize * 109

	canvas := giocanvas.NewCanvas(width, height, system.FrameEvent{})
	for e := range win.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			os.Exit(0)
		case system.FrameEvent:
			canvas.CenterRect(50, 50, 100, 100, black)
			canvas.Circle(100, 0, sunsize, yellow)
			canvas.Circle(30, 90, earthsize, blue)

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
	go sunearth("sun+earth", float32(w), float32(h))
	app.Main()
}

// concentric circles
package main

import (
	"flag"
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/unit"
	gc "github.com/ajstarks/giocanvas"
)

func concentric(s string, width, height float32) {
	defer os.Exit(0)
	win := app.NewWindow(app.Title(s), app.Size(unit.Px(width), unit.Px(height)))

	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			canvas := gc.NewCanvas(width, height, e)
			canvas.Background(gc.ColorLookup("white"))
			var r float32 = 50
			for g := uint8(0); g < 250; g += 50 {
				canvas.Circle(50, 50, r, color.NRGBA{g, g, g, 255})
				r -= 10
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
	go concentric("concentric", float32(w), float32(h))
	app.Main()
}

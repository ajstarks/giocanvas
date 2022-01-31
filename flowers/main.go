// flowers demonstrated transforms with "flowers"
package main

import (
	"flag"
	"image/color"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
	gc "github.com/ajstarks/giocanvas"
)

func flower(canvas *gc.Canvas, x, y, w, h float32, fill color.NRGBA) {
	var a float32
	for a = 0; a < math.Pi*2; a += math.Pi / 9 {
		stack := canvas.Rotate(x, y, a)
		canvas.Ellipse(x, y, w, h, fill)
		gc.EndTransform(stack)
	}
}

func work(title string, width, height float32) {
	win := app.NewWindow(app.Title(title), app.Size(unit.Px(width), unit.Px(height)))
	red := color.NRGBA{128, 0, 0, 100}
	blue := color.NRGBA{0, 0, 128, 100}
	green := color.NRGBA{0, 128, 0, 100}
	orange := gc.ColorLookup("orange")
	bgcolor := gc.ColorLookup("linen")

	for e := range win.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			os.Exit(0)
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, system.FrameEvent{})
			canvas.Background(bgcolor)
			flower(canvas, 10, 90, 5, 1, red)
			flower(canvas, 25, 75, 10, 1.5, green)
			flower(canvas, 50, 50, 15, 3.0, blue)
			flower(canvas, 80, 20, 20, 4.5, orange)
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
	go work("flowers", float32(w), float32(h))
	app.Main()
}

// transform tests affine transforms
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
	gc "github.com/ajstarks/giocanvas"
)

func transforms(title string, width, height float32) {
	defer os.Exit(0)
	win := app.NewWindow(app.Title(title), app.Size(unit.Px(width), unit.Px(height)))

	var midx, rectw, recth, recty, ts, ts2 float32
	midx = 50
	rectw = 40
	recth = rectw / 4
	ts = 5
	ts2 = ts / 3
	black := color.RGBA{0, 0, 0, 255}
	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			canvas := gc.NewCanvas(width, height, e)

			recty = 90
			canvas.CenterRect(midx, recty, rectw, recth, color.RGBA{128, 128, 128, 128})
			canvas.TextMid(midx, recty-ts2, ts, "Reference", black)

			recty = 70
			stack := canvas.Scale(midx, recty, 2)
			canvas.CenterRect(midx, recty, rectw, recth, color.RGBA{0, 0, 128, 128})
			canvas.TextMid(midx, recty-ts2, ts, "scale", black)
			gc.EndTransform(stack)

			recty = 50
			stack = canvas.Shear(midx, midx, math.Pi/4, 0)
			canvas.CenterRect(midx, recty, rectw, recth, color.RGBA{128, 0, 0, 128})
			canvas.TextMid(midx, recty-ts2, ts, "shear", black)
			gc.EndTransform(stack)

			stack = canvas.Translate(20, 85)
			canvas.CenterRect(midx, recty, rectw, recth, color.RGBA{0, 128, 0, 128})
			canvas.TextMid(midx, recty-ts2, ts, "translate", black)
			gc.EndTransform(stack)

			recty = 20
			stack = canvas.Rotate(midx, recty, math.Pi/4)
			canvas.CenterRect(midx, recty, rectw, recth, color.RGBA{255, 50, 0, 200})
			canvas.TextMid(midx, recty-ts2, ts, "rotate", black)
			gc.EndTransform(stack)

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
	go transforms("transforms", float32(w), float32(h))
	app.Main()
}

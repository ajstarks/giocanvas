// colorwall: inspired by Ellsworth Kelly's "Colors for a Large Wall, 1951'
package main

import (
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

var (
	bgcolor = color.NRGBA{0xbb, 0xbb, 0xbb, 0xff}
	layout  = [][]string{
		{"#000000", "#eeeeee", "#735976", "#eeeeee", "#000000", "#af5d23", "#eeeeee", "#366e93"}, // row 1
		{"#eeeeee", "#03342f", "#000000", "#eeeeee", "#ccb04d", "#eeeeee", "#a74e4a", "#000000"}, // row 2
		{"#000000", "#eeeeee", "#eeeeee", "#391a32", "#eeeeee", "#eeeeee", "#eeeeee", "#af5d23"}, // row 3
		{"#8a1f1b", "#eeeeee", "#366e93", "#eeeeee", "#5e825e", "#000000", "#391a32", "#eeeeee"}, // row 4
		{"#eeeeee", "#391a32", "#000000", "#eeeeee", "#eeeeee", "#8a1f1b", "#eeeeee", "#122e63"}, // row 5
		{"#03342f", "#eeeeee", "#eeeeee", "#366e93", "#eeeeee", "#eeeeee", "#03342f", "#000000"}, // row 6
		{"#eeeeee", "#a74e4a", "#5e825e", "#eeeeee", "#000000", "#735976", "#eeeeee", "#eeeeee"}, // row 7
		{"#000000", "#eeeeee", "#391a32", "#ccb04d", "#eeeeee", "#000000", "#a74e4a", "#000000"}, // row 8
	}
)

func colorwall(width, height float32) {
	w := new(app.Window)
	appsize := app.Size(unit.Dp(width), unit.Dp(height))
	w.Option(app.Title("colorwall: inspired by Ellsworth Kelly's “Colors on a Large Wall”, 1951"), appsize)

	var x, y, left, right, top, bottom, xincr, yincr float32
	left, right, bottom, top = 25, 85, 20, 80
	nr, nc := 8, 8

	xincr = (right - left) / float32(nr)
	yincr = (top - bottom) / float32(nc)

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, app.FrameEvent{})
			canvas.Background(bgcolor)
			y = top
			for i := 0; i < nr; i++ {
				x = left
				for j := 0; j < nc; j++ {
					canvas.Square(x, y, yincr-0.1, giocanvas.ColorLookup(layout[i][j]))
					x += xincr
				}
				y -= yincr
			}
			e.Frame(canvas.Context.Ops)

		case app.DestroyEvent:
			os.Exit(0)
		}
	}
}

func main() {
	go colorwall(1000, 1000)
	app.Main()
}

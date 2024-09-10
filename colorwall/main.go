// colorwall: inspired by Ellsworth Kelly's "Colors for a Large Wall, 1951'
package main

import (
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

var palette = []string{
	"#000000",
	"#8a1f1b", "#a74e4a", "#af5d23", "#ccb04d", "#03342f",
	"#122e63", "#391a32", "#735976", "#5e825e", "#366e93",
	"#eeeeee",
}

func displaypalette(canvas *giocanvas.Canvas, x, y, xincr, yincr float32) {
	for i := 0; i < len(palette); i++ {
		canvas.Square(x, y, yincr, giocanvas.ColorLookup(palette[i]))
		canvas.CText(x, y-yincr, 1.5, palette[i], giocanvas.ColorLookup("#000000"))
		x += xincr
	}
}

func colorwall(width, height float32, showpalette bool) {
	w := new(app.Window)
	appsize := app.Size(unit.Dp(width), unit.Dp(height))
	w.Option(app.Title(`colorwall: inspired by Ellsworth Kelly's "Colors on a Large Wall", 1951`), appsize)
	layout := [][]string{
		{"#000000", "#eeeeee", "#735976", "#eeeeee", "#000000", "#af5d23", "#eeeeee", "#366e93"}, // row 1
		{"#eeeeee", "#03342f", "#000000", "#eeeeee", "#ccb04d", "#eeeeee", "#a74e4a", "#000000"}, // row 2
		{"#000000", "#eeeeee", "#eeeeee", "#391a32", "#eeeeee", "#eeeeee", "#eeeeee", "#af5d23"}, // row 3
		{"#8a1f1b", "#eeeeee", "#366e93", "#eeeeee", "#5e825e", "#000000", "#391a32", "#eeeeee"}, // row 4
		{"#eeeeee", "#391a32", "#000000", "#eeeeee", "#eeeeee", "#8a1f1b", "#eeeeee", "#122e63"}, // row 5
		{"#03342f", "#eeeeee", "#eeeeee", "#366e93", "#eeeeee", "#eeeeee", "#03342f", "#000000"}, // row 6
		{"#eeeeee", "#a74e4a", "#5e825e", "#eeeeee", "#000000", "#735976", "#eeeeee", "#eeeeee"}, // row 7
		{"#000000", "#eeeeee", "#391a32", "#ccb04d", "#eeeeee", "#000000", "#a74e4a", "#000000"}, // row 8
	}
	var x, y, left, right, top, bottom, xincr, yincr, nc, nr float32
	left, right, bottom, top = 25, 85, 20, 80
	nr, nc = 8, 8

	xincr = (right - left) / nr
	yincr = (top - bottom) / nc
	bgcolor := giocanvas.ColorLookup("#dddddd")
	basecolor := giocanvas.ColorLookup("#bbbbbb")
	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, app.FrameEvent{})
			canvas.Background(bgcolor)
			canvas.Square(51.25, 53.75, 60.5, basecolor)
			y = top
			for i := 0; i < len(layout); i++ {
				row := layout[i]
				x = left
				for j := 0; j < len(row); j++ {
					canvas.Square(x, y, yincr-0.1, giocanvas.ColorLookup(row[j]))
					x += xincr
				}
				y -= yincr
			}
			if showpalette {
				displaypalette(canvas, 10, 10, xincr, yincr)
			}
			e.Frame(canvas.Context.Ops)

		case app.DestroyEvent:
			os.Exit(0)
		}
	}
}

func main() {
	go colorwall(1000, 1000, false)
	app.Main()
}

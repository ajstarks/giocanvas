// confetti -- random shapes
package main

import (
	"flag"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func alpha(s string, width, height float32, color string) {
	size := app.Size(unit.Px(width), unit.Px(height))
	title := app.Title(s)
	win := app.NewWindow(title, size)
	blue := giocanvas.ColorLookup("steelblue")
	gray := giocanvas.ColorLookup("gray")
	dotcolor := giocanvas.ColorLookup(color)
	var x, y, px, dotsize, interval float32
	y = 50
	px = 2
	dotsize = 0.8
	interval = dotsize * 2.4
	for e := range win.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			os.Exit(0)
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, e)
			canvas.CText(50, y+12, 1.5, "Alpha", blue)
			canvas.CText(50, y-18, 1.5, "% Alpha", gray)
			for x = 0; x <= 100; x += 2 {
				dotcolor.A = uint8((x / 100) * 255)
				canvas.Circle(px, y, dotsize, dotcolor)
				canvas.TextMid(px, y-8, 0.75, strconv.FormatFloat(float64(x), 'g', -1, 32), gray)
				canvas.TextMid(px, y+5, 0.75, strconv.FormatInt(int64(dotcolor.A), 10), blue)
				px += interval
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
	var color string
	flag.IntVar(&w, "width", 2400, "canvas width")
	flag.IntVar(&h, "height", 600, "canvas height")
	flag.StringVar(&color, "color", "black", "color")
	flag.Parse()
	go alpha("alpha", float32(w), float32(h), color)
	app.Main()
}

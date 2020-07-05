// tile makes a tiling visual
package main

import (
	"flag"
	"math/rand"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/unit"
	gc "github.com/ajstarks/giocanvas"
)

type options struct {
	width, height                      int
	left, right, top, bottom, lw, step float32
	leftcolor, rightcolor, bgcolor     string
}

func main() {
	var opts options
	var lw, top, bottom, left, right, step float64
	flag.IntVar(&opts.width, "width", 1000, "canvas width")
	flag.IntVar(&opts.height, "height", 1000, "canvas height")
	flag.Float64Var(&lw, "linewidth", 0.5, "line width")
	flag.Float64Var(&step, "step", 5, "step")
	flag.Float64Var(&top, "top", 90, "top")
	flag.Float64Var(&bottom, "bottom", 10, "bottom")
	flag.Float64Var(&left, "left", 10, "left")
	flag.Float64Var(&right, "right", 90, "right")
	flag.StringVar(&opts.leftcolor, "leftcolor", "maroon", "left stroke color")
	flag.StringVar(&opts.rightcolor, "rightcolor", "maroon", "right stroke color")
	flag.StringVar(&opts.bgcolor, "bgcolor", "linen", "stroke color")
	flag.Parse()
	opts.lw, opts.step = float32(lw), float32(step)
	opts.top, opts.bottom, opts.right, opts.left = float32(top), float32(bottom), float32(right), float32(left)

	rand.Seed(time.Now().Unix() % 1e6)
	go tile("tile", opts)
	app.Main()
}

func tile(title string, opts options) {
	defer os.Exit(0)
	width := float32(opts.width)
	height := float32(opts.height)
	leftColor := gc.ColorLookup(opts.leftcolor)
	rightColor := gc.ColorLookup(opts.rightcolor)
	bgcolor := gc.ColorLookup(opts.bgcolor)

	win := app.NewWindow(app.Title(title), app.Size(unit.Px(width), unit.Px(height)))

	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			canvas := gc.NewCanvas(width, height, e)

			canvas.Background(bgcolor)
			var x, y float32
			for x = opts.left; x < opts.right; x += opts.step {
				for y = opts.bottom; y < opts.top; y += opts.step {
					if rand.Float64() >= 0.5 {
						canvas.Line(x, y, x+opts.step, y+opts.step, opts.lw, leftColor)
					} else {
						canvas.Line(x+opts.step, y, x, y+opts.step, opts.lw, rightColor)
					}
				}
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

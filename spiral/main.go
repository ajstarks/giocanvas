package main

import (
	"flag"
	"image/color"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func spiral(canvas *giocanvas.Canvas, cx, cy, dotsize float32, a, b, n, incr float64, color color.NRGBA) {
	for t := 0.0; t <= n*math.Pi; t += incr {
		sint, cost := math.Sincos(t)
		r := a + (b * t)
		x := cx + float32(r*cost)
		y := cy + float32(r*sint)
		canvas.Circle(x, y, dotsize, color)
	}
}

func work(title string, width, height float32, a, b, n, incr, dotsize float64, color string) {
	defer os.Exit(0)
	win := app.NewWindow(app.Title(title), app.Size(unit.Dp(width), unit.Dp(height)))
	var cx, cy float32 = 50, 50
	spcolor := giocanvas.ColorLookup(color)
	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), e)
			spiral(canvas, cx, cy, float32(dotsize), a, b, n, incr, spcolor)
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {
	var w, h int
	var a, b, n, incr, dotsize float64
	var dotcolor string
	flag.IntVar(&w, "width", 1000, "canvas width")
	flag.IntVar(&h, "height", 1000, "canvas height")
	flag.Float64Var(&a, "a", 1.5, "a")
	flag.Float64Var(&b, "b", 0.8, "b")
	flag.Float64Var(&n, "n", 16.0, "n")
	flag.Float64Var(&incr, "incr", 0.01, "increment")
	flag.Float64Var(&dotsize, "dot", 0.75, "dotsize")
	flag.StringVar(&dotcolor, "color", "gray", "dot color")
	flag.Parse()
	go work("spiral", float32(w), float32(h), a, b, n, incr, dotsize, dotcolor)
	app.Main()
}

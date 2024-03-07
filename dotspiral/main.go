package main

import (
	"flag"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

type config struct {
	start, end, r, rincr, dincr, tincr, dotsize float64
	dotcolor, bgcolor                           string
}

func dotspiral(canvas *giocanvas.Canvas, cx, cy float32, c config) {
	r := float32(c.r)
	rincr := float32(c.rincr)
	tincr := float32(c.tincr)
	dincr := float32(c.dincr)
	dotsize := float32(c.dotsize)
	start := float32(c.start)
	end := float32(c.end)
	color := giocanvas.ColorLookup(c.dotcolor)
	canvas.Background(giocanvas.ColorLookup(c.bgcolor))
	for t := start; t <= end; t += tincr {
		px, py := canvas.PolarDegrees(cx, cy, r, t)
		canvas.Circle(px, py, dotsize, color)
		r += rincr
		dotsize += dincr
	}
}

func work(title string, width, height float32, c config) {
	w := &app.Window{}
	w.Option(app.Title(title), app.Size(unit.Dp(width), unit.Dp(height)))
	var cx, cy float32 = 50, 50
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), e)
			dotspiral(canvas, cx, cy, c)
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {
	var w, h int
	var c config
	flag.IntVar(&w, "width", 1000, "canvas width")
	flag.IntVar(&h, "height", 1000, "canvas height")
	flag.Float64Var(&c.start, "start", 180, "start angle")
	flag.Float64Var(&c.end, "end", 360, "end angle")
	flag.Float64Var(&c.r, "r", 10.0, "radius")
	flag.Float64Var(&c.rincr, "rincr", 1.0, "radius increment")
	flag.Float64Var(&c.tincr, "tincr", 10.0, "angle increment")
	flag.Float64Var(&c.dincr, "dincr", 0.5, "size increment")
	flag.Float64Var(&c.dotsize, "size", 0.5, "dotsize")
	flag.StringVar(&c.dotcolor, "color", "rgb(128,0,0,128)", "dot color")
	flag.StringVar(&c.bgcolor, "bgcolor", "white", "background color")
	flag.Parse()
	go work("dotspiral", float32(w), float32(h), c)
	app.Main()
}

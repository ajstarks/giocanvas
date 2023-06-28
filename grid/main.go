// grid makes a grid
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func main() {
	var cw, ch int
	var x1, x2, y1, y2, xincr, yincr, lw float64
	var color string
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.Float64Var(&x1, "x1", 0, "x begin")
	flag.Float64Var(&x2, "x2", 100, "x end")
	flag.Float64Var(&y1, "y1", 0, "y begin")
	flag.Float64Var(&y2, "y2", 100, "y end")
	flag.Float64Var(&xincr, "xincr", 10, "x increment")
	flag.Float64Var(&yincr, "yincr", 10, "y increment")
	flag.Float64Var(&lw, "lw", 0.2, "line width")
	flag.StringVar(&color, "color", "black", "color")
	flag.Parse()

	width := float32(cw)
	height := float32(ch)

	go func() {
		w := app.NewWindow(app.Title("grid"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := grid(w, width, height, float32(x1), float32(x2), float32(y1), float32(y2), float32(xincr), float32(yincr), float32(lw), color); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

func grid(w *app.Window, width, height, x1, x2, y1, y2, xincr, yincr, lw float32, gridcolor string) error {
	color := giocanvas.ColorLookup(gridcolor)
	ts := xincr / 3
	for {
		ev := <-w.Events()
		switch e := ev.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), system.FrameEvent{})
			for y := y1; y <= y2; y += yincr {
				canvas.HLine(x1, y, x2-x1, lw, color)
				canvas.CText(x1+ts, y-(ts/2), ts, fmt.Sprintf("%0.f", y), color)
			}
			for x := x1; x <= x2; x += xincr {
				canvas.VLine(x, y1, y2-y1, lw, color)
				canvas.CText(x, y1+ts, ts, fmt.Sprintf("%0.f", x), color)
			}
			e.Frame(canvas.Context.Ops)
		}
	}
}

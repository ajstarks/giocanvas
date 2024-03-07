// alpha
package main

import (
	"flag"
	"io"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func alpha(w *app.Window, color string) error {
	blue := giocanvas.ColorLookup("steelblue")
	gray := giocanvas.ColorLookup("gray")
	dotcolor := giocanvas.ColorLookup(color)
	var x, y, px, dotsize, interval float32
	y = 50
	px = 2
	dotsize = 0.8
	interval = dotsize * 2.4
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})
			canvas.CText(50, y+12, 1.5, "Alpha", blue)
			canvas.CText(50, y-18, 1.5, "% Alpha", gray)
			px, y = 2, 50
			for x = 0; x <= 100; x += 2 {
				dotcolor.A = uint8((x / 100) * 255)
				canvas.Circle(px, y, dotsize, dotcolor)
				canvas.TextMid(px, y-8, 0.75, strconv.FormatFloat(float64(x), 'g', -1, 32), gray)
				canvas.TextMid(px, y+5, 0.75, strconv.FormatInt(int64(dotcolor.A), 10), blue)
				px += interval

			}
			e.Frame(canvas.Context.Ops)

		}
	}
}

func main() {
	var cw, ch int
	var color string
	flag.IntVar(&cw, "width", 1800, "canvas width")
	flag.IntVar(&ch, "height", 600, "canvas height")
	flag.StringVar(&color, "color", "black", "color")
	flag.Parse()
	width := float32(cw)
	height := float32(ch)

	go func() {
		w := &app.Window{}
		w.Option(app.Title("alpha"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := alpha(w, color); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

// lines tests line drawing
package main

import (
	"flag"
	"image/color"
	"io"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func linetest(w *app.Window, width, height float32) error {
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			os.Exit(0)
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, system.FrameEvent{})
			var x, y, lw, ls float32
			lw = 0.1
			ls = 1
			for y = 5; y <= 95; y += 5 {
				canvas.Line(50, 50, 95, y, lw, color.NRGBA{128, 0, 0, 128})
				canvas.Line(50, 50, 5, y, lw, color.NRGBA{0, 0, 128, 128})
				canvas.Coord(95, y, ls, "", color.NRGBA{0, 0, 0, 255})
				canvas.Coord(5, y, ls, "", color.NRGBA{0, 0, 0, 255})
				lw += 0.1
			}

			lw = 0.1
			for x = 5; x <= 95; x += 5 {
				canvas.Line(50, 50, x, 95, lw, color.NRGBA{0, 128, 0, 128})
				canvas.Line(50, 50, x, 5, lw, color.NRGBA{0, 0, 0, 128})
				canvas.Coord(x, 95, ls, "", color.NRGBA{0, 0, 0, 255})
				canvas.Coord(x, 5, ls, "", color.NRGBA{0, 0, 0, 255})
				lw += 0.1
			}
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {
	var cw, ch int
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.Parse()
	width := float32(cw)
	height := float32(ch)

	go func() {
		w := app.NewWindow(app.Title("lines"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := linetest(w, width, height); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()

}

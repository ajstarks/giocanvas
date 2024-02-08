// eclipse illustrates the eclipse
package main

import (
	"flag"
	"image/color"
	"io"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func eclipse(w *app.Window, width, height float32) error {
	black := color.NRGBA{0, 0, 0, 255}
	white := color.NRGBA{255, 255, 255, 255}
	for {
		e := w.NextEvent()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})
			canvas.CenterRect(50, 50, 100, 100, black)
			var r float32 = 5.0
			var y float32 = 50.0
			var x float32 = 10.0
			for x = 10.0; x < 100.0; x += 15 {
				canvas.Circle(x, 50, r+0.5, white)
				canvas.Circle(x, y, r, black)
				y -= 2
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
		w := app.NewWindow(app.Title("eclipse"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := eclipse(w, width, height); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

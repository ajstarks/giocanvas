// concentric circles
package main

import (
	"flag"
	"image/color"
	"io"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	gc "github.com/ajstarks/giocanvas"
)

func concentric(w *app.Window) error {
	for {
		e := w.NextEvent()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			canvas := gc.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})
			canvas.Background(gc.ColorLookup("white"))
			var r float32 = 50
			for g := uint8(0); g < 250; g += 50 {
				canvas.Circle(50, 50, r, color.NRGBA{g, g, g, 255})
				r -= 10
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
		w := app.NewWindow(app.Title("concentric"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := concentric(w); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

// hello is the giocanvas hello, world
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

func main() {
	var cw, ch int
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.Parse()

	width := float32(cw)
	height := float32(ch)

	go func() {
		w := app.NewWindow(app.Title("hello"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := hello(w, width, height); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

func hello(w *app.Window, width, height float32) error {
	black := color.NRGBA{0, 0, 0, 255}
	blue := color.NRGBA{0, 0, 255, 255}
	white := color.NRGBA{255, 255, 255, 255}
	for {
		switch e := w.NextEvent().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})
			canvas.CenterRect(50, 50, 100, 100, black)
			canvas.Circle(50, 0, 50, blue)
			canvas.TextMid(50, 20, 10, "hello, world", white)
			canvas.CenterImage("earth.jpg", 50, 70, 1000, 1000, 30)
			e.Frame(canvas.Context.Ops)
		}
	}
}

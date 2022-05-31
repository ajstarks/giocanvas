// hello is the giocanvas hello, world
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
	for {
		ev := <-w.Events()
		switch e := ev.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, system.FrameEvent{})
			canvas.CenterRect(50, 50, 100, 100, color.NRGBA{0, 0, 0, 255})
			canvas.Circle(50, 0, 50, color.NRGBA{0, 0, 255, 255})
			canvas.TextMid(50, 20, 10, "hello, world", color.NRGBA{255, 255, 255, 255})
			canvas.CenterImage("earth.jpg", 50, 70, 1000, 1000, 30)
			e.Frame(canvas.Context.Ops)
		}
	}
}

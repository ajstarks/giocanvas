// hello is the giocanvas hello, world
package main

import (
	"flag"
	"image/color"
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
	go hello(cw, ch)
	app.Main()
}

func hello(width, height int) {
	w := &app.Window{}
	w.Option(app.Title("hello"), app.Size(unit.Dp(float32(width)), unit.Dp(float32(height))))
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})
			canvas.Background(color.NRGBA{0, 0, 0, 255})
			canvas.Image("earth.jpg", 100, 0, 1000, 1000, 100)
			canvas.Text(10, 70, 10, "hello, world", color.NRGBA{255, 255, 255, 255})
			e.Frame(canvas.Context.Ops)
		}
	}
}

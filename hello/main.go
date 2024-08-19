// hello is the giocanvas hello, world
package main

import (
	"os"

	"gioui.org/app"
	"github.com/ajstarks/giocanvas"
)

func main() {
	go hello()
	app.Main()
}

func hello() {
	black := giocanvas.ColorLookup("black")
	white := giocanvas.ColorLookup("white")
	w := new(app.Window)
	w.Option(app.Title("hello"))
	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			w, h := float32(e.Size.X), float32(e.Size.Y)
			canvas := giocanvas.NewCanvas(w, h, app.FrameEvent{})
			canvas.Background(black)
			canvas.Image("earth.jpg", 100, 0, 1000, 1000, 100)
			canvas.Text(10, 70, 10, "hello, world", white)
			e.Frame(canvas.Context.Ops)
		case app.DestroyEvent:
			os.Exit(0)
		}
	}
}

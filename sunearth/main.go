// sunearth shows the relative size of the Sun and Earth
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

func sunearth(w *app.Window) error {

	yellow := color.NRGBA{255, 248, 231, 255}
	blue := color.NRGBA{44, 77, 232, 255}
	black := color.NRGBA{0, 0, 0, 255}

	var earthsize float32 = 0.8
	sunsize := earthsize * 109

	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})
			canvas.CenterRect(50, 50, 100, 100, black)
			canvas.Circle(100, 0, sunsize, yellow)
			canvas.Circle(30, 90, earthsize, blue)
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
		w := &app.Window{}
		w.Option(app.Title("sun+earth"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := sunearth(w); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

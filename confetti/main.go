// confetti -- random shapes
package main

import (
	"flag"
	"image/color"
	"io"
	"math/rand"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func rn(n int) float32 {
	return float32(rand.Intn(n))
}

func rn8(n int) uint8 {
	return uint8(rand.Intn(n))
}

func confetti(w *app.Window, width, height float32, nshapes, maxsize int) error {
	for {
		e := w.NextEvent()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), system.FrameEvent{})
			canvas.CenterRect(50, 50, 100, 100, color.NRGBA{0, 0, 0, 255})
			for i := 0; i < nshapes; i++ {
				color := color.NRGBA{rn8(255), rn8(255), rn8(255), rn8(255)}
				x, y := rn(100), rn(100)
				w, h := rn(maxsize), rn(maxsize)
				if i%2 == 0 {
					canvas.Ellipse(x, y, w, h, color)
				} else {
					canvas.CenterRect(x, y, w, h, color)
				}
			}
			e.Frame(canvas.Context.Ops)
		}

	}
}

func main() {
	var cw, ch, nshapes, maxsize int
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.IntVar(&nshapes, "n", 500, "number of shapes")
	flag.IntVar(&maxsize, "size", 10, "max size")
	flag.Parse()

	width := float32(cw)
	height := float32(ch)
	go func() {
		w := app.NewWindow(app.Title("confetti"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := confetti(w, width, height, nshapes, maxsize); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

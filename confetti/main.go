// confetti -- random shapes
package main

import (
	"flag"
	"image/color"
	"math/rand"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
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

func confetti(s string, w, h, nshapes, maxsize int) {
	defer os.Exit(0)
	width := float32(w)
	height := float32(h)
	size := app.Size(unit.Px(width), unit.Px(height))
	title := app.Title(s)
	win := app.NewWindow(title, size)
	canvas := giocanvas.NewCanvas(width, height, system.FrameEvent{})
	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
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
		case key.Event:
			switch e.Name {
			case "Q", key.NameEscape:
				os.Exit(0)
			}
		}

	}
}

func main() {
	var w, h, nshapes, maxsize int
	flag.IntVar(&w, "width", 1000, "canvas width")
	flag.IntVar(&h, "height", 1000, "canvas height")
	flag.IntVar(&nshapes, "n", 500, "number of shapes")
	flag.IntVar(&maxsize, "size", 10, "max size")
	flag.Parse()
	go confetti("Confetti", w, h, nshapes, maxsize)
	app.Main()
}

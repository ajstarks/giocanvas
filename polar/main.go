package main

import (
	"flag"
	"image/color"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func polar(x, y, r, t float32) (float32, float32) {
	fr := float64(r)
	px := fr * math.Cos(float64(t))
	py := fr * math.Sin(float64(t))

	return x + float32(px), y + float32(py)
}

func circles(title string, width, height float32) {
	defer os.Exit(0)
	win := app.NewWindow(app.Title(title), app.Size(unit.Px(width), unit.Px(height)))
	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, e)
			canvas.Background(color.RGBA{0, 0, 0, 255})
			var theta, radius float32
			for radius = 5; radius < 50; radius += 5 {
				for theta = 180; theta <= 360; theta += 15 {
					x, y := canvas.PolarDegrees(50, 50, radius, theta)
					canvas.Circle(x, y, radius/12, color.RGBA{128, 0, 0, 100})
				}
				for theta = math.Pi / 16; theta < math.Pi; theta += math.Pi / 16 {
					x, y := canvas.Polar(50, 50, radius, theta)
					canvas.Circle(x, y, radius/12, color.RGBA{0, 0, 128, 100})
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
	var w, h int
	flag.IntVar(&w, "width", 1000, "canvas width")
	flag.IntVar(&h, "height", 1000, "canvas height")
	flag.Parse()
	go circles("circles", float32(w), float32(h))
	app.Main()
}

package main

import (
	"flag"
	"fmt"
	"image/color"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func arc(title string, width, height float32) {
	defer os.Exit(0)
	win := app.NewWindow(app.Title(title), app.Size(unit.Px(width), unit.Px(height)))
	var x, y, step float32
	size := width / 10
	y = height * 0.2
	step = size + 10
	radius := width / 20
	bgcolor := color.RGBA{0, 0, 0, 50}
	arcolor := color.RGBA{0, 0, 128, 100}
	centercolor := color.RGBA{128, 128, 128, 128}
	black := color.RGBA{0, 0, 0, 255}
	for e := range win.Events() {
		if e, ok := e.(system.FrameEvent); ok {
			canvas := giocanvas.NewCanvas(width, height, e)
			angle := 45.0 // math.Pi / 4
			canvas.CText(50, 50, 10, "Arcs", black)
			for x = size; x < width-size; x += step {
				canvas.AbsCircle(x, y, radius, bgcolor)
				canvas.AbsCircle(x, y, 5, centercolor)
				canvas.AbsArc(x, y, radius, 0, radians(angle), arcolor)
				canvas.AbsTextMid(x, y+100, 18, fmt.Sprintf("%.1f°", angle), black)
				canvas.AbsTextMid(x, y-80, 18, fmt.Sprintf("%.4f rad", radians(angle)), black)
				angle += 45 // math.Pi / 4
			}

			y = 20
			angle = math.Pi * 2
			for x = 10; x < 90; x += 11 {
				canvas.Circle(x, y, 5, bgcolor)
				canvas.Circle(x, y, 0.5, centercolor)
				canvas.Arc(x, y, 5, 0, angle, arcolor)
				canvas.CText(x, y-10, 1.8, fmt.Sprintf("%.1f°", degrees(angle)), black)
				canvas.CText(x, y+8, 1.8, fmt.Sprintf("%.4f rad", angle), black)
				angle -= math.Pi / 4
			}

			e.Frame(canvas.Context.Ops)
		}
	}
}

func radians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func degrees(radians float64) float64 {
	return radians * (180 / math.Pi)
}

func main() {
	var w, h int
	flag.IntVar(&w, "width", 1000, "canvas width")
	flag.IntVar(&h, "height", 1000, "canvas height")
	flag.Parse()
	go arc("arc", float32(w), float32(h))
	app.Main()
}

package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

const pi = 3.14159265358979323846264338327950288419716939937510582097494459

func radians(degrees float64) float64 {
	return degrees * (pi / 180)
}

func degrees(radians float64) float64 {
	return radians * (180 / pi)
}

func arc(title string, width, height float32) {
	defer os.Exit(0)
	win := app.NewWindow(app.Title(title), app.Size(unit.Px(width), unit.Px(height)))
	var x, y, step float32
	size := width / 10
	step = width * 0.11
	ts := width * .01800
	radius := width / 20
	dotsize := width * 0.005
	bgcolor := color.RGBA{0, 0, 0, 50}
	arcolor := color.RGBA{0, 0, 128, 100}
	centercolor := color.RGBA{128, 128, 128, 128}
	black := color.RGBA{0, 0, 0, 255}

	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, e)
			angle := 45.0 // Pi / 4
			y = height * 0.2
			canvas.CText(50, 50, 10, "Arcs", black)
			for x = size; x < width-size; x += step {
				canvas.AbsCircle(x, y, radius, bgcolor)
				canvas.AbsCircle(x, y, dotsize, centercolor)
				canvas.AbsArc(x, y, radius, 0, radians(angle), arcolor)
				canvas.AbsTextMid(x, y+(width*0.10), ts, fmt.Sprintf("%.1f°", angle), black)
				canvas.AbsTextMid(x, y-(width*0.08), ts, fmt.Sprintf("%.4f rad", radians(angle)), black)
				angle += 45 // Pi / 4
			}

			y = 20
			angle = pi * 2
			for x = 10; x < 90; x += 11 {
				canvas.Circle(x, y, 5, bgcolor)
				canvas.Circle(x, y, 0.5, centercolor)
				canvas.Arc(x, y, 5, 0, angle, arcolor)
				canvas.CText(x, y-10, 1.8, fmt.Sprintf("%.1f°", degrees(angle)), black)
				canvas.CText(x, y+8, 1.8, fmt.Sprintf("%.4f rad", angle), black)
				angle -= pi / 4
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
	go arc("arc", float32(w), float32(h))
	app.Main()
}

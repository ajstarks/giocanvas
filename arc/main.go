package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
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

func absarc(c *giocanvas.Canvas, x, y, x1, y1, x2, y2, angle float32, fillcolor color.NRGBA) {
	p := new(clip.Path)
	ops := c.Context.Ops
	//r := f32.Rect(0, 0, c.Width, c.Height)
	defer op.Push(c.Context.Ops).Pop()
	p.Begin(ops)
	p.Move(f32.Pt(x, y))
	p.Arc(f32.Pt(x1, y1), f32.Pt(x2, y2), angle)
	p.Outline().Add(ops)
	paint.ColorOp{Color: fillcolor}.Add(ops)
	paint.PaintOp{}.Add(ops)
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
	bgcolor := color.NRGBA{0, 0, 0, 50}
	arcolor := color.NRGBA{0, 0, 128, 100}
	centercolor := color.NRGBA{128, 128, 128, 128}
	black := color.NRGBA{0, 0, 0, 255}

	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, e)
			//absarc(canvas, 250, 500, 100, 0, 50, 20, math.Pi, black)
			//absarc(canvas, 255, 500, 195, 0, math.Pi, color.NRGBA{255, 255, 255, 255})

			//canvas.Arc(50, 50, 50, radians(0), radians(180), centercolor)
			//canvas.CText(50, 30, 3, "giocanvas", black)
			//canvas.CText(50, 70, 3, "Arc", centercolor)
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
				canvas.ArcLine(x, y, 5, 0, angle, 0.2, color.NRGBA{128, 0, 0, 255})
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

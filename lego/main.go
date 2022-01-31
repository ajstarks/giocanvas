// lego charts
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

type pdata struct {
	name  string
	value float64
	fill  string
}

func wdata(canvas *giocanvas.Canvas, x, y, left, step float32, n int, fillcolor color.NRGBA) (float32, float32) {
	edge := (((step * 0.3) + step) * 7) + left
	for i := 0; i < n; i++ {
		if x > edge {
			x = left
			y -= step
		}
		op := fillcolor.A
		canvas.Circle(x, y, step*0.3, fillcolor)
		fillcolor.A = op - 30
		canvas.Square(x, y, step*0.9, fillcolor)
		fillcolor.A = op
		x += step
	}
	return x, y
}

func waffle(canvas *giocanvas.Canvas, data []pdata, left, top, step float32) {
	x := left
	y := top
	for _, d := range data {
		px, py := wdata(canvas, x, y, left, step, int(d.value), giocanvas.ColorLookup(d.fill))
		x = px
		y = py
	}
	y -= step * 2
	for _, d := range data {
		canvas.Circle(left, y, step*0.3, giocanvas.ColorLookup(d.fill))
		canvas.Text(left+step, y-step*0.2, step*0.5, fmt.Sprintf("%s (%.d%%)", d.name, int(d.value)), giocanvas.ColorLookup("rgb(120,120,120"))
		y -= step
	}
}

func grid(canvas *giocanvas.Canvas, x, y, step float32, bcolor, dcolor color.NRGBA) {
	sx := x
	for r := 0; r < 10; r++ {
		x = sx
		for c := 0; c < 10; c++ {
			canvas.Square(x, y, step*0.9, bcolor)
			canvas.Circle(x, y, step*0.3, dcolor)
			x += step
		}
		y -= step
	}
}

func lego(title string, width, height float32) {
	win := app.NewWindow(app.Title(title), app.Size(unit.Px(width), unit.Px(height)))
	incar := []pdata{
		{name: "White", value: 39, fill: "rgb(160,82,45,120)"},
		{name: "Hispanic", value: 19, fill: "rgb(160,82,45,180)"},
		{name: "Black", value: 40, fill: "rgb(160,82,45)"},
		{name: "Other", value: 2, fill: "rgb(180,180,180)"},
	}

	pop := []pdata{
		{name: "White", value: 64, fill: "rgb(160,82,45,120)"},
		{name: "Hispanic", value: 16, fill: "rgb(160,82,45,180)"},
		{name: "Black", value: 13, fill: "rgb(160,82,45)"},
		{name: "Other", value: 7, fill: "rgb(180,180,180)"},
	}
	var c1, c2 float32 = 15, 60

	for e := range win.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			os.Exit(0)
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, system.FrameEvent{})
			canvas.Text(c1, 85, 3, "Incarceration Rate", giocanvas.ColorLookup("black"))
			canvas.Text(c2, 85, 3, "US Population", giocanvas.ColorLookup("black"))
			canvas.CText(50, 20, 1.5, "Source: Breaking Down Mass Incarceration in the 2010 Census: State-by-State Incarceration Rates by Race/Ethnicity", giocanvas.ColorLookup("gray"))
			waffle(canvas, incar, c1, 80, 3)
			waffle(canvas, pop, c2, 80, 3)
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
	go lego("lego chart", float32(w), float32(h))
	app.Main()
}

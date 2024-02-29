// lego charts
package main

import (
	"flag"
	"fmt"
	"image/color"
	"io"
	"os"

	"gioui.org/app"
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

func lego(w *app.Window) error {
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

	for {
		e := w.NextEvent()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})
			canvas.Text(c1, 85, 3, "Incarceration Rate", giocanvas.ColorLookup("black"))
			canvas.Text(c2, 85, 3, "US Population", giocanvas.ColorLookup("black"))
			canvas.CText(50, 20, 1.5, "Source: Breaking Down Mass Incarceration in the 2010 Census: State-by-State Incarceration Rates by Race/Ethnicity", giocanvas.ColorLookup("gray"))
			waffle(canvas, incar, c1, 80, 3)
			waffle(canvas, pop, c2, 80, 3)
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
		w := app.NewWindow(app.Title("lego"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := lego(w); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

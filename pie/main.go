package main

import (
	"flag"
	"fmt"
	"image/color"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

const fullcircle = math.Pi * 2

type piedata struct {
	name  string
	value float64
	color string
}

func datasum(data []piedata) float64 {
	sum := 0.0
	for _, d := range data {
		sum += d.value
	}
	return sum
}

func piechart(canvas *giocanvas.Canvas, x, y, r float32, data []piedata) {
	sum := datasum(data)
	a1 := 0.0
	labelr := r + 10
	ts := r / 10
	for _, d := range data {
		color := giocanvas.ColorLookup(d.color)
		pct := (d.value / sum)
		a2 := (fullcircle * pct) + a1
		mid := fullcircle - (a1 + (a2-a1)/2)
		canvas.Arc(x, y, r, a1, a2, color)
		tx, ty := canvas.Polar(x, y, labelr, float32(mid))
		lx, ly := canvas.Polar(x, y, labelr-ts, float32(mid))
		canvas.CText(tx, ty, ts, fmt.Sprintf("%s (%.2f%%)", d.name, pct*100), color)
		canvas.Line(x, y, lx, ly, 0.1, color)
		a1 = a2
	}
}

func pie(title string, width, height float32) {
	defer os.Exit(0)
	win := app.NewWindow(app.Title(title), app.Size(unit.Px(width), unit.Px(height)))

	data := []piedata{
		{name: "Chrome", value: 65.47, color: "rgb(211,57,53)"},
		{name: "Safari", value: 16.97, color: "rgb(13,107,202)"},
		{name: "Other", value: 11.25, color: "rgb(150,150,150)"},
		{name: "Firefox", value: 4.25, color: "rgb(228,158,21)"},
		{name: "IE/Edge", value: 2.06, color: "rgb(0,128,0)"},
	}

	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, e)
			canvas.CText(50, 92, 4, "Browser Market Share, 2020-06", color.RGBA{20, 20, 20, 255})
			canvas.CText(50, 5, 2, "Source: Statcounter Global Stats, July 2020", color.RGBA{100, 100, 100, 255})
			piechart(canvas, 50, 50, 25, data)
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
	go pie("piechart", float32(w), float32(h))
	app.Main()
}

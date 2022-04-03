// pie chart
package main

import (
	"flag"
	"fmt"
	"image/color"
	"io"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

const fullcircle = 3.14159265358979323846264338327950288419716939937510582097494459 * 2

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
		p := (d.value / sum)
		angle := p * fullcircle
		a2 := a1 + angle
		mid := fullcircle - (a1 + (a2-a1)/2)
		canvas.Arc(x, y, r, a1, a2, color)
		//fmt.Fprintf(os.Stderr, "p=%.2f a1=%.2f a2=%.2f\n", p, a1, a2)
		tx, ty := canvas.Polar(x, y, labelr, float32(mid))
		lx, ly := canvas.Polar(x, y, labelr-ts, float32(mid))
		canvas.CText(tx, ty, ts, fmt.Sprintf("%s (%.2f%%)", d.name, p*100), color)
		canvas.Line(x, y, lx, ly, 0.1, color)
		a1 = a2

	}
}

func pie(w *app.Window, width, height float32) error {

	data := []piedata{
		{name: "Chrome", value: 65.47, color: "rgb(211,57,53)"},
		{name: "Safari", value: 16.97, color: "rgb(13,107,202)"},
		{name: "Other", value: 11.25, color: "rgb(150,150,150)"},
		{name: "Firefox", value: 4.25, color: "rgb(228,158,21)"},
		{name: "IE/Edge", value: 2.06, color: "rgb(0,128,0)"},
	}

	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			os.Exit(0)
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, system.FrameEvent{})
			canvas.CText(50, 92, 4, "Browser Market Share, 2020-06", color.NRGBA{20, 20, 20, 255})
			canvas.CText(50, 5, 2, "Source: Statcounter Global Stats, July 2020", color.NRGBA{100, 100, 100, 255})
			piechart(canvas, 50, 50, 25, data)
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
		w := app.NewWindow(app.Title("pie"), app.Size(unit.Px(width), unit.Px(height)))
		if err := pie(w, width, height); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

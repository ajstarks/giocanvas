package main

import (
	"flag"
	"fmt"
	"image/color"
	"math"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

// ChartData defines data
type ChartData struct {
	name  string
	value float64
}

func minmax(data []ChartData) (float64, float64) {
	min := data[0].value
	max := data[0].value
	for _, d := range data {
		if d.value > max {
			max = d.value
		}
		if d.value < min {
			min = d.value
		}
	}
	return min, max
}

func chart(canvas *giocanvas.Canvas, x, y, width, height float32, data []ChartData, interval int, datacolor string) {
	min, max := minmax(data)
	for i, d := range data {
		xp := float32(giocanvas.MapRange(float64(i), 0, float64(len(data)-1), float64(x), float64(width)))
		yp := float32(giocanvas.MapRange(d.value, min, max, float64(y), float64(y-height)))

		canvas.CenterRect(xp, yp, 5, 5, giocanvas.ColorLookup(datacolor))
		if interval > 0 && i%interval == 0 {
			canvas.TextMid(xp, y+5, 20, d.name, color.RGBA{0, 0, 0, 255})
		}
	}
}

func main() {
	var w, h int
	flag.IntVar(&w, "width", 1200, "canvas width")
	flag.IntVar(&h, "height", 900, "canvas height")
	flag.Parse()
	width := float32(w)
	height := float32(h)
	size := app.Size(unit.Dp(width), unit.Dp(height))
	title := app.Title("barchart")

	var (
		sinedata   []ChartData
		cosinedata []ChartData
		d          ChartData
	)
	for x := 0.0; x <= 2*math.Pi; x += math.Pi / 32 {
		d.name = fmt.Sprintf("%.1f", x)
		d.value = math.Sin(x)
		sinedata = append(sinedata, d)
	}

	for x := 0.0; x <= 2*math.Pi; x += math.Pi / 32 {
		d.name = fmt.Sprintf("%.1f", x)
		d.value = math.Cos(x)
		cosinedata = append(cosinedata, d)
	}

	go func() {
		w := app.NewWindow(title, size)
		canvas := giocanvas.NewCanvas(width, height)
		chbottom := height * 0.80
		chleft := width * 0.10
		chwidth := width * 0.90
		chheight := height * 0.50
		titley := height * 0.05
		labely := height * 0.15
		for e := range w.Events() {
			if e, ok := e.(system.FrameEvent); ok {
				canvas.Context.Reset(e.Queue, e.Config, e.Size)
				canvas.Text(chleft, titley, 40, "Sine and Cosine", giocanvas.ColorLookup("black"))
				canvas.Text(chleft, labely, 30, "sin(x)", color.RGBA{128, 128, 128, 255})
				canvas.HLine(chleft+width*0.10, labely+20, 15, 15, color.RGBA{255, 0, 0, 255})
				labely += width * 0.05
				canvas.Text(chleft, labely, 30, "cos(x)", color.RGBA{128, 128, 128, 255})
				canvas.HLine(chleft+width*0.10, labely+20, 15, 15, color.RGBA{0, 0, 255, 255})

				chart(canvas, chleft, chbottom, chwidth, chheight, sinedata, 5, "red")
				chart(canvas, chleft, chbottom, chwidth, chheight, cosinedata, 0, "blue")
				canvas.Grid(width, height, 2, 20, color.RGBA{0, 0, 0, 128})
				e.Frame(canvas.Context.Ops)
			}
		}
	}()
	app.Main()
}

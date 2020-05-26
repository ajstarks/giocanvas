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

func chart(canvas *giocanvas.Canvas, x, y, width, height float32, data []ChartData, interval int, datacolor color.RGBA) {
	min, max := minmax(data)
	for i, d := range data {
		xp := float32(giocanvas.MapRange(float64(i), 0, float64(len(data)-1), float64(x), float64(width)))
		yp := float32(giocanvas.MapRange(d.value, min, max, float64(y), float64(height)))
		canvas.Rect(x, y, width-x, 0.1, color.RGBA{0, 0, 0, 128})
		canvas.CenterRect(xp, yp, 0.3, 0.5, datacolor)
		if interval > 0 && i%interval == 0 {
			canvas.TextMid(xp, y-3, 1.5, d.name, color.RGBA{0, 0, 0, 255})
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
	title := app.Title("sine and cosine")

	var (
		sinedata   []ChartData
		cosinedata []ChartData
		d          ChartData
	)
	for x := 0.0; x <= 2*math.Pi; x += 0.05 {
		d.name = fmt.Sprintf("%.2f", x)
		d.value = math.Sin(x)
		sinedata = append(sinedata, d)
		d.value = math.Cos(x)
		cosinedata = append(cosinedata, d)
	}

	go func() {
		w := app.NewWindow(title, size)
		canvas := giocanvas.NewCanvas(width, height)
		black := color.RGBA{0, 0, 0, 255}
		red := color.RGBA{255, 0, 0, 255}
		blue := color.RGBA{0, 0, 255, 255}
		gcolor := color.RGBA{0, 0, 0, 100}
		for e := range w.Events() {
			if e, ok := e.(system.FrameEvent); ok {
				canvas.Context.Reset(e.Queue, e.Config, e.Size)
				canvas.Text(10, 90, 3, "Sine and Cosine", black)
				canvas.Text(10, 84, 2.5, "sin(x)", red)
				canvas.Text(10, 79, 2.5, "cos(x)", blue)
				canvas.HLine(20, 85, 2, 1, red)
				canvas.HLine(20, 80, 2, 1, blue)
				chart(canvas, 10, 15, 90, 70, sinedata, 16, red)
				chart(canvas, 10, 15, 90, 70, cosinedata, 0, blue)
				canvas.Grid(0, 0, 100, 100, 0.1, 5, gcolor)
				e.Frame(canvas.Context.Ops)
			}
		}
	}()
	app.Main()
}

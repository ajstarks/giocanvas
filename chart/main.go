package main

import (
	"flag"
	"fmt"
	"image/color"
	"math"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/unit"
	gc "github.com/ajstarks/giocanvas"
)

// ChartData defines data
type ChartData struct {
	name  string
	value float64
}

// ChartOptions define all the components of a chart
type ChartOptions struct {
	showtitle, showscatter, showarea, showframe, showlegend bool
	title                                                   string
	xlabelInterval                                          int
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

func xaxis(canvas *gc.Canvas, x, y, width, height float32, interval int, data []ChartData) {
	for i, d := range data {
		xp := float32(gc.MapRange(float64(i), 0, float64(len(data)-1), float64(x), float64(width)))
		if interval > 0 && i%interval == 0 {
			canvas.TextMid(xp, y-3, 1.5, d.name, color.RGBA{0, 0, 0, 255})
			canvas.Line(xp, y, xp, height, 0.1, color.RGBA{0, 0, 0, 128})
		}
	}
}

func frame(canvas *gc.Canvas, x, y, width, height float32, color color.RGBA) {
	canvas.Rect(x, height, width-x, height-y, color)
}

func dotchart(canvas *gc.Canvas, x, y, width, height float32, data []ChartData, datacolor color.RGBA) {
	min, max := minmax(data)
	for i, d := range data {
		xp := float32(gc.MapRange(float64(i), 0, float64(len(data)-1), float64(x), float64(width)))
		yp := float32(gc.MapRange(d.value, min, max, float64(y), float64(height)))
		canvas.Circle(xp, yp, 0.3, datacolor)
	}
}

func barchart(canvas *gc.Canvas, x, y, width, height float32, data []ChartData, interval int, datacolor color.RGBA) {
	min, max := minmax(data)
	for i, d := range data {
		xp := float32(gc.MapRange(float64(i), 0, float64(len(data)-1), float64(x), float64(width)))
		yp := float32(gc.MapRange(d.value, min, max, float64(y), float64(height)))
		canvas.VLine(xp, y, yp-y, 0.1, datacolor)
	}
}

func areachart(canvas *gc.Canvas, x, y, width, height float32, data []ChartData, datacolor color.RGBA) {
	min, max := minmax(data)
	l := len(data)

	ax := make([]float32, l+2)
	ay := make([]float32, l+2)
	ax[0] = x
	ay[0] = y
	ax[l+1] = width
	ay[l+1] = y

	for i, d := range data {
		xp := float32(gc.MapRange(float64(i), 0, float64(len(data)-1), float64(x), float64(width)))
		yp := float32(gc.MapRange(d.value, min, max, float64(y), float64(height)))
		ax[i+1] = xp
		ay[i+1] = yp
	}
	canvas.Polygon(ax, ay, datacolor)
}

func chart(s string, w, h int, chartopts ChartOptions) {
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
	gofont.Register()
	win := app.NewWindow(title, size)
	black := color.RGBA{0, 0, 0, 255}
	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}
	framecolor := color.RGBA{0, 0, 0, 20}
	for e := range win.Events() {
		if e, ok := e.(system.FrameEvent); ok {
			canvas := gc.NewCanvas(width, height, e.Config, e.Queue, e.Size)
			if chartopts.showtitle {
				canvas.Text(10, 90, 3, chartopts.title, black)
			}
			if chartopts.showlegend {
				canvas.Text(10, 84, 2.5, "sin(x)", red)
				canvas.Text(10, 79, 2.5, "cos(x)", blue)
				canvas.HLine(20, 85, 2, 1, red)
				canvas.HLine(20, 80, 2, 1, blue)
			}
			if chartopts.xlabelInterval > 0 {
				xaxis(canvas, 10, 15, 90, 70, chartopts.xlabelInterval, sinedata)
			}
			if chartopts.showframe {
				frame(canvas, 10, 15, 90, 70, framecolor)
			}
			if chartopts.showscatter {
				dotchart(canvas, 10, 15, 90, 70, sinedata, red)
				dotchart(canvas, 10, 15, 90, 70, cosinedata, blue)
			}
			if chartopts.showarea {
				red.A = 100
				areachart(canvas, 10, 15, 90, 70, sinedata, red)
				blue.A = 100
				areachart(canvas, 10, 15, 90, 70, cosinedata, blue)
			}
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {
	var opts ChartOptions
	var w, h int
	flag.IntVar(&w, "width", 1200, "canvas width")
	flag.IntVar(&h, "height", 900, "canvas height")
	flag.IntVar(&opts.xlabelInterval, "xlabel", 0, "show x axis")
	flag.StringVar(&opts.title, "chartitle", "Sine and Cosine", "chart title")
	flag.BoolVar(&opts.showtitle, "title", true, "show title")
	flag.BoolVar(&opts.showlegend, "legend", true, "show legend")
	flag.BoolVar(&opts.showarea, "area", false, "show area chart")
	flag.BoolVar(&opts.showscatter, "scatter", false, "show scatter chart")
	flag.BoolVar(&opts.showframe, "frame", false, "show frame")
	flag.Parse()
	go chart("chart", w, h, opts)
	app.Main()
}

package main

import (
	"bufio"
	"flag"
	"fmt"
	"image/color"
	"io"
	"os"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/unit"
	gc "github.com/ajstarks/giocanvas"
)

// NameValue defines data
type NameValue struct {
	name  string
	note  string
	value float64
}

// ChartOptions define all the components of a chart
type ChartOptions struct {
	showtitle, showscatter, showarea, showframe, showlegend, showbar bool
	title, legend, color                                             string
	xlabelInterval                                                   int
}

func minmax(data []NameValue) (float64, float64) {
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

// DataRead reads tab separated values into a NameValue slice
func DataRead(r io.Reader) ([]NameValue, error) {
	var d NameValue
	var data []NameValue
	var err error
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) == 0 { // skip blank lines
			continue
		}
		if t[0] == '#' && len(t) > 2 { // process titles
			// title = strings.TrimSpace(t[1:])
			continue
		}
		fields := strings.Split(t, "\t")
		if len(fields) < 2 {
			continue
		}
		if len(fields) == 3 {
			d.note = fields[2]
		} else {
			d.note = ""
		}
		d.name = fields[0]
		d.value, err = strconv.ParseFloat(fields[1], 64)
		if err != nil {
			d.value = 0
		}
		data = append(data, d)
	}
	err = scanner.Err()
	return data, err
}

func xaxis(canvas *gc.Canvas, x, y, width, height float32, interval int, data []NameValue) {
	for i, d := range data {
		xp := float32(gc.MapRange(float64(i), 0, float64(len(data)-1), float64(x), float64(width)))
		if interval > 0 && i%interval == 0 {
			canvas.TextMid(xp, y-3, 1.5, d.name, color.RGBA{0, 0, 0, 255})
			canvas.Line(xp, y, xp, height, 0.1, color.RGBA{0, 0, 0, 128})
		}
	}
	canvas.Line(x, height, width, height, 0.1, color.RGBA{0, 0, 0, 128})
	canvas.Line(width, height, width, y, 0.1, color.RGBA{0, 0, 0, 128})
}

func frame(canvas *gc.Canvas, x, y, width, height float32, color color.RGBA) {
	canvas.Rect(x, height, width-x, height-y, color)
}

func dotchart(canvas *gc.Canvas, x, y, width, height float32, data []NameValue, datacolor color.RGBA) {
	min, max := minmax(data)
	for i, d := range data {
		xp := float32(gc.MapRange(float64(i), 0, float64(len(data)-1), float64(x), float64(width)))
		yp := float32(gc.MapRange(d.value, min, max, float64(y), float64(height)))
		canvas.Circle(xp, yp, 0.3, datacolor)
	}
}

func barchart(canvas *gc.Canvas, x, y, width, height float32, data []NameValue, datacolor color.RGBA) {
	min, max := minmax(data)
	for i, d := range data {
		xp := float32(gc.MapRange(float64(i), 0, float64(len(data)-1), float64(x), float64(width)))
		yp := float32(gc.MapRange(d.value, min, max, float64(y), float64(height)))
		canvas.VLine(xp, y, yp-y, 0.1, datacolor)
	}
}

func areachart(canvas *gc.Canvas, x, y, width, height float32, data []NameValue, datacolor color.RGBA) {
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
	datacolor.A = 128
	canvas.Polygon(ax, ay, datacolor)
}

func chart(s string, w, h int, data []NameValue, chartopts ChartOptions) {
	width := float32(w)
	height := float32(h)
	size := app.Size(unit.Px(width), unit.Px(height))
	title := app.Title(s)
	gofont.Register()
	win := app.NewWindow(title, size)
	black := color.RGBA{0, 0, 0, 255}
	datacolor := gc.ColorLookup(chartopts.color)
	framecolor := color.RGBA{0, 0, 0, 20}
	for e := range win.Events() {
		if e, ok := e.(system.FrameEvent); ok {
			canvas := gc.NewCanvas(width, height, e.Config, e.Queue, e.Size)
			if chartopts.showtitle {
				canvas.Text(10, 90, 3, chartopts.title, black)
			}
			if chartopts.showlegend {
				canvas.Text(10, 84, 2.5, chartopts.legend, datacolor)
				canvas.HLine(20, 85, 2, 1, datacolor)
			}
			if chartopts.xlabelInterval > 0 {
				xaxis(canvas, 10, 15, 90, 70, chartopts.xlabelInterval, data)
			}
			if chartopts.showframe {
				frame(canvas, 10, 15, 90, 70, framecolor)
			}
			if chartopts.showscatter {
				dotchart(canvas, 10, 15, 90, 70, data, datacolor)
			}
			if chartopts.showarea {
				areachart(canvas, 10, 15, 90, 70, data, datacolor)
			}
			if chartopts.showbar {
				barchart(canvas, 10, 15, 90, 70, data, datacolor)
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

	flag.StringVar(&opts.title, "chartitle", "", "chart title")
	flag.StringVar(&opts.legend, "chartlegend", "", "chart legend")
	flag.StringVar(&opts.color, "color", "lightsteelblue", "chart data color")

	flag.BoolVar(&opts.showtitle, "title", true, "show title")
	flag.BoolVar(&opts.showlegend, "legend", false, "show legend")
	flag.BoolVar(&opts.showbar, "bar", false, "show bar chart")
	flag.BoolVar(&opts.showarea, "area", false, "show area chart")
	flag.BoolVar(&opts.showscatter, "scatter", false, "show scatter chart")
	flag.BoolVar(&opts.showframe, "frame", false, "show frame")
	flag.Parse()

	data, err := DataRead(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	go chart("charts", w, h, data, opts)
	app.Main()
}

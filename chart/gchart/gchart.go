// fchart -- command line chart using fc chart packages
package main

import (
	"flag"
	"fmt"
	"image/color"
	"io"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
	"github.com/ajstarks/giocanvas/chart"
)

type chartOptions struct {
	top, bottom, left, right                                                 float64
	barwidth, linewidth, linespacing, dotsize, textsize, ty, frameOp, areaOp float64
	bgcolor, dcolor, chartitle, yaxfmt, yrange                               string
	xlabel                                                                   int
	zb, line, bar, hbar, scatter, area, showtitle, showgrid                  bool
}

func main() {

	// Command line options
	var opts chartOptions
	var width, height int

	flag.IntVar(&width, "w", 1000, "canvas width")
	flag.IntVar(&height, "h", 1000, "canvas height")
	flag.IntVar(&opts.xlabel, "xlabel", 1, "x-xaxis label")
	flag.Float64Var(&opts.barwidth, "barwidth", 0.5, "bar width")
	flag.Float64Var(&opts.linewidth, "linewidth", 0.25, "bar width")
	flag.Float64Var(&opts.linespacing, "ls", opts.barwidth*4, "bar width")
	flag.Float64Var(&opts.dotsize, "dotsize", 0.5, "bar width")
	flag.Float64Var(&opts.textsize, "textsize", 1.5, "bar width")
	flag.Float64Var(&opts.top, "top", 80, "bar width")
	flag.Float64Var(&opts.bottom, "bottom", 40, "bar width")
	flag.Float64Var(&opts.left, "left", 10, "bar width")
	flag.Float64Var(&opts.right, "right", 90, "bar width")
	flag.Float64Var(&opts.ty, "ty", 5, "title position relative to the top")
	flag.Float64Var(&opts.frameOp, "frame", 0, "frame opacity")
	flag.Float64Var(&opts.areaOp, "areaop", 50, "area opacity")
	flag.StringVar(&opts.yrange, "yrange", "", "y axis range (min,max,step")
	flag.StringVar(&opts.chartitle, "chartitle", "", "chart title")
	flag.StringVar(&opts.yaxfmt, "yfmt", "%v", "yaxis format")
	flag.StringVar(&opts.dcolor, "color", "steelblue", "color")
	flag.StringVar(&opts.bgcolor, "bgcolor", "white", "background color")
	flag.BoolVar(&opts.showtitle, "title", true, "show the title")
	flag.BoolVar(&opts.showgrid, "grid", false, "show y axis grid")
	flag.BoolVar(&opts.zb, "zero", true, "zero minumum")
	flag.BoolVar(&opts.area, "area", false, "area chart")
	flag.BoolVar(&opts.bar, "bar", false, "bar chart")
	flag.BoolVar(&opts.line, "line", false, "line chart")
	flag.BoolVar(&opts.hbar, "hbar", false, "horizontal bar")
	flag.BoolVar(&opts.scatter, "scatter", false, "scatter chart")
	flag.Parse()

	var input io.Reader
	var ferr error

	// Read from stdin or specified file
	if len(flag.Args()) == 0 {
		input = os.Stdin
	} else {
		input, ferr = os.Open(flag.Args()[0])
		if ferr != nil {
			fmt.Fprintf(os.Stderr, "%v\n", ferr)
			os.Exit(1)
		}
	}
	data, err := chart.DataRead(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}

	// make the chart
	go gchart("charts", width, height, data, opts)
	app.Main()
}

func gchart(s string, w, h int, data chart.ChartBox, opts chartOptions) {
	width := float32(w)
	height := float32(h)
	appsize := app.Size(unit.Px(width), unit.Px(height))
	apptitle := app.Title(fmt.Sprintf("Chart: %s", data.Title))
	win := app.NewWindow(apptitle, appsize)

	// Define the colors
	datacolor := giocanvas.ColorLookup(opts.dcolor)
	labelcolor := color.RGBA{100, 100, 100, 255}
	bgcolor := giocanvas.ColorLookup(opts.bgcolor)

	// Set the chart attributes
	data.Zerobased = opts.zb
	data.Top, data.Bottom, data.Left, data.Right = opts.top, opts.bottom, opts.left, opts.right

	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, e)
			canvas.Background(bgcolor)

			// Draw the data
			data.Color = datacolor
			if opts.frameOp > 0 {
				data.Frame(canvas, opts.frameOp)
			}
			if opts.line {
				data.Line(canvas, opts.linewidth)
			}
			if opts.bar {
				data.Bar(canvas, opts.barwidth)
			}
			if opts.scatter {
				data.Scatter(canvas, opts.dotsize)
			}
			if opts.hbar {
				data.HBar(canvas, opts.barwidth, opts.linespacing, opts.textsize)
			}
			if opts.area {
				data.Area(canvas, opts.areaOp)
			}

			// Draw labels, axes if specified
			data.Color = labelcolor
			if opts.line || opts.bar || opts.scatter {
				data.Label(canvas, opts.textsize, opts.xlabel)
				if len(opts.yrange) > 0 {
					var yaxmin, yaxmax, yaxstep float64
					n, err := fmt.Sscanf(opts.yrange, "%v,%v,%v", &yaxmin, &yaxmax, &yaxstep)
					if n == 3 && err == nil {
						data.YAxis(canvas, opts.textsize, yaxmin, yaxmax, yaxstep, opts.yaxfmt, opts.showgrid)
					}
				}
			}

			// Draw the chart titles
			if len(opts.chartitle) > 0 {
				data.Title = opts.chartitle
			}
			if opts.showtitle && len(data.Title) > 0 {
				data.CTitle(canvas, opts.textsize*2, opts.ty)
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

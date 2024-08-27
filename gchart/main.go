// gchart -- command line chart using giocanvas chart package
package main

import (
	"flag"
	"io"
	"os"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
	"github.com/ajstarks/giocanvas/chart"
)

type chartOptions struct {
	top, bottom, left, right                                                          float64
	barwidth, linewidth, linespacing, dotsize, textsize, piesize, ty, frameOp, areaOp float64
	bgcolor, dcolor, labelcolor, chartitle, yaxfmt, yrange                            string
	xlabel                                                                            int
	zb, line, bar, hbar, scatter, area, pie, lego, showtitle, showgrid                bool
}

func main() {

	// Command line options
	var opts chartOptions
	var width, height int

	flag.IntVar(&width, "w", 1000, "canvas width")
	flag.IntVar(&height, "h", 1000, "canvas height")
	flag.IntVar(&opts.xlabel, "xlabel", 1, "x-axis label interval")
	flag.Float64Var(&opts.barwidth, "barwidth", 0.5, "bar width")
	flag.Float64Var(&opts.linewidth, "linewidth", 0.25, "line width")
	flag.Float64Var(&opts.linespacing, "ls", opts.barwidth*4, "line spacing")
	flag.Float64Var(&opts.dotsize, "dotsize", 0.5, "dot size")
	flag.Float64Var(&opts.textsize, "textsize", 1.5, "text size")
	flag.Float64Var(&opts.top, "top", 80, "chart top")
	flag.Float64Var(&opts.bottom, "bottom", 20, "chart bottom")
	flag.Float64Var(&opts.left, "left", 20, "chart left")
	flag.Float64Var(&opts.right, "right", 80, "chart right")
	flag.Float64Var(&opts.ty, "ty", 5, "title position relative to the top")
	flag.Float64Var(&opts.frameOp, "frame", 0, "frame opacity")
	flag.Float64Var(&opts.areaOp, "areaop", 50, "area opacity")
	flag.Float64Var(&opts.piesize, "piesize", 20, "pie chart radius")
	flag.StringVar(&opts.yrange, "yrange", "", "y axis range (min,max,step)")
	flag.StringVar(&opts.chartitle, "chartitle", "", "chart title")
	flag.StringVar(&opts.yaxfmt, "yfmt", "%v", "yaxis format")
	flag.StringVar(&opts.dcolor, "color", "steelblue", "color")
	flag.StringVar(&opts.bgcolor, "bgcolor", "white", "background color")
	flag.StringVar(&opts.labelcolor, "labelcolor", "rgb(100,100,100)", "label color")
	flag.BoolVar(&opts.showtitle, "title", true, "show the title")
	flag.BoolVar(&opts.showgrid, "grid", false, "show y axis grid")
	flag.BoolVar(&opts.zb, "zero", true, "zero minumum")
	flag.BoolVar(&opts.lego, "lego", false, "lego chart")
	flag.BoolVar(&opts.area, "area", false, "area chart")
	flag.BoolVar(&opts.bar, "bar", false, "bar chart")
	flag.BoolVar(&opts.line, "line", false, "line chart")
	flag.BoolVar(&opts.hbar, "hbar", false, "horizontal bar")
	flag.BoolVar(&opts.scatter, "scatter", false, "scatter chart")
	flag.BoolVar(&opts.pie, "pie", false, "show a pie chart")
	flag.Parse()

	var input io.Reader
	var ferr error
	var infile string

	// Read from stdin or specified file
	if len(flag.Args()) == 0 {
		input = os.Stdin
		infile = "stdin"
	} else {
		infile = flag.Args()[0]
		input, ferr = os.Open(infile)
		if ferr != nil {
			perr("unable to open ", infile)
			os.Exit(1)
		}
	}
	// read the data
	data, err := chart.DataRead(input)
	if err != nil {
		perr("unable to read ", infile)
		os.Exit(2)
	}
	// must specify at least one of line, bar, hbar, scatter, area, pie, lego
	if !(opts.line || opts.scatter || opts.bar || opts.area || opts.hbar || opts.lego || opts.pie) {
		perr("pick a chart type (-line, -bar, -hbar, -area, -scatter, -lego, -pie)", infile)
		os.Exit(3)
	}
	// make the chart
	go gchart("charts", width, height, data, opts)
	app.Main()
}

// perr prints a filename and message to stderr
func perr(msg, file string) {
	io.WriteString(os.Stderr, file+": "+msg+"\n")
}

// string to floating point
func stof(s string) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

// yr parses the the yrange (max,min,step) string
func yr(yrange string, dmin, dmax float64) (float64, float64, float64) {
	var min, max, step float64
	min = dmin
	max = dmax
	step = max - min/5
	v := strings.Split(yrange, ",")
	switch len(v) {
	case 1:
		min = stof(v[0])
	case 2:
		min = stof(v[0])
		max = stof(v[1])
	case 3:
		min = stof(v[0])
		max = stof(v[1])
		step = stof(v[2])
	}
	return min, max, step
}

// gchart draws a chart
func gchart(s string, w, h int, data chart.ChartBox, opts chartOptions) {
	width := float32(w)
	height := float32(h)
	appsize := app.Size(unit.Dp(width), unit.Dp(height))
	apptitle := app.Title("Chart: " + data.Title)

	win := new(app.Window)
	win.Option(apptitle, appsize)

	// Define the colors
	datacolor := giocanvas.ColorLookup(opts.dcolor)
	labelcolor := giocanvas.ColorLookup(opts.labelcolor)
	bgcolor := giocanvas.ColorLookup(opts.bgcolor)

	// Set the chart attributes
	data.Zerobased = opts.zb
	data.Top, data.Bottom = opts.top, opts.bottom
	data.Left, data.Right = opts.left, opts.right

	for {
		switch e := win.Event().(type) {
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), e)
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
			if opts.pie {
				data.Pie(canvas, opts.piesize)
			}
			if opts.lego {
				data.Lego(canvas, opts.dotsize)
			}

			// Draw labels, axes if specified
			data.Color = labelcolor
			if opts.line || opts.bar || opts.scatter || opts.area {
				data.Label(canvas, opts.textsize, opts.xlabel)
				if len(opts.yrange) > 0 {
					yaxmin, yaxmax, yaxstep := yr(opts.yrange, data.Minvalue, data.Maxvalue)
					data.YAxis(canvas, opts.textsize, yaxmin, yaxmax, yaxstep, opts.yaxfmt, opts.showgrid)
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

		case app.DestroyEvent:
			os.Exit(0)
		}
	}
}

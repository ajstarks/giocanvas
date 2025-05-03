// component charts
package main

import (
	"bytes"
	"flag"
	"fmt"
	"image/color"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
	"github.com/ajstarks/giocanvas/chart"
)

func comp(canvas *giocanvas.Canvas) {
	var s1buf, s2buf bytes.Buffer
	begin := 0.0
	end := math.Pi * 4
	incr := 0.1

	// write to data sets, y=sin(x) and y=2*sin(x)
	for x := begin; x <= end; x += incr {
		if x == begin { // add data set labels
			fmt.Fprintln(&s1buf, "# y=sin(x)")
			fmt.Fprintln(&s2buf, "# y=2*sin(x)")
		}
		fmt.Fprintf(&s1buf, fmt.Sprintf("%.2f\t%f\n", x, math.Sin(x)))
		fmt.Fprintf(&s2buf, fmt.Sprintf("%.2f\t%f\n", x, 2*math.Sin(x)))
	}

	// read in data sets
	sr1 := bytes.NewReader(s1buf.Bytes())
	sr2 := bytes.NewReader(s2buf.Bytes())
	sine1, err := chart.DataRead(sr1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	sine2, err := chart.DataRead(sr2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}

	minv, maxv := -2.0, 2.0
	dotsize := 0.4
	frameOpacity := 5.0

	// set attributes for each data set
	sine1.Zerobased = false
	sine2.Zerobased = false

	sine1.MinMax(minv, maxv)
	sine2.MinMax(minv, maxv)

	sine1.Frame(canvas, frameOpacity)
	sine1.Label(canvas, 1.5, 10, "", "gray")
	sine1.YAxis(canvas, 1.5, minv, maxv, 0.5, "%0.2f", true)

	sine1.Color = color.NRGBA{128, 0, 0, 255}
	sine2.Color = color.NRGBA{0, 128, 0, 255}

	// chart the data on the same frame
	sine1.Scatter(canvas, dotsize)
	sine2.Scatter(canvas, dotsize)

	// using the same data sets, make separate charts
	sine1.Left = 10
	sine1.Right = sine1.Left + 30

	sine1.Top = 30
	sine2.Top = 30

	sine1.Bottom = 10
	sine2.Bottom = 10

	dotsize /= 2
	frameOpacity *= 2

	sine1.CTitle(canvas, 2, 2)
	sine1.Frame(canvas, frameOpacity)
	sine1.Scatter(canvas, dotsize)

	offset := 50.0
	sine2.Left = sine1.Left + offset
	sine2.Right = sine1.Right + offset

	sine2.CTitle(canvas, 2, 2)
	sine2.Frame(canvas, frameOpacity)
	sine2.Scatter(canvas, dotsize)

}

func sincos(width, height float32) error {
	w := new(app.Window)
	w.Option(app.Title("composite chart"), app.Size(unit.Dp(width), unit.Dp(height)))
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), e)
			comp(canvas)
			e.Frame(canvas.Context.Ops)
		case app.DestroyEvent:
			os.Exit(0)
		}
	}
}

func main() {
	var cw, ch int
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.Parse()
	go sincos(float32(cw), float32(ch))
	app.Main()
}

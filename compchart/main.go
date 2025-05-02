// component charts
package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
	"github.com/ajstarks/giocanvas/chart"
)

func comp(canvas *giocanvas.Canvas) {
	sr, err := os.Open("sine.d")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	sr2, err := os.Open("sine2.d")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	sine, err := chart.DataRead(sr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(3)
	}
	sine2, err := chart.DataRead(sr2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(4)
	}

	minv, maxv := -2.0, 2.0
	dotsize := 0.4
	frameOpacity := 5.0
	sine2.Zerobased, sine.Zerobased = false, false
	sine.Minvalue, sine.Maxvalue = minv, maxv
	sine2.Minvalue, sine2.Maxvalue = minv, maxv

	sine.Frame(canvas, frameOpacity)
	sine.Label(canvas, 1.5, 10, "", "gray")
	sine.YAxis(canvas, 1.5, minv, maxv, 0.5, "%0.2f", true)

	sine2.Color = color.NRGBA{0, 128, 0, 255}
	sine.Color = color.NRGBA{128, 0, 0, 255}
	sine2.Scatter(canvas, dotsize)
	sine.Scatter(canvas, dotsize)

	sine.Left = 10
	sine.Right = sine.Left + 30
	sine.Top, sine2.Top = 30, 30
	sine.Bottom, sine2.Bottom = 10, 10
	dotsize /= 2
	frameOpacity *= 2

	sine.CTitle(canvas, 2, 2)
	sine.Frame(canvas, frameOpacity)
	sine.Scatter(canvas, dotsize)

	offset := 50.0
	sine2.Left = sine.Left + offset
	sine2.Right = sine.Right + offset

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

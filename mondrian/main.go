// mondrian makes Composition II with Red Blue and Yellow by Piet Mondrian
package main

import (
	"flag"
	"image/color"
	"io"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func mondrian(w *app.Window, width, height float32) error {

	black := color.NRGBA{0, 0, 0, 255}
	white := color.NRGBA{255, 255, 255, 255}
	blue := color.NRGBA{0, 0, 255, 255}
	red := color.NRGBA{255, 0, 0, 255}
	yellow := color.NRGBA{255, 255, 0, 255}

	var third float32 = 100.0 / 3
	var border float32 = 1
	halft := third / 2
	qt := third / 4
	t2 := third * 2
	tq := 100.0 - qt
	t2h := t2 + halft

	for {
		e := w.NextEvent()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})
			canvas.Background(white)
			canvas.CenterRect(halft, halft, third, third, blue)      // lower left blue square
			canvas.CenterRect(t2, t2, t2, t2, red)                   // big red
			canvas.CenterRect(tq, qt, halft, halft, yellow)          // small yellow lower right
			canvas.Line(0, 0, 100, 0, border, black)                 // top border
			canvas.Line(0, 0, 0, 100, border, black)                 // left border
			canvas.Line(100, border/2, 100, 100, border, black)      // right border
			canvas.Line(0, 100, 100, 100, border, black)             // bottom border
			canvas.Line(t2h, halft, t2h+halft, halft, border, black) // top of yellow square
			canvas.Line(third, 100, third, 0, border, black)         //  first column border
			canvas.Line(t2h, 0, t2h, third, border, black)           // left of small right squares
			canvas.Line(0, third, 100, third, border, black)         // top of bottom squares
			canvas.Line(0, t2, third, t2, border, black)             // border between left white squares
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
		w := app.NewWindow(app.Title("mondrian"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := mondrian(w, width, height); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

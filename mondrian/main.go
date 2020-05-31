// mondrian makes Composition II with Red Blue and Yellow by Piet Mondrian
package main

import (
	"image/color"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func main() {

	width := float32(750)
	height := float32(750)

	black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}
	blue := color.RGBA{0, 0, 255, 255}
	red := color.RGBA{255, 0, 0, 255}
	yellow := color.RGBA{255, 255, 0, 255}

	third := float32(100.0) / 3

	halft := third / 2
	qt := third / 4
	t2 := third * 2
	tq := 100.0 - qt
	t2h := t2 + halft

	border := float32(1.0)
	b2 := border * 2

	size := app.Size(unit.Dp(width), unit.Dp(height))
	title := app.Title("Mondrian")

	go func() {
		w := app.NewWindow(title, size)
		for e := range w.Events() {
			if e, ok := e.(system.FrameEvent); ok {
				canvas := giocanvas.NewCanvas(width, height, e.Config, e.Queue, e.Size)

				canvas.Rect(0, 100, 100, 100, white)                // white background
				canvas.CenterRect(halft, halft, third, third, blue) // lower left blue square
				canvas.CenterRect(t2, t2, t2, t2, red)              // big red
				canvas.CenterRect(tq, qt, halft, halft, yellow)     // small yellow lower right
				canvas.VLine(0, 0, 100, b2, black)                  // left border
				canvas.VLine(100, 0, 100, b2, black)                // right border
				canvas.HLine(0, 0, 100, b2, black)                  // top border
				canvas.HLine(0, 100, 100, b2, black)                // bottom border
				canvas.HLine(t2h, halft, t2h+halft, border, black)  // top of yellow square
				canvas.VLine(third, 0, 100, border, black)          //  first column border
				canvas.VLine(t2h, 0, third, border, black)          // left of small right squares
				canvas.HLine(0, third, 100, border, black)          // top of bottom squares
				canvas.VLine(0, t2, third, border, black)           // border between left white squares

				e.Frame(canvas.Context.Ops)

			}
		}
	}()
	app.Main()
}

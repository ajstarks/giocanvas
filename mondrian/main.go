// mondrian makes Composition II with Red Blue and Yellow by Piet Mondrian
package main

import (
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func mondrian(s string, w, h int) {
	width, height := float32(w), float32(h)
	size := app.Size(unit.Px(width), unit.Px(height))
	title := app.Title(s)

	black := color.RGBA{0, 0, 0, 255}
	white := color.RGBA{255, 255, 255, 255}
	blue := color.RGBA{0, 0, 255, 255}
	red := color.RGBA{255, 0, 0, 255}
	yellow := color.RGBA{255, 255, 0, 255}

	var third float32 = 100.0 / 3
	var border float32 = 1
	halft := third / 2
	qt := third / 4
	t2 := third * 2
	tq := 100.0 - qt
	t2h := t2 + halft

	win := app.NewWindow(title, size)
	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, e)

			canvas.Rect(0, 100, 100, 100, white)                     // white background
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
		case key.Event:
			switch e.Name {
			case "Q", key.NameEscape:
				os.Exit(0)
			}

		}
	}
}

func main() {
	go mondrian("Mondrian", 1000, 1000)
	app.Main()
}

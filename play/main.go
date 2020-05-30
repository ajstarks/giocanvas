package main

import (
	"flag"
	"fmt"
	"image/color"
	"math/rand"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func fade(c *giocanvas.Canvas, y, width, size, interval float32, fint uint8, color color.RGBA) {
	for x := float32(0); x <= width; x += width / interval {
		c.CenterRect(x, y, size, size, color)
		c.TextMid(x, y+10, size, fmt.Sprintf("%v", color.A), c.TextColor)
		color.A -= fint
	}
}

func rn(n int) float32 {
	return float32(rand.Intn(n))
}

// randrect randomly  places random sized rectangles
func randrect(c *giocanvas.Canvas, n int) {
	color := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	for i := 0; i < n; i++ {
		x := rn(100)
		y := rn(100)
		w := rn(10)
		h := rn(10)
		if y < 50 || x > 50 {
			continue
		}
		color.R = uint8(rand.Intn(255))
		color.G = uint8(rand.Intn(255))
		color.B = uint8(rand.Intn(255))
		color.A = uint8(rand.Intn(255))
		c.CenterRect(x, y, w, h, color)
	}
}

func main() {
	var w, h int
	flag.IntVar(&w, "width", 1000, "canvas width")
	flag.IntVar(&h, "height", 1000, "canvas height")
	flag.Parse()
	width := float32(w)
	height := float32(h)
	size := app.Size(unit.Dp(width), unit.Dp(height))
	title := app.Title("Gio Canvas")
	tcolor := color.RGBA{128, 0, 0, 255}
	fcolor := color.RGBA{0, 0, 128, 255}
	labelcolor := color.RGBA{0, 0, 0, 255}
	labelsize := float32(2.5)

	go func() {
		w := app.NewWindow(title, size)
		canvas := giocanvas.NewCanvas(width, height)
		tcolor.A = 150
		fcolor.A = 150
		for e := range w.Events() {
			if e, ok := e.(system.FrameEvent); ok {

				canvas.Context.Reset(e.Queue, e.Config, e.Size)

				canvas.TextMid(90, 90, labelsize, "Rectangle", labelcolor)
				canvas.Rect(85, 85, 5, 20, tcolor)
				canvas.CenterRect(90, 65, 10, 10, fcolor)

				canvas.TextMid(90, 50, labelsize, "Image", labelcolor)
				canvas.CenterImage("follow.jpg", 90, 35, 816, 612, 30)

				// Curve
				canvas.TextMid(60, 90, labelsize, "Curve", labelcolor)
				canvas.Curve(45, 65, 75, 85, 75, 65, 0, tcolor)
				canvas.TextMid(45, 60, 1.2, "(45, 65)", fcolor)
				canvas.TextMid(75, 80, 1.2, "(75, 85)", fcolor)
				canvas.TextMid(75, 60, 1.2, "(75, 65)", fcolor)

				// Polygon
				fcolor.A = 150
				canvas.TextMid(60, 50, labelsize, "Polygon", labelcolor)
				xp := []float32{50, 60, 70, 70, 60, 50}
				yp := []float32{50, 40, 50, 25, 30, 25}
				canvas.Polygon(xp, yp, fcolor)

				// Lines
				fcolor.A = 255
				canvas.TextMid(25, 90, labelsize, "Line", labelcolor)
				for y := float32(5); y <= 90; y += 10 {
					canvas.Line(10, 50, 40, y, 0.2, fcolor)
				}

				// Grid
				tcolor.A = 255
				for x := float32(5); x <= 95; x += 5 {
					v := fmt.Sprintf("%v", x)
					canvas.TextMid(x, 2, 1.5, v, color.RGBA{0, 0, 0, 200})
					canvas.TextMid(2, x-0.75, 1.5, v, color.RGBA{0, 0, 0, 200})
				}
				canvas.Grid(0, 0, 100, 100, 0.1, 5, color.RGBA{0, 0, 0, 100})

				e.Frame(canvas.Context.Ops)
			}
		}
	}()
	app.Main()
}

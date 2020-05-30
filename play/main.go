package main

import (
	"flag"
	"fmt"
	"image/color"
	"math/rand"

	"gioui.org/app"
	"gioui.org/f32"
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

func coord(canvas *giocanvas.Canvas, x, y, size float32, color color.RGBA) {
	canvas.Square(x, y, size/2, color)
	canvas.TextMid(x, y+size, size, fmt.Sprintf("(%v,%v)", x, y), color)
}

func main() {
	var w, h int
	var showgrid bool
	flag.IntVar(&w, "width", 1600, "canvas width")
	flag.IntVar(&h, "height", 1000, "canvas height")
	flag.BoolVar(&showgrid, "grid", false, "show grid")
	flag.Parse()
	width := float32(w)
	height := float32(h)
	size := app.Size(unit.Dp(width), unit.Dp(height))
	title := app.Title("Gio Canvas")
	tcolor := color.RGBA{128, 0, 0, 255}
	fcolor := color.RGBA{0, 0, 128, 255}
	bgcolor := color.RGBA{255, 255, 255, 255}
	labelcolor := color.RGBA{0, 0, 0, 255}
	labelsize := float32(2)
	titlesize := labelsize * 2
	subsize := labelsize / 2

	go func() {
		w := app.NewWindow(title, size)
		canvas := giocanvas.NewCanvas(width, height)
		tcolor.A = 150
		fcolor.A = 150
		for e := range w.Events() {
			if e, ok := e.(system.FrameEvent); ok {

				canvas.Context.Reset(e.Queue, e.Config, e.Size)

				// Title
				canvas.CenterRect(50, 50, 100, 100, bgcolor)
				canvas.TextMid(50, 92, titlesize, "Gio Canvas API", labelcolor)

				// Lines
				fcolor.A = 255
				tcolor.A = 255
				lw := float32(0.1)
				canvas.TextMid(25, 80, labelsize, "Line", labelcolor)
				for y := float32(40); y <= 80; y += 5 {
					canvas.Line(25, 60, 40, y, lw, fcolor)
					canvas.Line(10, y, 25, 60, lw, fcolor)
				}

				// Curve
				fcolor.A = 100
				tcolor.A = 100
				start := f32.Point{X: 45, Y: 65}
				c1 := f32.Point{X: 75, Y: 85}
				end := f32.Point{X: 75, Y: 65}
				canvas.TextMid(60, 80, labelsize, "Quadratic Bezier Curve", labelcolor)
				canvas.Curve(start.X, start.Y, c1.X, c1.Y, end.X, end.Y, tcolor)
				coord(canvas, start.X, start.Y, subsize, labelcolor)
				coord(canvas, c1.X, c1.Y, subsize, labelcolor)
				coord(canvas, end.X, end.Y, subsize, labelcolor)

				// Cubic Bezier
				start = f32.Point{X: 45, Y: 40}
				c1 = f32.Point{X: 45, Y: 55}
				c2 := f32.Point{X: 60, Y: 50}
				end = f32.Point{X: 75, Y: 40}

				canvas.TextMid(60, 55, labelsize, "Cubic Bezier Curve", labelcolor)
				canvas.CubeCurve(start.X, start.Y, c1.X, c1.Y, c2.X, c2.Y, end.X, end.Y, fcolor)
				coord(canvas, start.X, start.Y, 1, labelcolor)
				coord(canvas, end.X, end.Y, subsize, labelcolor)
				coord(canvas, c1.X, c1.Y, subsize, labelcolor)
				coord(canvas, c2.X, c2.Y, subsize, labelcolor)

				// Polygon
				canvas.TextMid(60, 30, labelsize, "Polygon", labelcolor)
				xp := []float32{45, 60, 75, 75, 60, 45}
				yp := []float32{25, 20, 25, 5, 10, 5}
				for i := 0; i < len(xp); i++ {
					coord(canvas, xp[i], yp[i], 1, labelcolor)
				}
				canvas.Polygon(xp, yp, fcolor)

				// Rectangles
				canvas.TextMid(90, 80, labelsize, "Rectangle", labelcolor)
				canvas.Rect(85, 75, 5, 20, tcolor)
				canvas.CenterRect(90, 60, 10, 10, fcolor)

				// Image
				canvas.TextMid(90, 45, labelsize, "Image", labelcolor)
				canvas.Image("follow.jpg", 90, 30, 816, 612, 30)

				// Grid
				if showgrid {
					for x := float32(5); x <= 95; x += 5 {
						v := fmt.Sprintf("%v", x)
						canvas.TextMid(x, 2, 1.5, v, color.RGBA{0, 0, 0, 200})
						canvas.TextMid(2, x-0.75, 1.5, v, color.RGBA{0, 0, 0, 200})
					}
					canvas.Grid(0, 0, 100, 100, 0.1, 5, color.RGBA{0, 0, 0, 100})
				}
				e.Frame(canvas.Context.Ops)
			}
		}
	}()
	app.Main()
}

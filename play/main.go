package main

import (
	"flag"
	"fmt"
	"image/color"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func play(appname string, w, h int, showgrid bool) {
	width, height := float32(w), float32(h)
	size := app.Size(unit.Px(width), unit.Px(height))
	title := app.Title(appname)
	tcolor := color.RGBA{128, 0, 0, 255}
	fcolor := color.RGBA{0, 0, 128, 255}
	bgcolor := color.RGBA{255, 255, 255, 255}
	labelcolor := color.RGBA{50, 50, 50, 255}
	labelsize := float32(2)
	titlesize := labelsize * 2
	subsize := labelsize * 0.7
	var colx float32
	var lw float32 = 0.2
	gofont.Register()

	win := app.NewWindow(title, size)
	for e := range win.Events() {
		if e, ok := e.(system.FrameEvent); ok {
			canvas := giocanvas.NewCanvas(width, height, e.Config, e.Queue, e.Size)

			// Title
			canvas.Background(bgcolor)
			canvas.TextMid(50, 92, titlesize, appname, labelcolor)

			colx = 20
			// Lines
			canvas.TextMid(colx, 80, labelsize, "Line", labelcolor)
			canvas.Line(10, 70, colx+5, 65, lw, tcolor)
			canvas.Coord(10, 70, subsize, "P0", labelcolor)
			canvas.Coord(colx+5, 65, subsize, "P1", labelcolor)

			canvas.Line(colx, 70, 35, 75, lw, fcolor)
			canvas.Coord(colx, 70, subsize, "P0", labelcolor)
			canvas.Coord(35, 75, subsize, "P1", labelcolor)

			fcolor.A = 100
			tcolor.A = 100

			// Circle
			canvas.TextMid(colx, 55, labelsize, "Circle", labelcolor)
			canvas.Circle(colx, 45, 5, fcolor)
			canvas.Coord(colx, 45, subsize, "center", labelcolor)

			// Ellipse
			canvas.TextMid(colx, 30, labelsize, "Ellipse", labelcolor)
			ex := (colx / 100) * width
			canvas.AbsEllipse(ex, height*0.85, width*0.05, height*0.10, tcolor)
			canvas.AbsEllipse(ex, height*0.85, width*0.10, height*0.05, fcolor)
			canvas.Coord(colx, 15, subsize, "center", labelcolor)

			// Quadradic Bezier
			start := f32.Point{X: 45, Y: 65}
			c1 := f32.Point{X: 70, Y: 85}
			end := f32.Point{X: 70, Y: 65}
			canvas.TextMid(60, 80, labelsize, "Quadratic Bezier Curve", labelcolor)
			canvas.Curve(start.X, start.Y, c1.X, c1.Y, end.X, end.Y, tcolor)
			canvas.Coord(start.X, start.Y, subsize, "start", labelcolor)
			canvas.Coord(c1.X, c1.Y, subsize, "control", labelcolor)
			canvas.Coord(end.X, end.Y, subsize, "end", labelcolor)

			colx += 40
			// Cubic Bezier
			start = f32.Point{X: 45, Y: 40}
			c1 = f32.Point{X: 45, Y: 55}
			c2 := f32.Point{X: colx, Y: 50}
			end = f32.Point{X: 70, Y: 40}
			canvas.TextMid(colx, 55, labelsize, "Cubic Bezier Curve", labelcolor)
			canvas.CubeCurve(start.X, start.Y, c1.X, c1.Y, c2.X, c2.Y, end.X, end.Y, fcolor)
			canvas.Coord(start.X, start.Y, subsize, "start", labelcolor)
			canvas.Coord(end.X, end.Y, subsize, "end", labelcolor)
			canvas.Coord(c1.X, c1.Y, subsize, "control 1", labelcolor)
			canvas.Coord(c2.X, c2.Y, subsize, "control 2", labelcolor)

			// Polygon
			canvas.TextMid(colx, 30, labelsize, "Polygon", labelcolor)
			xp := []float32{45, 60, 70, 70, 60, 45}
			yp := []float32{25, 20, 25, 5, 10, 5}
			for i := 0; i < len(xp); i++ {
				canvas.Coord(xp[i], yp[i], subsize, fmt.Sprintf("P%d", i), labelcolor)
			}
			canvas.Polygon(xp, yp, fcolor)

			colx += 30
			// Rectangles
			canvas.TextMid(colx, 80, labelsize, "Rectangle", labelcolor)
			canvas.CenterRect(colx, 70, 5, 15, fcolor)
			canvas.Coord(colx, 70, subsize, "center", labelcolor)

			// Square
			canvas.TextMid(colx, 55, labelsize, "Square", labelcolor)
			canvas.Square(colx, 45, 10, tcolor)
			canvas.Coord(colx, 45, subsize, "center", labelcolor)

			// Image
			canvas.TextMid(colx, 30, labelsize, "Image", labelcolor)
			canvas.Image("earth.jpg", colx, 15, int(width*.15), int(width*.15), 100)
			canvas.Coord(colx, 15, subsize, "", color.RGBA{255, 255, 255, 255})

			// Grid
			if showgrid {
				gridcolor := color.RGBA{0, 0, 128, 128}
				var gridsize float32 = 1.2
				for x := float32(5); x <= 95; x += 5 {
					v := fmt.Sprintf("%v", x)
					canvas.TextMid(x, 2, gridsize, v, gridcolor)
					canvas.TextMid(2, x-0.75, gridsize, v, gridcolor)
				}
				canvas.Grid(0, 0, 100, 100, 0.1, 5, gridcolor)
			}
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {
	var w, h int
	var showgrid bool
	flag.IntVar(&w, "width", 1600, "canvas width")
	flag.IntVar(&h, "height", 1000, "canvas height")
	flag.BoolVar(&showgrid, "grid", false, "show grid")
	flag.Parse()
	go play("Gio Canvas API", w, h, showgrid)
	app.Main()
}

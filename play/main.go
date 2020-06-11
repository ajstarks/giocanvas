package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func play(appname string, w, h int, showgrid bool) {
	width, height := float32(w), float32(h)
	size := app.Size(unit.Px(width), unit.Px(height))
	title := app.Title(appname)
	tcolor := color.RGBA{128, 0, 0, 100}
	fcolor := color.RGBA{0, 0, 128, 100}
	bgcolor := color.RGBA{255, 255, 255, 255}
	labelcolor := color.RGBA{50, 50, 50, 255}

	var colx float32
	var lw float32 = 0.2
	var labelsize float32 = 2
	titlesize := labelsize * 2
	subsize := labelsize * 0.7
	subtitle := `A canvas API for Gio applications using high-level objects and a percentage-based coordinate system (https://github.com/ajstarks/giocanvas)`

	win := app.NewWindow(title, size)
	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, e)

			// Title
			canvas.Background(bgcolor)

			colx = 10
			canvas.Text(colx, 92, titlesize, appname, labelcolor)
			canvas.TextWrap(colx+35, 95, titlesize*0.3, 40, subtitle, labelcolor)

			colx = 20
			// Lines
			canvas.TextMid(colx, 80, labelsize, "Line", labelcolor)
			canvas.Line(10, 70, colx+5, 65, lw, tcolor)
			canvas.Coord(10, 70, subsize, "P0", labelcolor)
			canvas.Coord(colx+5, 65, subsize, "P1", labelcolor)

			canvas.Line(colx, 70, 35, 75, lw, fcolor)
			canvas.Coord(colx, 70, subsize, "P0", labelcolor)
			canvas.Coord(35, 75, subsize, "P1", labelcolor)

			// Circle
			canvas.TextMid(colx, 55, labelsize, "Circle", labelcolor)
			canvas.Circle(colx, 45, 5, fcolor)
			canvas.Coord(colx, 45, subsize, "center", labelcolor)

			// Ellipse
			canvas.TextMid(colx, 30, labelsize, "Ellipse", labelcolor)
			canvas.Ellipse(colx, 15, 5, 10, tcolor)
			canvas.Ellipse(colx, 15, 10, 5, fcolor)
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
		case key.Event:
			switch e.Name {
			case "Q", key.NameEscape:
				os.Exit(0)
			}
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

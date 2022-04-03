// play demos the giocanvas API
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func getImage(name string) image.Image {
	r, err := os.Open(name)
	if err != nil {
		return nil
	}
	im, _, err := image.Decode(r)
	if err != nil {
		return nil
	}
	return im
}

func play(w *app.Window, width, height float32, showgrid bool) error {
	tcolor := color.NRGBA{128, 0, 0, 150}
	fcolor := color.NRGBA{0, 0, 128, 150}
	bgcolor := color.NRGBA{255, 255, 255, 255}
	labelcolor := color.NRGBA{50, 50, 50, 255}

	var colx float32
	var lw float32 = 0.2
	var labelsize float32 = 2
	titlesize := labelsize * 2
	subsize := labelsize * 0.7
	subtitle := `A canvas API for Gio applications using high-level objects and a percentage-based coordinate system (https://github.com/ajstarks/giocanvas)`
	logoimg := getImage("logo.png")
	const pi = 3.14159265
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, system.FrameEvent{})

			// Title
			canvas.Background(bgcolor)
			canvas.Img(logoimg, 5, 95, 400, 400, 20)

			colx = 20
			canvas.TextMid(colx, 92, titlesize, "Canvas API", labelcolor)
			canvas.TextWrap(colx+15, 95, titlesize*0.3, 50, subtitle, labelcolor)

			// Lines
			canvas.TextMid(colx, 80, labelsize, "Line", labelcolor)
			canvas.Line(10, 70, colx+5, 65, lw, tcolor)
			canvas.Coord(10, 70, subsize, "P0", labelcolor)
			canvas.Coord(colx+5, 65, subsize, "P1", labelcolor)

			canvas.Line(colx, 70, 35, 75, lw, fcolor)
			canvas.Coord(colx, 70, subsize, "P0", labelcolor)
			canvas.Coord(35, 75, subsize, "P1", labelcolor)

			// Circle
			cx1 := colx - 10
			cx2 := colx + 10
			canvas.TextMid(cx1, 55, labelsize, "Circle", labelcolor)
			canvas.Circle(cx1, 45, 5, fcolor)
			canvas.Coord(cx1, 45, subsize, "center", labelcolor)

			// Arc
			canvas.TextMid(cx2, 55, labelsize, "Arc", labelcolor)
			canvas.Arc(cx2, 45, 5, 0, 3*pi/4, tcolor)
			canvas.Coord(cx2, 45, subsize, "center", labelcolor)

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
			canvas.StrokedCurve(start.X, start.Y, c1.X, c1.Y, end.X, end.Y, lw, tcolor)
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
			canvas.StrokedCubeCurve(start.X, start.Y, c1.X, c1.Y, c2.X, c2.Y, end.X, end.Y, lw, fcolor)
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
			canvas.Image("earth.jpg", colx, 15, 1000, 1000, 15)
			canvas.Coord(colx, 15, subsize, "", color.NRGBA{255, 255, 255, 255})

			// Grid
			if showgrid {
				gridcolor := color.NRGBA{0, 0, 128, 50}
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
	var cw, ch int
	var showgrid bool
	flag.IntVar(&cw, "width", 1600, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.BoolVar(&showgrid, "grid", false, "show grid")
	flag.Parse()
	width := float32(cw)
	height := float32(ch)

	go func() {
		w := app.NewWindow(app.Title("Canvas API"), app.Size(unit.Px(width), unit.Px(height)))
		if err := play(w, width, height, showgrid); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()

}

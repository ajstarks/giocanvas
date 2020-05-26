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
	flag.IntVar(&w, "width", 1200, "canvas width")
	flag.IntVar(&h, "height", 900, "canvas height")
	flag.Parse()
	width := float32(w)
	height := float32(h)
	size := app.Size(unit.Dp(width), unit.Dp(height))
	title := app.Title("Gio Canvas")
	tcolor := color.RGBA{128, 0, 0, 255}
	fcolor := color.RGBA{0, 0, 128, 255}

	go func() {
		w := app.NewWindow(title, size)
		canvas := giocanvas.NewCanvas(width, height)
		for e := range w.Events() {
			if e, ok := e.(system.FrameEvent); ok {
				canvas.Context.Reset(e.Queue, e.Config, e.Size)
				canvas.Text(60, 45, 3, "Follow your dreams", tcolor)
				canvas.Text(10, 45, 3, "Random", tcolor)
				canvas.CenterImage("follow.jpg", 75, 75, 816, 612, 75)
				fade(canvas, 25, 100, 2.0, 20, 10, fcolor)
				randrect(canvas, 200)
				for x := float32(0); x <= 100; x += 10 {
					v := fmt.Sprintf("%v", x)
					canvas.Text(x, 20, 1.5, v, tcolor)
					canvas.TextMid(x, 15, 1.5, v, tcolor)
					canvas.TextEnd(x, 10, 1.5, v, tcolor)
				}
				e.Frame(canvas.Context.Ops)
			}
		}
	}()
	app.Main()
}

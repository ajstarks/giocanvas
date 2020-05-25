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

// randrect randomly  places random sized rectangles
func randrect(c *giocanvas.Canvas, w, h float32, n int) {
	color := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	for i := 0; i < n; i++ {
		x := rand.Float32() * w
		y := rand.Float32() * h
		w := rand.Float32() * 100
		h := rand.Float32() * 100
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
	gcolor := color.RGBA{0, 0, 0, 128}
	fcolor := color.RGBA{0, 0, 128, 255}

	go func() {
		w := app.NewWindow(title, size)
		canvas := giocanvas.NewCanvas(width, height)
		for e := range w.Events() {
			if e, ok := e.(system.FrameEvent); ok {
				canvas.Context.Reset(e.Queue, e.Config, e.Size)
				canvas.Text(width*0.60, height*0.55, 40, "Follow your dreams", tcolor)
				canvas.Text(width*0.10, height*0.55, 40, "Random", tcolor)
				canvas.CenterImage("follow.jpg", width*0.75, height*0.25, 816, 612, 75)
				randrect(canvas, width/2, height/2, 300)
				canvas.Grid(width, height, 2, 20, gcolor)
				fade(canvas, height*0.75, width, 20, 20, 10, fcolor)

				for x := float32(0); x <= canvas.Width; x += canvas.Width / 20 {
					v := fmt.Sprintf("%v", x)
					canvas.Text(x, height*0.80, 15, v, tcolor)
					canvas.TextMid(x, height*0.85, 15, v, tcolor)
					canvas.TextEnd(x, height*0.90, 15, v, tcolor)
				}
				e.Frame(canvas.Context.Ops)
			}
		}
	}()
	app.Main()
}

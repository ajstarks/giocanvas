// hello is the giocanvas hello, world
package main

import (
	"flag"
	"image"
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func getimage(s string) (image.Image, error) {
	i, err := os.Open(s)
	if err != nil {
		return nil, err
	}
	im, _, err := image.Decode(i)
	if err != nil {
		return nil, err
	}
	i.Close()
	return im, nil
}

func images(title string, width, height float32) {
	defer os.Exit(0)
	win := app.NewWindow(app.Title(title), app.Size(unit.Px(width), unit.Px(height)))
	canvas := giocanvas.NewCanvas(width, height, system.FrameEvent{})
	im, err := getimage("earth.jpg")
	if err != nil {
		return
	}
	bgcolor := color.NRGBA{0, 0, 0, 255}
	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			var x, y, scale float32
			scale = 2.0
			canvas.Background(bgcolor)
			canvas.Grid(0, 0, 100, 100, 0.1, 10, color.NRGBA{128, 128, 128, 255})
			for x = 20; x <= 80; x += 10 {
				y = x
				canvas.Img(im, x, y, 1000, 1000, scale)
				scale += 2.0
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
	flag.IntVar(&w, "width", 1000, "canvas width")
	flag.IntVar(&h, "height", 1000, "canvas height")
	flag.Parse()
	go images("images", float32(w), float32(h))
	app.Main()
}

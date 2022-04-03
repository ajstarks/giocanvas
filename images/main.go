// hello is the giocanvas hello, world
package main

import (
	"flag"
	"image"
	"image/color"
	"io"
	"os"

	"gioui.org/app"
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

func images(w *app.Window, width, height float32) error {
	im, err := getimage("earth.jpg")
	if err != nil {
		return err
	}
	bgcolor := color.NRGBA{0, 0, 0, 255}
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, system.FrameEvent{})
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
		}
	}
}

func main() {
	var cw, ch int
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.Parse()
	width, height := float32(cw), float32(ch)
	go func() {
		w := app.NewWindow(app.Title("images"), app.Size(unit.Px(width), unit.Px(height)))
		if err := images(w, width, height); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

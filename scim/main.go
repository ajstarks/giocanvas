// scalable image
package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
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

func scimage(w *app.Window, filename string, width, height float32) error {
	im, err := getimage(filename)
	if err != nil {
		return err
	}
	imw := im.Bounds().Dx()
	imh := im.Bounds().Dy()
	bgcolor := color.NRGBA{0, 0, 0, 255}
	fgcolor := color.NRGBA{255, 255, 255, 255}
	gridcolor := fgcolor
	gridcolor.A = 75
	for {
		e := w.NextEvent()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), system.FrameEvent{})
			scale := (float32(e.Size.X) / float32(imw)) * 100
			canvas.Background(bgcolor)
			canvas.Img(im, 50, 50, imw, imh, scale)
			canvas.CText(50, 50, 5, "Scaled Image", fgcolor)
			canvas.Grid(0, 0, 100, 100, 0.1, 10, gridcolor)
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {
	var cw, ch int
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.Parse()
	filenames := flag.Args()
	var imfile string
	if len(filenames) == 0 {
		imfile = "earth.jpg"
	} else {
		imfile = filenames[0]
	}
	width, height := float32(cw), float32(ch)
	go func() {
		w := app.NewWindow(app.Title("Scaled Image"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := scimage(w, imfile, width, height); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create the window: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

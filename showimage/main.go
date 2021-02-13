// showimage shows an image
package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

// imageinfo opens an image file, returning an image.Image, with dimensions
func imageinfo(imagefile string, w, h int) (image.Image, int, int, error) {
	f, err := os.Open(imagefile)
	if err != nil {
		return nil, 0, 0, err
	}
	im, _, err := image.Decode(f)
	if err != nil {
		return nil, 0, 0, err
	}
	if w == 0 {
		w = im.Bounds().Dx()
	}
	if h == 0 {
		h = im.Bounds().Dy()
	}
	f.Close()
	return im, w, h, nil
}

// showimage shows an image, centered on the canvas at the specified scale and size
func showimage(title string, im image.Image, width, height int, scale float64) {
	sw, sh, sc := float32(width), float32(height), float32(scale)
	if sc != 100 {
		sw *= sc / 100
		sh *= sc / 100
	}
	win := app.NewWindow(app.Title(title), app.Size(unit.Px(sw), unit.Px(sh)))
	for e := range win.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			return
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(sw, sh, system.FrameEvent{})
			canvas.Img(im, 50, 50, width, height, sc)
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
	var (
		w, h  int
		scale float64
		err   error
		im    image.Image
	)
	flag.IntVar(&w, "width", 0, "canvas width")
	flag.IntVar(&h, "height", 0, "canvas height")
	flag.Float64Var(&scale, "scale", 100, "scale (0-100)")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "specify an image file (JPEG, PNG, or GIF)")
		os.Exit(1)
	}
	for _, imagefile := range args {
		im, w, h, err = imageinfo(imagefile, w, h)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		go showimage(imagefile, im, w, h, scale)
	}
	app.Main()
}

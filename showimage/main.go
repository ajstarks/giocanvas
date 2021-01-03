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
			os.Exit(0)
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
	var w, h int
	var scale float64
	flag.IntVar(&w, "width", 0, "canvas width")
	flag.IntVar(&h, "height", 0, "canvas height")
	flag.Float64Var(&scale, "scale", 100, "scale (0-100)")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "specify an image file")
		os.Exit(1)
	}
	imagefile := args[0]
	f, ferr := os.Open(imagefile)
	if ferr != nil {
		fmt.Fprintf(os.Stderr, "%v\n", ferr)
		os.Exit(2)
	}
	im, _, err := image.Decode(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagefile, err)
		os.Exit(3)
	}
	imw := im.Bounds().Dx()
	imh := im.Bounds().Dy()
	if w == 0 {
		w = imw
	}
	if h == 0 {
		h = imh
	}
	f.Close()
	go showimage(imagefile, im, w, h, scale)
	app.Main()
}

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
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

// imageinfo opens an image file, returning an image.Image, with dimensions
func imageinfo(imagefile string) (image.Image, int, int, error) {
	f, err := os.Open(imagefile)
	if err != nil {
		return nil, 0, 0, err
	}
	im, _, err := image.Decode(f)
	if err != nil {
		return nil, 0, 0, err
	}
	f.Close()
	return im, im.Bounds().Dx(), im.Bounds().Dy(), nil
}

// showimage shows an image, centered on the canvas at the specified scale and size
func showimage(w *app.Window, im image.Image, width, height int, scale float32) error {
	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})
			scale = (float32(e.Size.X) / float32(width)) * 100
			canvas.Img(im, 50, 50, width, height, scale)
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {
	var (
		width, height int
		scale         float64
		err           error
		im            image.Image
	)
	flag.Float64Var(&scale, "scale", 100, "scale (0-100)")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "specify an image file (JPEG, PNG, or GIF)")
		os.Exit(1)
	}
	for _, imagefile := range args {
		im, width, height, err = imageinfo(imagefile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}

		sc, sw, sh := float32(scale), float32(width), float32(height)
		if scale != 100 {
			sw *= sc / 100
			sh *= sc / 100
		}
		w := &app.Window{}
		w.Option(app.Title(imagefile), app.Size(unit.Dp(sw), unit.Dp(sh)))
		if err := showimage(w, im, width, height, sc); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}

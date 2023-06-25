// transform tests affine transforms
package main

import (
	"flag"
	"image/color"
	"io"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	gc "github.com/ajstarks/giocanvas"
)

func transforms(w *app.Window, width, height float32) error {
	var midx, rectw, recth, recty, ts, ts2 float32
	midx = 50
	rectw = 40
	recth = rectw / 4
	ts = 5
	ts2 = ts / 3
	textcolor := color.NRGBA{0, 0, 0, 255}
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			canvas := gc.NewCanvas(float32(e.Size.X), float32(e.Size.Y), system.FrameEvent{})
			recty = 90
			canvas.CenterRect(midx, recty, rectw, recth, color.NRGBA{128, 128, 128, 128})
			canvas.TextMid(midx, recty-ts2, ts, "Reference", textcolor)

			recty = 70
			stack := canvas.Scale(midx, recty, 2)
			canvas.CenterRect(midx, recty, rectw, recth, color.NRGBA{0, 0, 128, 128})
			canvas.TextMid(midx, recty-ts2, ts, "scale", textcolor)
			gc.EndTransform(stack)

			recty = 50
			stack = canvas.Shear(midx, midx, math.Pi/4, 0)
			canvas.CenterRect(midx, recty, rectw, recth, color.NRGBA{128, 0, 0, 128})
			canvas.TextMid(midx, recty-ts2, ts, "shear", textcolor)
			gc.EndTransform(stack)

			stack = canvas.Translate(20, 85)
			canvas.CenterRect(midx, recty, rectw, recth, color.NRGBA{0, 128, 0, 128})
			canvas.TextMid(midx, recty-ts2, ts, "translate", textcolor)
			gc.EndTransform(stack)

			recty = 20
			stack = canvas.Rotate(midx, recty, math.Pi/4)
			canvas.CenterRect(midx, recty, rectw, recth, color.NRGBA{255, 50, 0, 200})
			canvas.TextMid(midx, recty-ts2, ts, "rotate", textcolor)
			gc.EndTransform(stack)
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {
	var cw, ch int
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.Parse()
	width := float32(cw)
	height := float32(ch)

	go func() {
		w := app.NewWindow(app.Title("transform"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := transforms(w, width, height); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()

}

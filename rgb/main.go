// rgb shows RGB values
package main

import (
	"flag"
	"image/color"
	"io"
	"os"

	"gioui.org/app"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

func main() {
	var cw, ch int
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.Parse()

	width := float32(cw)
	height := float32(ch)

	go func() {
		w := &app.Window{}
		w.Option(app.Title("rgb"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := rgb(w); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()

}

func rgb(w *app.Window) error {
	colortab := []string{
		"orange",
		"rgb(100)",
		"rgb(100,50)",
		"rgb(100,50,2)",
		"rgb(100,50,2,100)",
		"hsv(0,70,50)",
		"hsv(0,70,50,50)",
		"#aa",
		"#aabb",
		"#aabbcc",
		"#aabbcc64",
		"rgb()",
		"hsv()",
		"#",
		"#error",
		"nonsense",
	}
	var x, y float32
	for {
		ev := w.Event()
		switch e := ev.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})
			x, y = 50, 95
			for _, c := range colortab {
				canvas.EText(x-10, y, 3, c, color.NRGBA{0, 0, 0, 255})
				canvas.Circle(x, y+1, 2, giocanvas.ColorLookup(c))
				y -= 6
			}
			e.Frame(canvas.Context.Ops)
		}
	}
}

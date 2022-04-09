// rgb shows RGB values
package main

import (
	"flag"
	"image/color"
	"io"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
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
		w := app.NewWindow(app.Title("rgb"), app.Size(unit.Px(width), unit.Px(height)))
		if err := rgb(w, width, height); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

func rgb(w *app.Window, width, height float32) error {
	colortab := []string{
		"orange",
		"rgb(100)",
		"rgb(100,50)",
		"rgb(100,50,2)",
		"rgb(100,50,2,100)",
		"#aa",
		"#aabb",
		"#aabbcc",
		"#aabbcc64",
		"rgb()",
		"#",
		"#error",
		"nonsense",
	}
	var x, y float32
	for {
		ev := <-w.Events()
		switch e := ev.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(width, height, system.FrameEvent{})
			x, y = 50, 90
			for _, c := range colortab {
				canvas.EText(x-10, y, 3, c, color.NRGBA{0, 0, 0, 255})
				canvas.Circle(x, y+1, 2, giocanvas.ColorLookup(c))
				y -= 7
			}
			e.Frame(canvas.Context.Ops)
		}
	}
}

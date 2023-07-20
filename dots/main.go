// dots: draw with dots
package main

import (
	"flag"
	"io"
	"os"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

type coord struct {
	X float32
	Y float32
}

func main() {
	var cw, ch, nc int
	var bgcolor string
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.IntVar(&nc, "nc", 100, "number of dots")
	flag.StringVar(&bgcolor, "bgcolor", "black", "background color")
	flag.Parse()

	// kick off the application
	go func() {
		w := app.NewWindow(app.Title("dots"), app.Size(unit.Dp(cw), unit.Dp(ch)))
		if err := dots(w, nc, bgcolor); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

// pctcoord converts device coordinates to canvas percents
func pctcoord(x, y, width, height float32) (float32, float32) {
	return 100 * (x / width), 100 - (100 * (y / height))
}

var pressed bool
var colorindex int
var coordindex int
var dotsize float32 = 0.5

func kbpointer(q event.Queue, width, height float32, coords []coord) {
	for _, ev := range q.Events(pressed) {
		// keyboard events
		if k, ok := ev.(key.Event); ok {
			switch k.State {
			case key.Press:
				switch k.Name {
				case key.NameUpArrow, key.NameRightArrow:
					dotsize += 0.1
				case key.NameDownArrow, key.NameLeftArrow:
					dotsize -= 0.1
				case key.NameEscape, "Q":
					os.Exit(0)
				}
			}
		}
		// pointer events
		if p, ok := ev.(pointer.Event); ok {
			switch p.Type {
			case pointer.Drag:
				coords[coordindex].X, coords[coordindex].Y = pctcoord(p.Position.X, p.Position.Y, width, height)
				coordindex++
				if coordindex == len(coords) {
					coordindex = 0
				}
			case pointer.Press:
				switch p.Buttons {
				case pointer.ButtonSecondary:
					dotsize += 0.1
				case pointer.ButtonTertiary:
					dotsize -= 0.1
				}
				pressed = true
			}
		}
	}
}

func dots(w *app.Window, nc int, bgcolor string) error {
	palette := []string{
		"#aaaaaaaa",
		"#aa0000aa",
		"#00aa00aa",
		"#0000aaaa",
		"#ffd821aa",
		"#234ad5aa",
		"#ffad5e00",
		"#000000aa",
	}
	np := len(palette)
	coordinates := make([]coord, nc)
	coordinates[0].X, coordinates[0].Y = 50, 50
	bg := giocanvas.ColorLookup(bgcolor)
	for {
		ev := <-w.Events()
		switch e := ev.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			w, h := float32(e.Size.X), float32(e.Size.Y)
			canvas := giocanvas.NewCanvas(w, h, system.FrameEvent{})
			key.InputOp{Tag: pressed}.Add(canvas.Context.Ops)
			pointer.InputOp{
				Tag:   pressed,
				Grab:  false,
				Types: pointer.Press | pointer.Move | pointer.Drag}.Add(canvas.Context.Ops)
			if dotsize < 0.1 {
				dotsize = 0.1
			}
			if dotsize > 5 {
				dotsize = 5
			}
			cs := dotsize
			ci := 0
			canvas.Background(bg)
			for i := 0; i < nc; i++ {
				x1, y1 := coordinates[i].X, coordinates[i].Y
				c := giocanvas.ColorLookup(palette[ci%np])
				canvas.Circle(x1, y1, cs, c)
				cs += 0.01
				ci++
			}
			kbpointer(e.Queue, w, h, coordinates)
			e.Frame(canvas.Context.Ops)
		}
	}
}

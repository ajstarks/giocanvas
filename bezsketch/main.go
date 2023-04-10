// bezsketch: sketch quadradic Bezier curves
package main

import (
	"flag"
	"image/color"
	"io"
	"os"
	"strconv"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

type config struct {
	precision                                                          int
	width, height, textsize, coordsize, curvesize                      float32
	curvecolor, bgcolor, textcolor, begincolor, controlcolor, endcolor color.NRGBA
}

func main() {
	var cw, ch int
	var ts, cs, cus float64
	var curvecolor, bgcolor, textcolor, begincolor, endcolor, controlcolor string
	var cfg config

	// set up command flags
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.IntVar(&cfg.precision, "precision", 0, "coordinate precision")
	flag.StringVar(&bgcolor, "bgcolor", "white", "background color")
	flag.StringVar(&textcolor, "textcolor", "black", "text color")
	flag.StringVar(&begincolor, "begincolor", "green", "begin coordinate color")
	flag.StringVar(&endcolor, "endcolor", "red", "end coordinate color")
	flag.StringVar(&controlcolor, "controlcolor", "gray", "control coordinate color")
	flag.StringVar(&curvecolor, "curvecolor", "#aaaaaaaa", "curve color")
	flag.Float64Var(&ts, "tsize", 2.5, "text size")
	flag.Float64Var(&cs, "csize", 1.25, "coordinate size")
	flag.Float64Var(&cus, "curvesize", 1.0, "curve size")
	flag.Parse()

	cfg.width = float32(cw)
	cfg.height = float32(ch)
	cfg.textsize = float32(ts)
	cfg.coordsize = float32(cs)
	cfg.curvesize = float32(cus)
	cfg.bgcolor = giocanvas.ColorLookup(bgcolor)
	cfg.textcolor = giocanvas.ColorLookup(textcolor)
	cfg.curvecolor = giocanvas.ColorLookup(curvecolor)
	cfg.begincolor = giocanvas.ColorLookup(begincolor)
	cfg.endcolor = giocanvas.ColorLookup(endcolor)
	cfg.controlcolor = giocanvas.ColorLookup(controlcolor)

	go func() {
		w := app.NewWindow(app.Title("bezsketch"), app.Size(unit.Dp(cfg.width), unit.Dp(cfg.height)))
		if err := bezsketch(w, cfg); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

var pressed, dogrid bool
var mouseX, mouseY float32
var bx, by, ex, ey, cx, cy float32

// pctcoord converts device coordinates to canvas percents
func pctcoord(x, y, width, height float32) (float32, float32) {
	return 100 * (x / width), 100 - (100 * (y / height))
}

// ftoa converts float to string, with leading space
func ftoa(x float32, prec int) string {
	return " " + strconv.FormatFloat(float64(x), 'f', prec, 32)
}

// textcoord displays a labelled coordinate
func textcoord(canvas *giocanvas.Canvas, x, y float32, color color.NRGBA, cfg config) {
	size := cfg.coordsize
	ts := cfg.textsize
	prec := cfg.precision
	canvas.Circle(x, y, size, color)
	canvas.Circle(x, y, size/2, cfg.bgcolor)
	canvas.TextMid(x, y+ts, ts, ftoa(x, prec)+","+ftoa(y, prec), cfg.textcolor)
}

// pctPointerPos records and processes the pointer events in percent coordinates
func pctPointerPos(q event.Queue, cfg config) {
	width, height := cfg.width, cfg.height
	prec := cfg.precision
	for _, ev := range q.Events(pressed) {
		if k, ok := ev.(key.Event); ok {

			switch k.State {
			case key.Press:
				switch k.Name {
				case "G":
					dogrid = !dogrid
				case "C":
					io.WriteString(os.Stdout, "curve"+ftoa(bx, prec)+ftoa(by, prec)+ftoa(cx, prec)+ftoa(cy, prec)+ftoa(ex, prec)+ftoa(ey, prec)+"\n")
				case key.NameRightArrow:
					switch k.Modifiers {
					case 0:
						bx += 1
					case 4:
						ex += 1
					}
				case key.NameLeftArrow:
					switch k.Modifiers {
					case 0:
						bx -= 1
					case 4:
						ex -= 1
					}
				case key.NameUpArrow:
					switch k.Modifiers {
					case 0:
						by += 1
					case 4:
						ey += 1
					}
				case key.NameDownArrow:
					switch k.Modifiers {
					case 0:
						by -= 1
					case 4:
						ey -= 1
					}
				case key.NameEscape, "Q":
					os.Exit(0)
				}
			}
		}
		if p, ok := ev.(pointer.Event); ok {
			switch p.Type {
			case pointer.Move:
				mouseX, mouseY = pctcoord(p.Position.X, p.Position.Y, width, height)
			case pointer.Press:
				switch p.Buttons {
				case pointer.ButtonPrimary:
					bx, by = pctcoord(p.Position.X, p.Position.Y, width, height)
				case pointer.ButtonSecondary:
					ex, ey = pctcoord(p.Position.X, p.Position.Y, width, height)
				case pointer.ButtonTertiary:
					io.WriteString(os.Stdout, "curve"+ftoa(bx, prec)+ftoa(by, prec)+ftoa(cx, prec)+ftoa(cy, prec)+ftoa(ex, prec)+ftoa(ey, prec)+"\n")
				}
				pressed = true
			}
		}
	}
}

// bezsketch sketches quadratic bezier curves: left pointer press defines the begin point,
// right pointer press defines the end point, middle pointer press shows the curve spec,
// pointer move defines the control point
func bezsketch(w *app.Window, cfg config) error {
	bx, by = 25.0, 50.0
	ex, ey = 75.0, 50.0
	cx, cy = 10, 10
	begincolor, endcolor, controlcolor := cfg.begincolor, cfg.endcolor, cfg.controlcolor

	// event loop
	for {
		ev := <-w.Events()
		switch e := ev.(type) {
		// return an error on close
		case system.DestroyEvent:
			return e.Err

		// for each frame: register press and move events, draw coordinates, and the curve,
		// track the pointer position for the control point, show curve spec on middle click
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(cfg.width, cfg.height, system.FrameEvent{})
			key.InputOp{Tag: pressed, Hint: key.HintAny, Keys: "g|G"}.Add(canvas.Context.Ops)
			pointer.InputOp{Tag: pressed, Grab: false, Types: pointer.Press | pointer.Move}.Add(canvas.Context.Ops)
			canvas.Background(cfg.bgcolor)
			if dogrid {
				canvas.Grid(0, 0, 100, 100, 0.1, 5, cfg.textcolor)
			}
			textcoord(canvas, bx, by, begincolor, cfg)
			textcoord(canvas, ex, ey, endcolor, cfg)
			textcoord(canvas, cx, cy, controlcolor, cfg)
			canvas.QuadStrokedCurve(bx, by, cx, cy, ex, ey, cfg.curvesize, cfg.curvecolor)
			pctPointerPos(e.Queue, cfg)
			cx, cy = mouseX, mouseY
			e.Frame(canvas.Context.Ops)
		}
	}
}

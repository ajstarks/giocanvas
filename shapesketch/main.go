// shapesketch: sketch quadradic Bezier curves
package main

import (
	"flag"
	"image/color"
	"io"
	"math"
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

// configuration options
type config struct {
	precision                                                                     int
	width, height, textsize, coordsize, linesize, shapesize, stepsize             float32
	linecolor, bgcolor, textcolor, begincolor, shapecolor, currentcolor, endcolor color.NRGBA
}

func main() {
	var cw, ch int
	var ts, cs, ls, ss float64
	var shapecolor, bgcolor, textcolor, begincolor, endcolor, currentcolor string
	var cfg config

	// set up command flags
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.IntVar(&cfg.precision, "precision", 0, "coordinate precision")
	flag.StringVar(&bgcolor, "bgcolor", "white", "background color")
	flag.StringVar(&textcolor, "textcolor", "black", "text color")
	flag.StringVar(&begincolor, "begincolor", "green", "begin coordinate color")
	flag.StringVar(&endcolor, "endcolor", "red", "end coordinate color")
	flag.StringVar(&currentcolor, "currentcolor", "gray", "current coordinate color")
	flag.StringVar(&shapecolor, "shapecolor", "#22222255", "curve color")
	flag.Float64Var(&ts, "tsize", 2.5, "text size")
	flag.Float64Var(&cs, "csize", 1.25, "coordinate size")
	flag.Float64Var(&ls, "lsize", 1.0, "line size")
	flag.Float64Var(&ss, "ssize", 0.5, "step size")
	flag.Parse()

	cfg.width = float32(cw)
	cfg.height = float32(ch)
	cfg.textsize = float32(ts)
	cfg.coordsize = float32(cs)
	cfg.linesize = float32(ls)
	cfg.stepsize = float32(ss)
	cfg.bgcolor = giocanvas.ColorLookup(bgcolor)
	cfg.textcolor = giocanvas.ColorLookup(textcolor)
	cfg.linecolor = giocanvas.ColorLookup(shapecolor)
	cfg.begincolor = giocanvas.ColorLookup(begincolor)
	cfg.endcolor = giocanvas.ColorLookup(endcolor)
	cfg.shapecolor = giocanvas.ColorLookup(shapecolor)

	// kick off the application
	go func() {
		w := app.NewWindow(app.Title("shapesketch"), app.Size(unit.Dp(cfg.width), unit.Dp(cfg.height)))
		if err := shapesketch(w, cfg); err != nil {
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
var shape string

// pctcoord converts device coordinates to canvas percents
func pctcoord(x, y, width, height float32) (float32, float32) {
	return 100 * (x / width), 100 - (100 * (y / height))
}

// ftoa converts float to string, with leading space
func ftoa(x float32, prec int) string {
	return " " + strconv.FormatFloat(float64(x), 'f', prec, 32)
}

// grid conditionalls shows a grid
func grid(canvas *giocanvas.Canvas, interval float32, color color.NRGBA) {
	if dogrid {
		canvas.Grid(0, 0, 100, 100, 0.1, interval, color)
	}
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

func puts(s string) {
	io.WriteString(os.Stdout, s)
}

// deckspec displays the decksh curve specification
func deckspec(prec int) {
	switch shape {
	case "arc":
		r, a1, a2 := arcElements()
		d1 := rad2deg(float32(a1))
		d2 := rad2deg(float32(a2))
		puts("arc" + ftoa(bx, prec) + ftoa(by, prec) + ftoa(r, prec) + ftoa(r, prec) + ftoa(d1, prec) + ftoa(d2, prec) + "\n")
	case "bezier":
		puts("curve" + ftoa(bx, prec) + ftoa(by, prec) + ftoa(cx, prec) + ftoa(cy, prec) + ftoa(ex, prec) + ftoa(ey, prec) + "\n")
	case "line":
		puts("line" + ftoa(bx, prec) + ftoa(by, prec) + ftoa(ex, prec) + ftoa(ey, prec) + "\n")
	case "rect", "ellipse":
		dx := float64(cx-bx) * 2
		dy := float64(cy-by) * 2
		rw, rh := float32(math.Abs(dx)), float32(math.Abs(dy))
		puts(shape + ftoa(bx, prec) + ftoa(by, prec) + ftoa(rw, prec) + ftoa(rh, prec) + "\n")
	case "square":
		dx := float64(cx - bx)
		sw := float32(math.Abs(dx) * 2)
		puts(shape + ftoa(bx, prec) + ftoa(by, prec) + ftoa(sw, prec) + "\n")
	case "circle":
		puts("circle" + ftoa(bx, prec) + ftoa(by, prec) + ftoa(dist(bx, by, cx, cy), prec) + "\n")
	}
}

// dist computes the distance between (x1, y1) and (x2, y2)
func dist(x1, y1, x2, y2 float32) float32 {
	x := float64(x2 - x1)
	y := float64(y2 - y1)
	return float32(math.Sqrt(x*x + y*y))
}

func rad2deg(r float32) float32 {
	d := float32((180 / math.Pi) * r)
	if d < 0 {
		return 360 + d
	}
	return d
}

func arcElements() (float32, float64, float64) {
	r := dist(bx, by, cx, cy)
	dx1 := float64(ex - bx)
	dy1 := float64(ey - by)
	dx2 := float64(cx - bx)
	dy2 := float64(cy - by)
	a1 := math.Atan2(dy1, dx1)
	a2 := math.Atan2(dy2, dx2)
	return r, a1, a2
}

// pct returns the percentage of its input
func pct(p float32, m float32) float32 {
	return ((p / 100.0) * m)
}

// dimen returns canvas dimensions from percentages
// (converting from x increasing left-right, y increasing top-bottom)
func dimen(xp, yp, w, h float32) (float32, float32) {
	return pct(xp, w), pct(100-yp, h)
}

func arc(canvas *giocanvas.Canvas, size float32, fillcolor color.NRGBA) {
}

// kbpointer processes the keyboard events and pointer events in percent coordinates
func kbpointer(q event.Queue, cfg config) {
	width, height := cfg.width, cfg.height
	prec := cfg.precision
	stepsize := cfg.stepsize
	for _, ev := range q.Events(pressed) {
		// keyboard events
		if k, ok := ev.(key.Event); ok {
			switch k.State {
			case key.Press:
				switch k.Name {
				case "G":
					dogrid = !dogrid
				case "A":
					shape = "arc"
				case "L":
					shape = "line"
				case "B":
					shape = "bezier"
				case "C":
					shape = "circle"
				case "S":
					shape = "square"
				case "R":
					shape = "rect"
				case "E":
					shape = "ellipse"
				case "D":
					deckspec(prec)
				case key.NameRightArrow:
					switch k.Modifiers {
					case 0:
						bx += stepsize
					case key.ModCtrl:
						ex += stepsize
					}
				case key.NameLeftArrow:
					switch k.Modifiers {
					case 0:
						bx -= stepsize
					case key.ModCtrl:
						ex -= stepsize
					}
				case key.NameUpArrow:
					switch k.Modifiers {
					case 0:
						by += stepsize
					case key.ModCtrl:
						ey += stepsize
					}
				case key.NameDownArrow:
					switch k.Modifiers {
					case 0:
						by -= stepsize
					case key.ModCtrl:
						ey -= stepsize
					}
				case key.NameEscape, "Q":
					os.Exit(0)
				}
			}
		}
		// pointer events
		if p, ok := ev.(pointer.Event); ok {
			switch p.Kind {
			case pointer.Move:
				mouseX, mouseY = pctcoord(p.Position.X, p.Position.Y, width, height)
			case pointer.Press:
				switch p.Buttons {
				case pointer.ButtonPrimary:
					bx, by = pctcoord(p.Position.X, p.Position.Y, width, height)
				case pointer.ButtonSecondary:
					ex, ey = pctcoord(p.Position.X, p.Position.Y, width, height)
				case pointer.ButtonTertiary:
					deckspec(prec)
				}
				pressed = true
			}
		}
	}
}

// shapesketch sketches shapes
// left pointer press defines the begin point, right pointer press defines the end point,
// pointer move defines the current point, arrow keys (plain and shift) adjust begin and end points,
// "G" toggles a grid
// "D" shows the decksh spec
// "L" line
// "B" bezier
// "C" circle
// "R" rectangle
// "E" ellipese
// "S" square
// "Q" or Esc quits
func shapesketch(w *app.Window, cfg config) error {
	// initial values
	bx, by = 25.0, 50.0
	ex, ey = 75.0, 50.0
	cx, cy = 10, 10
	shape = "bezier"
	begincolor, endcolor, shapecolor := cfg.begincolor, cfg.endcolor, cfg.shapecolor

	// app loop
	for {
		ev := w.NextEvent()
		switch e := ev.(type) {
		// return an error on close
		case system.DestroyEvent:
			return e.Err
		// for each frame: register keyboard, pointer press and move events, draw coordinates and
		// specified shapes. Track the pointer position for the current point.
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), system.FrameEvent{})
			key.InputOp{Tag: pressed}.Add(canvas.Context.Ops)
			pointer.InputOp{Tag: pressed, Grab: false, Kinds: pointer.Press | pointer.Move}.Add(canvas.Context.Ops)
			canvas.Background(cfg.bgcolor)
			grid(canvas, 5, cfg.textcolor)
			// draw specified shape
			switch shape {
			case "line":
				textcoord(canvas, bx, by, begincolor, cfg)
				textcoord(canvas, ex, ey, endcolor, cfg)
				canvas.Line(bx, by, ex, ey, cfg.linesize, cfg.shapecolor)
			case "bezier":
				textcoord(canvas, bx, by, begincolor, cfg)
				textcoord(canvas, ex, ey, endcolor, cfg)
				textcoord(canvas, cx, cy, shapecolor, cfg)
				canvas.QuadStrokedCurve(bx, by, cx, cy, ex, ey, cfg.linesize, cfg.shapecolor)
			case "circle":
				textcoord(canvas, bx, by, begincolor, cfg)
				textcoord(canvas, cx, cy, shapecolor, cfg)
				canvas.Circle(bx, by, dist(bx, by, cx, cy), cfg.shapecolor)
			case "square":
				textcoord(canvas, bx, by, begincolor, cfg)
				textcoord(canvas, cx, cy, shapecolor, cfg)
				canvas.Square(bx, by, (cx-bx)*2, cfg.shapecolor)
			case "ellipse":
				textcoord(canvas, bx, by, begincolor, cfg)
				textcoord(canvas, cx, cy, shapecolor, cfg)
				canvas.Ellipse(bx, by, (cx-bx)*2, (cy-by)*2, cfg.shapecolor)
			case "rect":
				textcoord(canvas, bx, by, begincolor, cfg)
				textcoord(canvas, cx, cy, shapecolor, cfg)
				canvas.CenterRect(bx, by, (cx-bx)*2, (cy-by)*2, cfg.shapecolor)
			case "arc":
				textcoord(canvas, bx, by, begincolor, cfg)
				textcoord(canvas, ex, ey, endcolor, cfg)
				textcoord(canvas, cx, cy, shapecolor, cfg)
				r, a2, a1 := arcElements()
				for t := a1; t < a2; t += math.Pi / 256 {
					px, py := canvas.Polar(bx, by, r, float32(t))
					canvas.Line(bx, by, px, py, cfg.linesize, cfg.shapecolor)
				}
			}
			kbpointer(e.Queue, cfg)
			cx, cy = mouseX, mouseY
			e.Frame(canvas.Context.Ops)
		}
	}
}

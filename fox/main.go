// fox  - after Fox I by Anni Albers
package main

import (
	"flag"
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/input"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

const stepsize = 0.5
const minbound = 10
const maxbound = 95
const minstep = 2.0
const maxstep = 20.0
const defaultstep = 5.0
const shadowshift = 0.5
const rangefmt = "%v,%v,%v"

var colorpalette rgbpalette

// config holds configuration parameters
type config struct {
	hue1, hue2                               float64
	beginx, beginy, endx, endy, xstep, ystep float32
	bgcolor, color                           string
}

// random returns a random number between a range
func random(min, max float64) float64 {
	return vmap(rand.Float64(), 0, 1, min, max)
}

// vmap maps one interval to another
func vmap(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

// hsv makes a color given hue, saturation, value
func hsv(hue, sat, value int) color.NRGBA {
	return giocanvas.ColorLookup(fmt.Sprintf("hsv(%d,%d,%d)", hue, sat, value))
}

func randhsv(h1, h2 float64) color.NRGBA {
	return giocanvas.ColorLookup(fmt.Sprintf("hsv(%d,100,100)", int(random(h1, h2))))
}

// parseHues parses a color string: if the string is of the form "h1:h2",
// where h1, and h2 are numbers between 0 and 360, they are a range of hues.
// Otherwise, set to -1 for invalid entries (use named colors instead)
func parseHues(color string) (float64, float64) {
	var h1, h2 float64 = -1.0, -1.0
	hb := strings.Split(color, ":")
	if len(hb) == 2 {
		var err error
		h1, err = strconv.ParseFloat(hb[0], 64)
		if err != nil {
			h1 = -1
		}
		h2, err = strconv.ParseFloat(hb[1], 64)
		if err != nil {
			h2 = -1
		}
	}
	return h1, h2
}

func randpalette() string {
	n := rand.Intn(len(colorpalette))
	i := 0
	for p := range colorpalette {
		if i == n {
			return p
		}
		i++
	}
	return "2-bit-grayscale"
}

var pressed bool
var gbx, gby, gex, gey, gstep, gxstep, gystep float32
var pencolor string

// kbpointer processes the keyboard events and pointer events
func kbpointer(q input.Source, context *op.Ops) {

	for {
		ev, ok := q.Event(
			key.Filter{Optional: key.ModCtrl},
			pointer.Filter{Target: &pressed, Kinds: pointer.Press | pointer.Release})
		if !ok {
			break
		}
		switch e := ev.(type) {

		// keyboard events
		case key.Event:

			switch e.State {
			case key.Press:
				switch e.Name {
				case key.NameRightArrow:
					gbx += stepsize
					gex -= stepsize
				case key.NameLeftArrow:
					gbx -= stepsize
					gex += stepsize
				case key.NameDownArrow:
					gby -= stepsize
					gey += stepsize
				case key.NameUpArrow:
					gby += stepsize
					gey -= stepsize
				case "P":
					pencolor = randpalette()
				case "R":
					gbx, gby = minbound, minbound
					gex, gey = maxbound, maxbound
					gxstep, gystep = minbound, minbound
				case key.NameEscape, "Q":
					os.Exit(0)
				}
			}
		case pointer.Event:
			switch e.Kind {
			case pointer.Press:
				switch e.Buttons {
				case pointer.ButtonPrimary:
					gxstep += stepsize
					gystep += stepsize
				case pointer.ButtonSecondary:
					gxstep -= stepsize
					gystep -= stepsize
				}
			}
		}
	}
	event.Op(context, &pressed)
}

func triangle(canvas *giocanvas.Canvas, x, y, width, height float32, tcolor string, hue1, hue2 float64, direction string) {
	xp := make([]float32, 3)
	yp := make([]float32, 3)
	w2 := width / 2
	h2 := height / 2
	switch direction {
	case "u": // up
		xp[0], xp[1], xp[2] = x, x-w2, x+w2
		yp[0], yp[1], yp[2] = y+h2, y-h2, y-h2
	case "d": // down
		xp[0], xp[1], xp[2] = x, x-w2, x+w2
		yp[0], yp[1], yp[2] = y-h2, y+h2, y+h2
	case "l": // left
		xp[0], xp[1], xp[2] = x-w2, x+w2, x+w2
		yp[0], yp[1], yp[2] = y, y+h2, y-h2
	case "r": // right
		xp[0], xp[1], xp[2] = x+w2, x-w2, x-w2
		yp[0], yp[1], yp[2] = y, y+h2, y-h2
	case "ne": // northeast
		xp[0], xp[1], xp[2] = x-w2, x-w2, x+w2
		yp[0], yp[1], yp[2] = y-h2, y+h2, y+h2
	case "nw": // northwest
		xp[0], xp[1], xp[2] = x-w2, x+w2, x-w2
		yp[0], yp[1], yp[2] = y-h2, y+h2, y+h2
	case "sw": // southwest
		xp[0], xp[1], xp[2] = x+w2, x-w2, x-w2
		yp[0], yp[1], yp[2] = y-h2, y-h2, y+h2
	case "se": // southeast
		xp[0], xp[1], xp[2] = x-w2, x+w2, x+w2
		yp[0], yp[1], yp[2] = y-h2, y-h2, y+h2
	}

	var fillcolor color.NRGBA
	fillcolor = giocanvas.ColorLookup(tcolor) // default to named color
	if hue1 > -1 && hue2 > -1 {               // use hue
		fillcolor = randhsv(hue1, hue2)
	}
	if c, ok := colorpalette[tcolor]; ok { // use a palette
		fillcolor = c[rand.Intn(len(c))]
	}
	canvas.Polygon(xp, yp, fillcolor)
}

// fox makes...
func fox(w *app.Window, width, height float32, cfg config) error {
	bg := giocanvas.ColorLookup(cfg.bgcolor)
	pencolor = cfg.color
	gxstep = cfg.xstep
	gystep = cfg.ystep
	gbx, gex, gby, gey = cfg.beginx, cfg.endx, cfg.beginy, cfg.endy
	var directions = []string{"u", "d", "l", "r", "nw", "ne", "sw", "se"}

	for {
		switch e := w.NextEvent().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:

			if gxstep < minstep || gxstep > maxstep {
				gxstep = defaultstep
			}
			if gystep < minstep || gystep > maxstep {
				gystep = defaultstep
			}
			if gbx < minbound || gbx > maxbound {
				gbx = minbound
			}
			if gex < minbound || gex > maxbound {
				gex = maxbound
			}
			if gby < minbound || gby > maxbound {
				gby = minbound
			}
			if gey < minbound || gey > maxbound {
				gey = maxbound
			}
			if gbx > gex {
				gbx = minbound
				gex = maxbound
			}
			if gby > gey {
				gby = minbound
				gey = maxbound
			}
			//fmt.Fprintf(os.Stderr, "x=(%v,%v,%v) y=(%v,%v,%v)\n", gbx, gex, gxstep, gby, gey, gystep)

			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})
			canvas.Background(bg)
			for y := gby; y < gey; y += gystep {
				for x := gbx; x < gex; x += gxstep {
					w := float32(random(minstep, float64(gxstep)))
					h := float32(random(minstep, float64(gystep)))
					triangle(canvas, x, y, w, h, pencolor, cfg.hue1, cfg.hue2, directions[rand.Intn(len(directions))])
					triangle(canvas, x+shadowshift, y-shadowshift, w, h, pencolor, cfg.hue1, cfg.hue2, directions[rand.Intn(4)])

				}
			}
			kbpointer(e.Source, canvas.Context.Ops)
			e.Frame(canvas.Context.Ops)
		}
	}
}

func usage() {
	defrange := fmt.Sprintf(rangefmt, minbound, maxbound, defaultstep)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Option      Default     Description\n")
	fmt.Fprintf(os.Stderr, ".............................................................\n")
	fmt.Fprintf(os.Stderr, "-help       false       show usage\n")
	fmt.Fprintf(os.Stderr, "-width      1000        canvas width\n")
	fmt.Fprintf(os.Stderr, "-height     1000        canvas height\n")
	fmt.Fprintf(os.Stderr, "-w          "+defrange+"     percent begin,end,step for the width\n")
	fmt.Fprintf(os.Stderr, "-h          "+defrange+"     percent begin,end,step for the height\n")
	fmt.Fprintf(os.Stderr, "-p          \"\"          palette file\n")
	fmt.Fprintf(os.Stderr, "-bgcolor    white       background color\n")
	fmt.Fprintf(os.Stderr, "-color      gray        color name, h1:h2, or palette:\n\n")
	for p, k := range colorpalette {
		fmt.Fprintf(os.Stderr, "%-20s\t", p)
		end := len(k) - 1
		for i := 0; i < end; i++ {
			fmt.Fprintf(os.Stderr, "#%02x%02x%02x ", k[i].R, k[i].G, k[i].B)
		}
		fmt.Fprintf(os.Stderr, "#%02x%02x%02x\n", k[end].R, k[end].G, k[end].B)
	}
	os.Exit(1)
}

func parserange(s string) (float32, float32, float32) {
	v := strings.Split(s, ",")
	if len(v) == 3 {
		min, err := strconv.ParseFloat(v[0], 32)
		if err != nil {
			min = 5
		}
		max, err := strconv.ParseFloat(v[1], 32)
		if err != nil {
			max = 100
		}
		step, err := strconv.ParseFloat(v[2], 32)
		if err != nil {
			step = 5
		}
		return float32(min), float32(max), float32(step)
	}
	return minbound, maxbound, defaultstep
}

// convert the built-in palette from string based to rgb
func convertpalette() {
	colorpalette = make(rgbpalette)
	for name, value := range palette {
		colors := make([]color.NRGBA, len(value))
		for i, c := range value {
			x, _ := strconv.ParseUint(c[1:], 16, 32)
			r, g, b := rgb(uint32(x))
			colors[i] = color.NRGBA{R: r, G: g, B: b, A: 0xff}
		}
		colorpalette[name] = colors
	}
}

// load a palette from a file
func loadpalette(pfile string) {
	var err error
	colorpalette, err = LoadRGBPalette(pfile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func main() {
	var cw, ch int
	var cfg config
	var showhelp bool
	var xconfig, yconfig, pfile string
	defrange := fmt.Sprintf(rangefmt, minbound, maxbound, defaultstep)
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.BoolVar(&showhelp, "help", false, "show usage")
	flag.StringVar(&xconfig, "w", defrange, "horizontal config (min,max,step)")
	flag.StringVar(&yconfig, "h", defrange, "vertical config (min,max,step)")
	flag.StringVar(&pfile, "p", "", "palette file")
	flag.StringVar(&cfg.bgcolor, "bgcolor", "white", "background color")
	flag.StringVar(&cfg.color, "color", "gray", "pen color; named color, palette, or h1:h2 for a random hue range hsv(h1:h2, 100, 100)")
	flag.Parse()

	if len(pfile) > 0 {
		loadpalette(pfile)
	} else {
		convertpalette()
	}
	if showhelp {
		usage()
	}
	cfg.beginx, cfg.endx, cfg.xstep = parserange(xconfig)
	cfg.beginy, cfg.endy, cfg.ystep = parserange(yconfig)
	cfg.hue1, cfg.hue2 = parseHues(cfg.color)
	width, height := float32(cw), float32(ch)

	go func() {
		w := app.NewWindow(app.Title("fox"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := fox(w, width, height, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create the window: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

// fox  - after Fox I bu Anni Albers
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
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
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

var palette = map[string][]color.NRGBA{
	"nostalgia":              {{R: 0xd0, G: 0xd0, B: 0x58, A: 0xff}, {R: 0xa0, G: 0xa8, B: 0x40, A: 0xff}, {R: 0x70, G: 0x80, B: 0x28, A: 0xff}, {R: 0x40, G: 0x50, B: 0x10, A: 0xff}},
	"spacehaze":              {{R: 0xf8, G: 0xe3, B: 0xc4, A: 0xff}, {R: 0xcc, G: 0x34, B: 0x95, A: 0xff}, {R: 0x6b, G: 0x1f, B: 0xb1, A: 0xff}, {R: 0x0b, G: 0x06, B: 0x30, A: 0xff}},
	"dark-mode":              {{R: 0x21, G: 0x21, B: 0x21, A: 0xff}, {R: 0x45, G: 0x45, B: 0x45, A: 0xff}, {R: 0x78, G: 0x78, B: 0x78, A: 0xff}, {R: 0xa8, G: 0xa5, B: 0xa5, A: 0xff}},
	"autumn-decay":           {{R: 0x31, G: 0x36, B: 0x38, A: 0xff}, {R: 0x57, G: 0x47, B: 0x29, A: 0xff}, {R: 0x97, G: 0x53, B: 0x30, A: 0xff}, {R: 0xc5, G: 0x79, B: 0x38, A: 0xff}, {R: 0xff, G: 0xad, B: 0x3b, A: 0xff}, {R: 0xff, G: 0xe5, B: 0x96, A: 0xff}},
	"kirokaze-gameboy":       {{R: 0x33, G: 0x2c, B: 0x50, A: 0xff}, {R: 0x46, G: 0x87, B: 0x8f, A: 0xff}, {R: 0x94, G: 0xe3, B: 0x44, A: 0xff}, {R: 0xe2, G: 0xf3, B: 0xe4, A: 0xff}},
	"moonlight-gb":           {{R: 0x0f, G: 0x05, B: 0x2d, A: 0xff}, {R: 0x20, G: 0x36, B: 0x71, A: 0xff}, {R: 0x36, G: 0x86, B: 0x8f, A: 0xff}, {R: 0x5f, G: 0xc7, B: 0x5d, A: 0xff}},
	"mist-gb":                {{R: 0x2d, G: 0x1b, B: 0x00, A: 0xff}, {R: 0x1e, G: 0x60, B: 0x6e, A: 0xff}, {R: 0x5a, G: 0xb9, B: 0xa8, A: 0xff}, {R: 0xc4, G: 0xf0, B: 0xc2, A: 0xff}},
	"arq4":                   {{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, {R: 0x67, G: 0x72, B: 0xa9, A: 0xff}, {R: 0x3a, G: 0x32, B: 0x77, A: 0xff}, {R: 0x00, G: 0x00, B: 0x00, A: 0xff}},
	"pen-n-paper":            {{R: 0xe4, G: 0xdb, B: 0xba, A: 0xff}, {R: 0xa4, G: 0x92, B: 0x9a, A: 0xff}, {R: 0x4f, G: 0x3a, B: 0x54, A: 0xff}, {R: 0x26, G: 0x0d, B: 0x1c, A: 0xff}},
	"hollow":                 {{R: 0x0f, G: 0x0f, B: 0x1b, A: 0xff}, {R: 0x56, G: 0x5a, B: 0x75, A: 0xff}, {R: 0xc6, G: 0xb7, B: 0xbe, A: 0xff}, {R: 0xfa, G: 0xfb, B: 0xf6, A: 0xff}},
	"pokemon-sgb":            {{R: 0x18, G: 0x10, B: 0x10, A: 0xff}, {R: 0x84, G: 0x73, B: 0x9c, A: 0xff}, {R: 0xf7, G: 0xb5, B: 0x8c, A: 0xff}, {R: 0xff, G: 0xef, B: 0xff, A: 0xff}},
	"kankei4":                {{R: 0xff, G: 0xff, B: 0xff, A: 0xff}, {R: 0xf4, G: 0x2e, B: 0x1f, A: 0xff}, {R: 0x2f, G: 0x25, B: 0x6b, A: 0xff}, {R: 0x06, G: 0x06, B: 0x08, A: 0xff}},
	"polished-gold":          {{R: 0x00, G: 0x00, B: 0x00, A: 0xff}, {R: 0x36, G: 0x1c, B: 0x1b, A: 0xff}, {R: 0x75, G: 0x42, B: 0x32, A: 0xff}, {R: 0xcd, G: 0x89, B: 0x4a, A: 0xff}, {R: 0xe6, G: 0xb9, B: 0x83, A: 0xff}, {R: 0xff, G: 0xf8, B: 0xbc, A: 0xff}, {R: 0xff, G: 0xff, B: 0xff, A: 0xff}, {R: 0x2d, G: 0x24, B: 0x33, A: 0xff}, {R: 0x4f, G: 0x42, B: 0x54, A: 0xff}, {R: 0xb0, G: 0x92, B: 0xa7, A: 0xff}},
	"2-bit-grayscale":        {{R: 0x00, G: 0x00, B: 0x00, A: 0xff}, {R: 0x67, G: 0x67, B: 0x67, A: 0xff}, {R: 0xb6, G: 0xb6, B: 0xb6, A: 0xff}, {R: 0xff, G: 0xff, B: 0xff, A: 0xff}},
	"nintendo-super-gameboy": {{R: 0x33, G: 0x1e, B: 0x50, A: 0xff}, {R: 0xa6, G: 0x37, B: 0x25, A: 0xff}, {R: 0xd6, G: 0x8e, B: 0x49, A: 0xff}, {R: 0xf7, G: 0xe7, B: 0xc6, A: 0xff}},
	"rustic-gb":              {{R: 0x2c, G: 0x21, B: 0x37, A: 0xff}, {R: 0x76, G: 0x44, B: 0x62, A: 0xff}, {R: 0xed, G: 0xb4, B: 0xa1, A: 0xff}, {R: 0xa9, G: 0x68, B: 0x68, A: 0xff}},
	"links-awakening-sgb":    {{R: 0x5a, G: 0x39, B: 0x21, A: 0xff}, {R: 0x6b, G: 0x8c, B: 0x42, A: 0xff}, {R: 0x7b, G: 0xc6, B: 0x7b, A: 0xff}, {R: 0xff, G: 0xff, B: 0xb5, A: 0xff}},
	"blk-aqu4":               {{R: 0x00, G: 0x2b, B: 0x59, A: 0xff}, {R: 0x00, G: 0x5f, B: 0x8c, A: 0xff}, {R: 0x00, G: 0xb9, B: 0xbe, A: 0xff}, {R: 0x9f, G: 0xf4, B: 0xe5, A: 0xff}},
	"ice-cream-gb":           {{R: 0x7c, G: 0x3f, B: 0x58, A: 0xff}, {R: 0xeb, G: 0x6b, B: 0x6f, A: 0xff}, {R: 0xf9, G: 0xa8, B: 0x75, A: 0xff}, {R: 0xff, G: 0xf6, B: 0xd3, A: 0xff}},
	"ayy4":                   {{R: 0x00, G: 0x30, B: 0x3b, A: 0xff}, {R: 0xff, G: 0x77, B: 0x77, A: 0xff}, {R: 0xff, G: 0xce, B: 0x96, A: 0xff}, {R: 0xf1, G: 0xf2, B: 0xda, A: 0xff}},
	"nintendo-gameboy-bgb":   {{R: 0x08, G: 0x18, B: 0x20, A: 0xff}, {R: 0x34, G: 0x68, B: 0x56, A: 0xff}, {R: 0x88, G: 0xc0, B: 0x70, A: 0xff}, {R: 0xe0, G: 0xf8, B: 0xd0, A: 0xff}},
	"blu-scribbles":          {{R: 0x05, G: 0x18, B: 0x33, A: 0xff}, {R: 0x0a, G: 0x4f, B: 0x66, A: 0xff}, {R: 0x0f, G: 0x99, B: 0x8e, A: 0xff}, {R: 0x12, G: 0xcc, B: 0x7f, A: 0xff}},
	"ajstarks":               {{R: 0xaa, G: 0x00, B: 0x00, A: 0xff}, {R: 0xaa, G: 0xaa, B: 0xaa, A: 0xff}, {R: 0x00, G: 0x00, B: 0x00, A: 0xff}, {R: 0xff, G: 0xff, B: 0xff, A: 0xff}},
	"funk-it-up":             {{R: 0xe4, G: 0xff, B: 0xff, A: 0xff}, {R: 0xe6, G: 0x34, B: 0x10, A: 0xff}, {R: 0xa2, G: 0x37, B: 0x37, A: 0xff}, {R: 0xff, G: 0xec, B: 0x40, A: 0xff}, {R: 0x81, G: 0x91, B: 0x3b, A: 0xff}, {R: 0x26, G: 0xf6, B: 0x75, A: 0xff}, {R: 0x4c, G: 0x71, B: 0x4e, A: 0xff}, {R: 0x40, G: 0xeb, B: 0xda, A: 0xff}, {R: 0x39, G: 0x4e, B: 0x4e, A: 0xff}, {R: 0x0a, G: 0x0a, B: 0x0a, A: 0xff}},
	"2-bit-demichrome":       {{R: 0x21, G: 0x1e, B: 0x20, A: 0xff}, {R: 0x55, G: 0x55, B: 0x68, A: 0xff}, {R: 0xa0, G: 0xa0, B: 0x8b, A: 0xff}, {R: 0xe9, G: 0xef, B: 0xec, A: 0xff}},
	"red-brick":              {{R: 0xef, G: 0xf9, B: 0xd6, A: 0xff}, {R: 0xba, G: 0x50, B: 0x44, A: 0xff}, {R: 0x7a, G: 0x1c, B: 0x4b, A: 0xff}, {R: 0x1b, G: 0x03, B: 0x26, A: 0xff}},
}

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
	n := rand.Intn(len(palette))
	i := 0
	for p := range palette {
		if i == n {
			println(p)
			return p
		}
		i++
	}
	return "ajstarks"
}

var pressed bool
var gbx, gby, gex, gey, gstep, gxstep, gystep float32
var pencolor string

// kbpointer processes the keyboard events and pointer events
func kbpointer(q event.Queue) {

	for _, ev := range q.Events(pressed) {
		// keyboard events
		if k, ok := ev.(key.Event); ok {
			switch k.State {
			case key.Press:
				switch k.Name {
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
		}
		// pointer events
		if p, ok := ev.(pointer.Event); ok {
			switch p.Kind {
			case pointer.Press:
				switch p.Buttons {
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
	}

	var fillcolor color.NRGBA
	fillcolor = giocanvas.ColorLookup(tcolor) // default to named color
	if hue1 > -1 && hue2 > -1 {               // use hue
		fillcolor = randhsv(hue1, hue2)
	}
	if c, ok := palette[tcolor]; ok { // use a palette
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
	var directions = []string{"u", "d", "l", "r"}

	for {
		e := w.NextEvent()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), system.FrameEvent{})
			key.InputOp{Tag: pressed}.Add(canvas.Context.Ops)
			pointer.InputOp{Tag: pressed, Kinds: pointer.Press | pointer.Release}.Add(canvas.Context.Ops)
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
			canvas.Background(bg)
			for y := gby; y < gey; y += gystep {
				for x := gbx; x < gex; x += gxstep {
					w := float32(random(minstep, float64(gxstep)))
					h := float32(random(minstep, float64(gystep)))
					triangle(canvas, x, y, w, h, pencolor, cfg.hue1, cfg.hue2, directions[rand.Intn(4)])
					triangle(canvas, x+shadowshift, y-shadowshift, w, h, pencolor, cfg.hue1, cfg.hue2, directions[rand.Intn(4)])

				}
			}
			kbpointer(e.Queue)
			e.Frame(canvas.Context.Ops)
		}
	}
}

func usage() {
	defrange := fmt.Sprintf(rangefmt, minbound, maxbound, defaultstep)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Option      Default     Description\n")
	fmt.Fprintf(os.Stderr, ".....................................................\n")
	fmt.Fprintf(os.Stderr, "-help       false       show usage\n")
	fmt.Fprintf(os.Stderr, "-width      1000        canvas width\n")
	fmt.Fprintf(os.Stderr, "-height     1000        canvas height\n")
	fmt.Fprintf(os.Stderr, "-w          "+defrange+"     percent begin,end,step for the width\n")
	fmt.Fprintf(os.Stderr, "-h          "+defrange+"     percent begin,end,step for the height\n")
	fmt.Fprintf(os.Stderr, "-bgcolor    white       background color\n")
	fmt.Fprintf(os.Stderr, "-color      gray        color name, h1:h2, or palette:\n\n")
	for p, k := range palette {
		fmt.Fprintf(os.Stderr, "%-20s\t", p)
		end := len(k) - 1
		for i := 0; i < end; i++ {
			fmt.Fprintf(os.Stderr, "#%02x%02x%02x, ", k[i].R, k[i].G, k[i].B)
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

func main() {
	var cw, ch int
	var cfg config
	var showhelp bool
	var xconfig, yconfig string
	defrange := fmt.Sprintf(rangefmt, minbound, maxbound, defaultstep)
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.BoolVar(&showhelp, "help", false, "show usage")
	flag.StringVar(&xconfig, "w", defrange, "horizontal config (min,max,step)")
	flag.StringVar(&yconfig, "h", defrange, "vertical config (min,max,step)")
	flag.StringVar(&cfg.bgcolor, "bgcolor", "white", "background color")
	flag.StringVar(&cfg.color, "color", "gray", "pen color; named color, palette, or h1:h2 for a random hue range hsv(h1:h2, 100, 100)")
	flag.Parse()
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

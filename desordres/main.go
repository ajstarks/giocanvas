// desordres - after Des Ordres by Vera Molnar
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
	tiles, maxlw, h1, h2 float64
	bgcolor, color       string
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

// csquare makes squares, with possibly random colors, centered at (x,y)
func csquare(canvas *giocanvas.Canvas, x, y, size, maxlw, h1, h2 float64, linecolor string) {
	var color color.NRGBA

	if c, ok := palette[linecolor]; ok { // use a palette
		color = c[rand.Intn(len(c)-1)]
	} else if h1 > -1 && h2 > -1 { // hue range set
		color = hsv(int(random(h1, h2)), 100, 100)
	} else {
		color = giocanvas.ColorLookup(linecolor)
	}
	// define the corners
	hs := size / 2
	tlx, tly := float32(x-hs), float32(y+hs)
	trx, try := float32(x+hs), float32(y+hs)
	blx, bly := float32(x-hs), float32(y-hs)
	brx, bry := float32(x+hs), float32(y-hs)
	// make the boundaries
	lw := float32(random(0.1, maxlw))
	ll := float32(size)
	canvas.HLine(tlx, tly, ll, lw, color)
	canvas.HLine(blx, bly, ll, lw, color)
	canvas.VLine(blx, bly, ll, lw, color)
	canvas.VLine(brx, bry, ll, lw, color)
	// make the corners
	canvas.Square(tlx, tly, lw, color)
	canvas.Square(blx, bly, lw, color)
	canvas.Square(brx, bry, lw, color)
	canvas.Square(trx, try, lw, color)
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

// tiles makes a tile of concentric squares
func tiles(canvas *giocanvas.Canvas, x, y, minsize, maxsize, maxlw, h1, h2 float64, color string) {
	step := random(1, 5)
	for v := minsize; v < maxsize; v += step {
		csquare(canvas, x, y, v, maxlw, h1, h2, color)
	}
}

func randpalette() string {
	n := rand.Intn(len(palette) - 1)
	i := 0
	for p := range palette {
		if i == n {
			return p
		}
		i++
	}
	return "ajstarks"
}

var pressed bool
var tilesize float64
var pencolor string

const stepsize = 1.0
const mintile = 1.0
const maxtile = 20.0

// kbpointer processes the keyboard events and pointer events
func kbpointer(q input.Source, context *op.Ops) {

	for {
		e, ok := q.Event(
			key.Filter{Optional: key.ModCtrl},
			pointer.Filter{Target: &pressed, Kinds: pointer.Press},
		)
		if !ok {
			break
		}
		switch e := e.(type) {
		case key.Event: // keyboard events

			switch e.State {
			case key.Press:
				switch e.Name {
				case key.NameHome:
					tilesize = mintile
				case key.NameEnd:
					tilesize = maxtile
				case key.NameLeftArrow, key.NameDownArrow, "-":
					tilesize -= stepsize
				case key.NameRightArrow, key.NameUpArrow, "+":
					tilesize += stepsize
				case "P":
					pencolor = randpalette()
				case key.NameEscape, "Q":
					os.Exit(0)
				}
			}

		case pointer.Event: // pointer events

			switch e.Kind {
			case pointer.Press:
				switch e.Buttons {
				case pointer.ButtonPrimary:
					tilesize += stepsize
				case pointer.ButtonSecondary:
					tilesize -= stepsize
				case pointer.ButtonTertiary:
					tilesize = 10
				}
			}
		}
	}
	event.Op(context, &pressed)
}

// desordres makes tiles of random conentric squares
func desordres(w *app.Window, width, height float32, cfg config) error {
	bg := giocanvas.ColorLookup(cfg.bgcolor)
	maxlw := cfg.maxlw
	h1, h2 := parseHues(cfg.color) // set hue range, or named color
	pencolor = cfg.color
	tilesize = cfg.tiles
	var top, left, size float64
	for {
		e := w.NextEvent()
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), app.FrameEvent{})

			canvas.Background(bg)
			if tilesize < mintile {
				tilesize = mintile
			}
			if tilesize > maxtile {
				tilesize = maxtile
			}
			size = 100 / tilesize  // size of each tile
			top = 100 - (size / 2) // top of the beginning row
			left = 100 - top       // left of the beginning row
			for y := top; y > 0; y -= size {
				for x := left; x < 100; x += size {
					tiles(canvas, x, y, 2, size, maxlw, h1, h2, pencolor)
				}
			}
			kbpointer(e.Source, canvas.Context.Ops)
			e.Frame(canvas.Context.Ops)
		}
	}
}

func usage() {
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Option      Default    Description\n")
	fmt.Fprintf(os.Stderr, ".....................................................\n")
	fmt.Fprintf(os.Stderr, "-help       false      show usage\n")
	fmt.Fprintf(os.Stderr, "-width      1000       canvas width\n")
	fmt.Fprintf(os.Stderr, "-height     1000       canvas height\n")
	fmt.Fprintf(os.Stderr, "-tiles      10         number of tiles/row\n")
	fmt.Fprintf(os.Stderr, "-maxlw      1          maximim line thickness\n")
	fmt.Fprintf(os.Stderr, "-bgcolor    white      background color\n")
	fmt.Fprintf(os.Stderr, "-color      gray       color name, h1:h2, or palette:\n")
	for p, k := range palette {
		fmt.Fprintf(os.Stderr, "                       %-30s\t%v\n", p, k)
	}
	os.Exit(1)

}

func main() {
	var cw, ch int
	var cfg config
	var showhelp bool

	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.BoolVar(&showhelp, "help", false, "show usage")
	flag.Float64Var(&cfg.tiles, "tiles", 10, "tiles/row")
	flag.Float64Var(&cfg.maxlw, "maxlw", 1, "maximum line thickness")
	flag.StringVar(&cfg.bgcolor, "bgcolor", "white", "background color")
	flag.StringVar(&cfg.color, "color", "gray", "pen color; named color, or h1:h2 for a random hue range hsv(h1:h2, 100, 100)")
	flag.Parse()

	width, height := float32(cw), float32(ch)
	if width != height {
		fmt.Fprintln(os.Stderr, "width and height must be the same")
		os.Exit(1)
	}

	if showhelp {
		usage()
	}

	go func() {
		w := app.NewWindow(app.Title("desordres"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := desordres(w, width, height, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create the window: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

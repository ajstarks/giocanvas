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

var colorpalette rgbpalette

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

	if c, ok := colorpalette[linecolor]; ok { // use a palette
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
	n := rand.Intn(len(colorpalette) - 1)
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
func desordres(w *app.Window, cfg config) error {
	bg := giocanvas.ColorLookup(cfg.bgcolor)
	maxlw := cfg.maxlw
	h1, h2 := parseHues(cfg.color) // set hue range, or named color
	pencolor = cfg.color
	tilesize = cfg.tiles
	var top, left, size float64
	for {
		e := w.Event()
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
	fmt.Fprintf(os.Stderr, "Option      Default     Description\n")
	fmt.Fprintf(os.Stderr, ".....................................................\n")
	fmt.Fprintf(os.Stderr, "-help       false       show usage\n")
	fmt.Fprintf(os.Stderr, "-width      1000        canvas width\n")
	fmt.Fprintf(os.Stderr, "-height     1000        canvas height\n")
	fmt.Fprintf(os.Stderr, "-tiles      10          number of tiles/row\n")
	fmt.Fprintf(os.Stderr, "-maxlw      1           maximim line thickness\n")
	fmt.Fprintf(os.Stderr, "-bgcolor    white       background color\n")
	fmt.Fprintf(os.Stderr, "-p          \"\"          palette file\n")
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
	var pfile string

	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.BoolVar(&showhelp, "help", false, "show usage")
	flag.StringVar(&pfile, "p", "", "palette file")
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
	if len(pfile) > 0 {
		loadpalette(pfile)
	} else {
		convertpalette()
	}
	if showhelp {
		usage()
	}

	go func() {
		w := &app.Window{}
		w.Option(app.Title("desordres"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := desordres(w, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create the window: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

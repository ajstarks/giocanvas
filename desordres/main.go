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
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

var palette = map[string][]string{
	"kirokaze-gameboy":       {"#332c50", "#46878f", "#94e344", "#e2f3e4"},
	"ice-cream-gb":           {"#7c3f58", "#eb6b6f", "#f9a875", "#fff6d3"},
	"2-bit-demichrome":       {"#211e20", "#555568", "#a0a08b", "#e9efec"},
	"mist-gb":                {"#2d1b00", "#1e606e", "#5ab9a8", "#c4f0c2"},
	"rustic-gb":              {"#2c2137", "#764462", "#edb4a1", "#a96868"},
	"2-bit-grayscale":        {"#000000", "#676767", "#b6b6b6", "#ffffff"},
	"hollow":                 {"#0f0f1b", "#565a75", "#c6b7be", "#fafbf6"},
	"ayy4":                   {"#00303b", "#ff7777", "#ffce96", "#f1f2da"},
	"nintendo-gameboy-bgb":   {"#081820", "#346856", "#88c070", "#e0f8d0"},
	"red-brick":              {"#eff9d6", "#ba5044", "#7a1c4b", "#1b0326"},
	"nostalgia":              {"#d0d058", "#a0a840", "#708028", "#405010"},
	"spacehaze":              {"#f8e3c4", "#cc3495", "#6b1fb1", "#0b0630"},
	"moonlight-gb":           {"#0f052d", "#203671", "#36868f", "#5fc75d"},
	"links-awakening-sgb":    {"#5a3921", "#6b8c42", "#7bc67b", "#ffffb5"},
	"arq4":                   {"#ffffff", "#6772a9", "#3a3277", "#000000"},
	"blk-aqu4":               {"#002b59", "#005f8c", "#00b9be", "#9ff4e5"},
	"pokemon-sgb":            {"#181010", "#84739c", "#f7b58c", "#ffefff"},
	"nintendo-super-gameboy": {"#331e50", "#a63725", "#d68e49", "#f7e7c6"},
	"blu-scribbles":          {"#051833", "#0a4f66", "#0f998e", "#12cc7f"},
	"kankei4":                {"#ffffff", "#f42e1f", "#2f256b", "#060608"},
	"dark-mode":              {"#212121", "#454545", "#787878", "#a8a5a5"},
	"ajstarks":               {"#aa0000", "#aaaaaa", "#000000", "#ffffff"},
	"pen-n-paper":            {"#e4dbba", "#a4929a", "#4f3a54", "#260d1c"},
	"autumn-decay":           {"#313638", "#574729", "#975330", "#c57938", "#ffad3b", "#ffe596"},
	"polished-gold":          {"#000000", "#361c1b", "#754232", "#cd894a", "#e6b983", "#fff8bc", "#ffffff", "#2d2433", "#4f4254", "#b092a7"},
	"funk-it-up":             {"#e4ffff", "#e63410", "#a23737", "#ffec40", "#81913b", "#26f675", "#4c714e", "#40ebda", "#394e4e", "#0a0a0a"},
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

	color = giocanvas.ColorLookup(linecolor)

	if c, ok := palette[linecolor]; ok { // use a palette
		color = giocanvas.ColorLookup(c[rand.Intn(len(c)-1)])
	}
	if h1 > -1 && h2 > -1 { // hue range set
		color = hsv(int(random(h1, h2)), 100, 100)
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
func kbpointer(q event.Queue) {

	for _, ev := range q.Events(pressed) {
		// keyboard events
		if k, ok := ev.(key.Event); ok {
			switch k.State {
			case key.Press:
				switch k.Name {
				case key.NameHome:
					tilesize = mintile
				case key.NameEnd:
					tilesize = maxtile
				case key.NameLeftArrow, key.NameDownArrow, "-":
				case key.NameRightArrow, key.NameUpArrow, "+":
					tilesize += stepsize
				case "P":
					pencolor = randpalette()
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
					tilesize += stepsize
				case pointer.ButtonSecondary:
					tilesize -= stepsize
				case pointer.ButtonTertiary:
					tilesize = 10
				}
			}
		}
	}
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
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), system.FrameEvent{})
			key.InputOp{Tag: pressed}.Add(canvas.Context.Ops)
			pointer.InputOp{Tag: pressed, Grab: false, Kinds: pointer.Press}.Add(canvas.Context.Ops)

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
			kbpointer(e.Queue)
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
	flag.StringVar(&cfg.color, "color", "gray", "pen color; named color, palette, or h1:h2 for a random hue range hsv(h1:h2, 100, 100)")
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

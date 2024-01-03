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
	"gioui.org/io/system"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

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
	lw := float32(random(0.1, maxlw))
	ll := float32(size)
	var color color.NRGBA
	if h1 > -1 && h2 > -1 { // hue range set
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

var pressed bool
var tilesize float64

const stepsize = 1.0

// kbpointer processes the keyboard events and pointer events in percent coordinates
func kbpointer(q event.Queue, cfg config) {

	for _, ev := range q.Events(pressed) {
		// keyboard events
		if k, ok := ev.(key.Event); ok {
			switch k.State {
			case key.Press:
				switch k.Name {
				case key.NameLeftArrow:
					switch k.Modifiers {
					case 0:
						tilesize -= stepsize
					case key.ModCtrl:
						tilesize -= stepsize
					}
				case key.NameRightArrow:
					switch k.Modifiers {
					case 0:
						tilesize += stepsize
					case key.ModCtrl:
						tilesize += stepsize
					}
				case key.NameUpArrow:
					switch k.Modifiers {
					case 0:
						tilesize += stepsize
					case key.ModCtrl:
						tilesize += stepsize
					}
				case key.NameDownArrow:
					switch k.Modifiers {
					case 0:
						tilesize -= stepsize
					case key.ModCtrl:
						tilesize -= stepsize
					}
				case key.NameEscape, "Q":
					os.Exit(0)
				}
			}
		}
		pressed = true
	}
}

// desordres makes tiles of random conentric squares
func desordres(w *app.Window, width, height float32, cfg config) error {
	bg := giocanvas.ColorLookup(cfg.bgcolor)
	maxlw := cfg.maxlw
	h1, h2 := parseHues(cfg.color) // set hue range, or named color
	color := cfg.color
	tilesize = cfg.tiles
	var top, left, size float64
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return e.Err
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), system.FrameEvent{})
			key.InputOp{Tag: pressed}.Add(canvas.Context.Ops)
			canvas.Background(bg)
			if tilesize < 1 {
				tilesize = 20
			}
			if tilesize > 20 {
				tilesize = 1
			}
			size = 100 / tilesize  // size of each tile
			top = 100 - (size / 2) // top of the beginning row
			left = 100 - top       // left of the beginning row
			for y := top; y > 0; y -= size {
				for x := left; x < 100; x += size {
					tiles(canvas, x, y, 2, size, maxlw, h1, h2, color)
				}
			}
			kbpointer(e.Queue, cfg)
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {
	var cw, ch int
	var cfg config
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.Float64Var(&cfg.tiles, "tiles", 10, "tiles/row")
	flag.Float64Var(&cfg.maxlw, "maxlw", 1, "maximum line thickness")
	flag.StringVar(&cfg.bgcolor, "bgcolor", "white", "background color")
	flag.StringVar(&cfg.color, "color", "gray", "pen color")
	flag.Parse()

	width, height := float32(cw), float32(ch)
	if width != height {
		fmt.Fprintln(os.Stderr, "width and height must be the same")
		os.Exit(1)
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

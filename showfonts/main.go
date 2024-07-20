// showfonts: show fonts on a gio canvas
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/unit"
	"github.com/ajstarks/giocanvas"
)

// config options
type config struct {
	text, bgcolor, txcolor string
	width, height          int
	ts                     float64
}

// basename returns the basename of a path, removing extension
func basename(s, ext string) string {
	s = filepath.Base(s)
	i := strings.Index(s, ext)
	if i > 0 {
		return s[0:i]
	}
	return s
}

// loadfonts reads ttf files and returns a Gio font collection
func loadfonts(fonts []string) ([]font.FontFace, error) {
	collection := []font.FontFace{}
	fc := font.FontFace{}
	for _, v := range fonts {
		fontdata, err := os.ReadFile(v)
		if err != nil {
			return collection, fmt.Errorf("%s: %v\n", v, err)
		}
		face, err := opentype.Parse(fontdata)
		if err != nil {
			return collection, fmt.Errorf("%s: %v\n", v, err)
		}
		fc.Font.Typeface = font.Typeface(basename(v, ".ttf"))
		fc.Face = face
		collection = append(collection, fc)
	}
	return collection, nil
}

// showfonts displays a list of fonts in the named font, with optional message
func showfonts(title string, fontnames []string, cfg config) {
	ts := float32(cfg.ts)
	cw := float32(cfg.width)
	ch := float32(cfg.height)
	message := cfg.text
	// compute text sizing
	const (
		top    = 90.0
		bottom = 100 - top
		left   = 5.0
		right  = 100 - left
		mid    = 50.0
	)
	var y, yskip float32

	nf := float32(len(fontnames) - 1)
	if ts == 0 && nf > 0 {
		ts = (top - bottom) / nf
		ts /= 1.4
	}
	if ts > 10 || nf < 2 {
		ts = 5
	}
	yskip = ts * 1.4

	// load fonts
	fc, err := loadfonts(fontnames)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	// set backfound and text colors
	bg := giocanvas.ColorLookup(cfg.bgcolor)
	fg := giocanvas.ColorLookup(cfg.txcolor)

	// run the app
	w := &app.Window{}
	w.Option(app.Title(title), app.Size(unit.Dp(cw), unit.Dp(ch)))
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			canvas := giocanvas.NewCanvasFonts(float32(e.Size.X), float32(e.Size.Y), fc, app.FrameEvent{})
			canvas.Background(bg)
			y = top
			// show fonts in a vertical list
			for _, s := range fontnames {
				name := basename(s, ".ttf")
				canvas.Theme.Face = font.Typeface(name) // set the font that was preloaded by name
				if len(message) > 0 {
					canvas.Text(left, y, ts, message, fg)
					canvas.TextEnd(right, y, ts*.4, name, fg)
				} else {
					canvas.TextMid(mid, y, ts, name, fg)
				}
				y -= yskip
			}
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {
	var options config
	flag.StringVar(&options.text, "text", "", "text to show")
	flag.StringVar(&options.bgcolor, "bgcolor", "white", "background color")
	flag.StringVar(&options.txcolor, "txcolor", "black", "text color")
	flag.IntVar(&options.width, "width", 1000, "canvas width")
	flag.IntVar(&options.height, "height", 1500, "canvas height")
	flag.Float64Var(&options.ts, "ts", 0, "text size (0 for autoscale)")
	flag.Parse()
	go showfonts("showfonts", flag.Args(), options)
	app.Main()
}

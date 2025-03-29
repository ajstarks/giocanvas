// elections: show election results on a state grid
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/io/event"
	"gioui.org/io/input"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/unit"
	gc "github.com/ajstarks/giocanvas"
)

// Data file structure
type egrid struct {
	name       string
	party      string
	row        int
	col        int
	population int
}

// One election "frame"
type election struct {
	title    string
	min, max int
	data     []egrid
}

type options struct {
	width, height                                   int
	top, left, rowsize, colsize                     float64
	bgcolor, textcolor, shape, sansfont, symbolfont string
}

// character map for the Stateface font
var statemap = map[string]string{
	"AL": "B",
	"AK": "A",
	"AZ": "D",
	"AR": "C",
	"CA": "E",
	"CO": "F",
	"CT": "G",
	"DE": "H",
	"FL": "I",
	"GA": "J",
	"HI": "K",
	"ID": "M",
	"IL": "N",
	"IN": "O",
	"IA": "L",
	"KS": "P",
	"KY": "Q",
	"LA": "R",
	"ME": "U",
	"MD": "T",
	"MA": "S",
	"MI": "V",
	"MN": "W",
	"MS": "Y",
	"MO": "X",
	"MT": "Z",
	"NE": "c",
	"NV": "g",
	"NH": "d",
	"NJ": "e",
	"NM": "f",
	"NY": "h",
	"NC": "a",
	"ND": "b",
	"OH": "i",
	"OK": "j",
	"OR": "k",
	"PA": "l",
	"RI": "m",
	"SC": "n",
	"SD": "o",
	"TN": "p",
	"TX": "q",
	"UT": "r",
	"VT": "t",
	"VA": "s",
	"WA": "u",
	"WV": "w",
	"WI": "v",
	"WY": "x",
	"DC": "y",
}

var fontmap = map[string]string{"sans": "Go-Regular", "symbol": "stateface"}
var partyColors = map[string]string{"r": "red", "d": "blue", "i": "gray", "w": "peru", "dr": "purple", "f": "green"}

// maprange maps one range into another
func maprange(value, low1, high1, low2, high2 float64) float64 {
	return low2 + (high2-low2)*(value-low1)/(high1-low1)
}

// area computes the area of a circle
func area(v float64) float64 {
	return math.Sqrt((v / math.Pi)) * 2
}

// atoi converts a string to an integer
func atoi(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

// readData reads election data into the data structure
func readData(r io.Reader) (election, error) {
	var d egrid
	var data []egrid
	var min, max int
	title := ""
	scanner := bufio.NewScanner(r)
	min, max = math.MaxInt32, -math.MaxInt32
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) == 0 { // skip blank lines
			continue
		}
		if t[0] == '#' && len(t) > 2 { // get the title
			title = t[2:]
			continue
		}
		fields := strings.Split(t, "\t")
		if len(fields) < 5 { // skip incomplete records
			continue
		}
		// name,col,row,party,population
		d.name = fields[0]
		d.col = atoi(fields[1])
		d.row = atoi(fields[2])
		d.party = fields[3]
		d.population = atoi(fields[4])
		data = append(data, d)
		// compute min, max
		if d.population > max {
			max = d.population
		}
		if d.population < min {
			min = d.population
		}
	}
	var e election
	e.title = title
	e.min = min
	e.max = max
	e.data = data
	return e, scanner.Err()
}

func process(canvas *gc.Canvas, opts options, e election) {
	beginPage(canvas, opts.bgcolor)
	fmin, fmax := float64(e.min), float64(e.max)
	amin, amax := area(fmin), area(fmax)
	sumpop := 0
	for _, d := range e.data {
		sumpop += d.population
		x := opts.left + (float64(d.row) * opts.colsize)
		y := opts.top - (float64(d.col) * opts.rowsize)
		fpop := float64(d.population)
		apop := area(fpop)

		// defaults
		txcolor := "white"
		txsize := 1.2
		font := "sans"
		name := d.name

		switch opts.shape {
		case "c": // circle
			r := maprange(apop, amin, amax, 2, opts.colsize)
			circle(canvas, x, y, r, partyColors[d.party])
		case "h": // hexagom
			r := maprange(apop, amin, amax, 2, opts.colsize)
			hexagon(canvas, x, y, r/2, partyColors[d.party])
		case "s": // square
			r := maprange(fpop, fmin, fmax, 2, opts.colsize)
			square(canvas, x, y, r, partyColors[d.party])
		case "l": // lines
			r := maprange(apop, amin, amax, 2, opts.colsize)
			polylines(canvas, x, y, r/2, 0.25, partyColors[d.party])
			txcolor = partyColors[d.party]
		case "p": // plain text
			txcolor = partyColors[d.party]
			txsize = maprange(fpop, fmin, fmax, 2, opts.colsize*0.75)
		case "g": // geographic
			txcolor = partyColors[d.party]
			name = statemap[d.name]
			font = "symbol"
			txsize = maprange(fpop, fmin, fmax, 2, opts.colsize)
		default:
			r := maprange(apop, amin, amax, 2, opts.colsize)
			circle(canvas, x, y, r, partyColors[d.party])
		}
		ctext(canvas, x, y-0.5, txsize, name, font, txcolor)
	}
	showtitle(canvas, e.title, sumpop, opts.top+15, opts.textcolor)
	endPage(canvas)
}

func million(n int) string {
	s := strconv.FormatInt(int64(n), 10)
	p := len(s)
	return s[0:p-6] + "," + s[p-6:p-3] + "," + s[p-3:p]
}

func partycand(s, def string) (string, string) {
	var party, cand string
	f := strings.Split(s, ":")
	if len(f) > 1 {
		party = f[1]
		cand = f[0]
	} else {
		party = def
		cand = s
	}
	return party, cand
}

// showtitle shows the title and subhead
func showtitle(canvas *gc.Canvas, s string, pop int, top float64, textcolor string) {
	fields := strings.Fields(s) // year, democratic, republican, third-party (optional)
	if len(fields) < 2 {
		return
	}
	suby := top - 7
	ctext(canvas, 50, top, 3.6, fields[0]+" US Presidential Election", "sans", textcolor)
	ctext(canvas, 85, 5, 1.5, "Population: "+million(pop), "sans", textcolor)

	var party string
	var cand string
	if len(fields) > 1 {
		party, cand = partycand(fields[1], "d")
		legend(canvas, 20, suby, 2.0, cand, partyColors[party], textcolor)
	}
	if len(fields) > 2 {
		party, cand = partycand(fields[2], "r")
		legend(canvas, 80, suby, 2.0, cand, partyColors[party], textcolor)
	}
	if len(fields) > 3 {
		party, cand = partycand(fields[3], "i")
		legend(canvas, 50, suby, 2.0, cand, partyColors[party], textcolor)
	}
}

// circle makes a circle
func circle(canvas *gc.Canvas, x, y, r float64, color string) {
	cx, cy, cr := float32(x), float32(y), float32(r)
	canvas.Circle(cx, cy, cr/2, gc.ColorLookup(color))
}

// ctext makes centered text
func ctext(canvas *gc.Canvas, x, y, size float64, s string, fontname string, color string) {
	canvas.Theme.Face = font.Typeface(fontmap[fontname])
	tx, ty, ts := float32(x), float32(y), float32(size)
	canvas.CText(tx, ty, ts, s, gc.ColorLookup(color))
}

// ltext makes left-aligned text
func ltext(canvas *gc.Canvas, x, y, size float64, s string, color string) {
	tx, ty, ts := float32(x), float32(y), float32(size)
	canvas.Text(tx, ty, ts, s, gc.ColorLookup(color))
}

// square makes a square centered ar (x,y), width w.
func square(canvas *gc.Canvas, x, y, w float64, color string) {
	canvas.Square(float32(x), float32(y), float32(w), gc.ColorLookup(color))
}

// hexagon makes a filled hexagon centered at (cx, cy), size is the subscribed circle radius r
func hexagon(canvas *gc.Canvas, cx, cy, r float64, color string) {
	// construct a polygon with points at these angles
	angles := []float32{30, 90, 150, 210, 270, 330}
	px := make([]float32, len(angles))
	py := make([]float32, len(angles))
	x, y, rad := float32(cx), float32(cy), float32(r)
	for i, a := range angles {
		px[i], py[i] = canvas.PolarDegrees(x, y, rad, a)
	}
	canvas.Polygon(px, py, gc.ColorLookup(color))
}

// polylines makes a outlined hexagon, centered at (cx, cy), size is the subscribed circle radius r
func polylines(canvas *gc.Canvas, cx, cy, r, lw float64, color string) {
	// construct a polygon with points at these angles
	angles := []float32{30, 90, 150, 210, 270, 330} // square: []float64{45, 135, 225, 315}
	px := make([]float32, len(angles))
	py := make([]float32, len(angles))
	x, y, rad := float32(cx), float32(cy), float32(r)
	linewidth := float32(lw)
	for i, a := range angles {
		px[i], py[i] = canvas.PolarDegrees(x, y, rad, a)
	}
	lx := len(px) - 1
	for i := 0; i < lx; i++ {
		canvas.Line(px[i], py[i], px[i+1], py[i+1], linewidth, gc.ColorLookup(color))
	}
	canvas.Line(px[0], py[0], px[lx], py[lx], linewidth, gc.ColorLookup(color))
}

// legend makes the subtitle
func legend(canvas *gc.Canvas, x, y, ts float64, s string, color, textcolor string) {
	ltext(canvas, x, y, ts, s, textcolor)
	circle(canvas, x-ts, y+ts/3, ts/2, color)
}

// beginPage starts a page
func beginPage(canvas *gc.Canvas, bgcolor string) {
	canvas.Rect(50, 50, 100, 100, gc.ColorLookup(bgcolor))
}

// endPage ends a page
func endPage(canvas *gc.Canvas) {
	ctext(canvas, 50, 5, 1.5, "The area of a circle denotes state population: source U.S. Census", "sans", "gray")
}

var pressed bool
var electionNumber int

func kbpointer(q input.Source, context *op.Ops, ns int) {
	nev := 0
	for {
		e, ok := q.Event(
			key.Filter{Optional: key.ModCtrl},
			pointer.Filter{Target: &pressed, Kinds: pointer.Press | pointer.Move | pointer.Release | pointer.Scroll},
		)
		if !ok {
			break
		}
		switch e := e.(type) {

		case key.Event:
			switch e.State {
			case key.Press:
				switch e.Name {
				// emacs bindings
				case "A", "1": // first slide
					if e.Modifiers == 0 || e.Modifiers == key.ModCtrl {
						electionNumber = 0
					}
				case "E": // last slide
					if e.Modifiers == 0 || e.Modifiers == key.ModCtrl {
						electionNumber = ns
					}
				case "B": // back a slide
					if e.Modifiers == 0 || e.Modifiers == key.ModCtrl {
						electionNumber--
					}
				case "F": // forward a slide
					if e.Modifiers == 0 || e.Modifiers == key.ModCtrl {
						electionNumber++
					}
				case "P": // previous slide
					if e.Modifiers == 0 || e.Modifiers == key.ModCtrl {
						electionNumber--
					}
				case "N": // next slide
					if e.Modifiers == 0 || e.Modifiers == key.ModCtrl {
						electionNumber++
					}
				case "^", "⇱": // first slide
					electionNumber = 0
				case "$", "⇲": // last slide
					electionNumber = ns
				case key.NameSpace, "⏎":
					if e.Modifiers == 0 {
						electionNumber++
					} else {
						electionNumber--
					}
				case key.NameRightArrow, key.NamePageDown, key.NameDownArrow, "K":
					electionNumber++
				case key.NameLeftArrow, key.NamePageUp, key.NameUpArrow, "J":
					electionNumber--
				case key.NameEscape, "Q":
					os.Exit(0)
				}
			}

		case pointer.Event:
			switch e.Kind {
			case pointer.Scroll:
				nev++
				if e.Scroll.Y > 0 && nev == 2 {
					electionNumber--
				}
				if e.Scroll.Y == 0 && nev == 2 {
					electionNumber++
				}
			case pointer.Press:
				switch e.Buttons {
				case pointer.ButtonPrimary:
					electionNumber++
				case pointer.ButtonSecondary:
					electionNumber--
				case pointer.ButtonTertiary:
					electionNumber = 0
				}
			}
		}
	}
	event.Op(context, &pressed)
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

func elect(title string, opts options, elections []election) {
	fc, err := loadfonts([]string{"Go-Bold.ttf", "stateface.ttf"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	cw, ch := float32(opts.width), float32(opts.height)
	w := &app.Window{}
	w.Option(app.Title(title), app.Size(unit.Dp(cw), unit.Dp(ch)))
	ne := len(elections) - 1

	for {
		e := w.Event()
		switch e := e.(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			canvas := gc.NewCanvasFonts(float32(e.Size.X), float32(e.Size.Y), fc, app.FrameEvent{})
			if electionNumber > ne {
				electionNumber = 0
			}
			if electionNumber < 0 {
				electionNumber = ne
			}
			process(canvas, opts, elections[electionNumber])
			kbpointer(e.Source, canvas.Context.Ops, ne)
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {

	// command line options
	var opts options
	flag.Float64Var(&opts.top, "top", 75, "map top value (canvas %)")
	flag.Float64Var(&opts.left, "left", 15, "map left value (canvas %)")
	flag.Float64Var(&opts.rowsize, "rowsize", 9, "rowsize (canvas %)")
	flag.Float64Var(&opts.colsize, "colsize", 7, "column size (canvas %)")
	flag.IntVar(&opts.width, "width", 1200, "canvas width")
	flag.IntVar(&opts.height, "height", 900, "canvas height")
	flag.StringVar(&opts.sansfont, "sans", "Go-Regular", "sans font")
	flag.StringVar(&opts.symbolfont, "symbol", "stateface", "symbol font")
	flag.StringVar(&opts.bgcolor, "bgcolor", "black", "background color")
	flag.StringVar(&opts.textcolor, "textcolor", "white", "text color")
	flag.StringVar(&opts.shape, "shape", "c", "shape for states:\n\"c\": circle,\n\"h\": hexagon,\n\"s\": square\n\"l\": line\n\"g\": geographic\n\"p\": plain text")
	flag.Parse()

	// Read in the data
	var elections []election
	for _, f := range flag.Args() {
		r, err := os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		e, err := readData(r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		elections = append(elections, e)
		r.Close()
	}
	if len(elections) < 1 {
		fmt.Fprintln(os.Stderr, "no data read")
		os.Exit(1)
	}

	go elect("elections", opts, elections)
	app.Main()
}

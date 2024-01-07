// elections: show election results on a state grid
package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/unit"
	gc "github.com/ajstarks/giocanvas"
)

// Data file structure
type egrid struct {
	name       string
	party      string
	row        int
	col        int
	population int64
}

// One election "frame"
type election struct {
	title    string
	min, max int64
	data     []egrid
}

type options struct {
	width, height               int
	top, left, rowsize, colsize float64
	bgcolor, textcolor          string
}

var partyColors = map[string]string{"r": "red", "d": "blue", "i": "gray"}

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

// atoi64 converts a string to an 64-bit integer
func atoi64(s string) int64 {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return v
}

// readData reads election data into the data structure
func readData(r io.Reader) (election, error) {
	var d egrid
	var data []egrid
	var min, max int64
	title := ""
	scanner := bufio.NewScanner(r)
	min, max = math.MaxInt64, -math.MaxInt64
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
		d.population = atoi64(fields[4])
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

// process walks the data, making the visualization
func process(canvas *gc.Canvas, opts options, e election) {
	amin := area(float64(e.min))
	amax := area(float64(e.max))
	beginPage(canvas, opts.bgcolor)
	showtitle(canvas, e.title, opts.top+15, opts.textcolor)
	for _, d := range e.data {
		x := opts.left + (float64(d.row) * opts.colsize)
		y := opts.top - (float64(d.col) * opts.rowsize)
		r := maprange(area(float64(d.population)), amin, amax, 2, opts.colsize)
		circle(canvas, x, y, r, partyColors[d.party])
		ctext(canvas, x, y-0.5, 1.2, d.name, "white")
	}
	endPage(canvas)
}

// showtitle shows the title and subhead
func showtitle(canvas *gc.Canvas, s string, top float64, textcolor string) {
	fields := strings.Fields(s) // year, democratic, republican, third-party (optional)
	if len(fields) < 3 {
		return
	}
	suby := top - 7
	ctext(canvas, 50, top, 3.6, fields[0]+" US Presidential Election", textcolor)
	legend(canvas, 20, suby, 2.0, fields[1], partyColors["d"], textcolor)
	legend(canvas, 80, suby, 2.0, fields[2], partyColors["r"], textcolor)
	if len(fields) > 3 {
		legend(canvas, 50, suby, 2.0, fields[3], partyColors["i"], textcolor)
	}
}

// circle makes a circle
func circle(canvas *gc.Canvas, x, y, r float64, color string) {
	cx, cy, cr := float32(x), float32(y), float32(r)
	canvas.Circle(cx, cy, cr/2, gc.ColorLookup(color))
}

// ctext makes centered text
func ctext(canvas *gc.Canvas, x, y, size float64, s string, color string) {
	tx, ty, ts := float32(x), float32(y), float32(size)
	canvas.CText(tx, ty, ts, s, gc.ColorLookup(color))
}

// ltext makes left-aligned text
func ltext(canvas *gc.Canvas, x, y, size float64, s string, color string) {
	tx, ty, ts := float32(x), float32(y), float32(size)
	canvas.Text(tx, ty, ts, s, gc.ColorLookup(color))
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
	ctext(canvas, 50, 5, 1.5, "The area of a circle denotes state population: source U.S. Census", "gray")
}

var pressed bool
var electionNumber int

func kbpointer(q event.Queue, ns int) {
	nev := 0
	for _, ev := range q.Events(pressed) {
		if k, ok := ev.(key.Event); ok {
			switch k.State {
			case key.Press:
				switch k.Name {
				// emacs bindings
				case "A", "1": // first slide
					if k.Modifiers == 0 || k.Modifiers == key.ModCtrl {
						electionNumber = 0
					}
				case "E": // last slide
					if k.Modifiers == 0 || k.Modifiers == key.ModCtrl {
						electionNumber = ns
					}
				case "B": // back a slide
					if k.Modifiers == 0 || k.Modifiers == key.ModCtrl {
						electionNumber--
					}
				case "F": // forward a slide
					if k.Modifiers == 0 || k.Modifiers == key.ModCtrl {
						electionNumber++
					}
				case "P": // previous slide
					if k.Modifiers == 0 || k.Modifiers == key.ModCtrl {
						electionNumber--
					}
				case "N": // next slide
					if k.Modifiers == 0 || k.Modifiers == key.ModCtrl {
						electionNumber++
					}
				case "^", "⇱": // first slide
					electionNumber = 0
				case "$", "⇲": // last slide
					electionNumber = ns
				case key.NameSpace, "⏎":
					if k.Modifiers == 0 {
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
		}
		if p, ok := ev.(pointer.Event); ok {
			switch p.Kind {
			case pointer.Scroll:
				nev++
				if p.Scroll.Y > 0 && nev == 2 {
					electionNumber--
				}
				if p.Scroll.Y == 0 && nev == 2 {
					electionNumber++
				}
			case pointer.Press:
				switch p.Buttons {
				case pointer.ButtonPrimary:
					electionNumber++
				case pointer.ButtonSecondary:
					electionNumber--
				case pointer.ButtonTertiary:
					electionNumber = 0
				}
				pressed = true
			}
		}
	}

}

func elect(title string, opts options, elections []election) {
	cw, ch := float32(opts.width), float32(opts.height)
	win := app.NewWindow(app.Title(title), app.Size(unit.Dp(cw), unit.Dp(ch)))
	ne := len(elections) - 1

	for {
		e := win.NextEvent()
		switch e := e.(type) {
		case system.DestroyEvent:
			os.Exit(0)
		case system.FrameEvent:
			canvas := gc.NewCanvas(float32(e.Size.X), float32(e.Size.Y), system.FrameEvent{})
			key.InputOp{Tag: pressed}.Add(canvas.Context.Ops)
			pointer.InputOp{
				Tag:          pressed,
				Grab:         false,
				Kinds:        pointer.Press | pointer.Scroll,
				ScrollBounds: image.Rect(0, 0, e.Size.X, e.Size.Y)}.Add(canvas.Context.Ops)
			if electionNumber > ne {
				electionNumber = 0
			}
			if electionNumber < 0 {
				electionNumber = ne
			}
			process(canvas, opts, elections[electionNumber])
			kbpointer(e.Queue, ne)
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {

	// command line options
	var opts options
	flag.Float64Var(&opts.top, "top", 75, "map top value (canvas %)")
	flag.Float64Var(&opts.left, "left", 7, "map left value (canvas %)")
	flag.Float64Var(&opts.rowsize, "rowsize", 9, "rowsize (canvas %)")
	flag.Float64Var(&opts.colsize, "colsize", 7, "column size (canvas %)")
	flag.IntVar(&opts.width, "width", 1200, "canvas width")
	flag.IntVar(&opts.height, "height", 900, "canvas height")
	flag.StringVar(&opts.bgcolor, "bgcolor", "black", "background color")
	flag.StringVar(&opts.textcolor, "textcolor", "white", "text color")
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

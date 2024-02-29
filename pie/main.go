// pie chart
// read data from files with this format:
// one line per item, fields (name, value, color; tab-separated)
// lines beginning with '#' are the title
// for example:
//
// # Desktop Browser Market Share 2021-09
// Chrome	67.17	red
// Edge	9.33	green
// Firefox	7.87	orange
// Safari	9.63	blue
// Other	5.99	gray

package main

import (
	"bufio"
	"flag"
	"fmt"
	"image/color"
	"io"
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

const fullcircle = 3.14159265358979323846264338327950288419716939937510582097494459 * 2

type piedata struct {
	name  string
	value float64
	color string
}

func datasum(data []piedata) float64 {
	sum := 0.0
	for _, d := range data {
		sum += d.value
	}
	return sum
}

// readpie reads tab separated values
func readpie(filename string) ([]piedata, string, error) {
	var d piedata
	var data []piedata
	var err error
	var title string

	r, err := os.Open(filename)
	if err != nil {
		return data, "", err
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		t := scanner.Text()
		if len(t) == 0 { // skip blank lines
			continue
		}
		if t[0] == '#' && len(t) > 2 { // process titles
			title = strings.TrimSpace(t[1:])
			continue
		}
		fields := strings.Split(t, "\t")
		if len(fields) < 3 {
			continue
		}
		d.name = fields[0]
		d.color = fields[2]
		d.value, err = strconv.ParseFloat(fields[1], 64)
		if err != nil {
			d.value = 0
		}
		data = append(data, d)
	}
	r.Close()
	err = scanner.Err()
	return data, title, err
}

// piechart draws the piechart
func piechart(canvas *giocanvas.Canvas, x, y, r float32, data []piedata) {
	sum := datasum(data)
	a1 := 0.0
	labelr := r + 10
	ts := r / 10
	for _, d := range data {
		color := giocanvas.ColorLookup(d.color)
		p := (d.value / sum)
		angle := p * fullcircle
		a2 := a1 + angle
		mid := fullcircle - (a1 + (a2-a1)/2)
		canvas.Arc(x, y, r, a1, a2, color)
		tx, ty := canvas.Polar(x, y, labelr, float32(mid))
		lx, ly := canvas.Polar(x, y, labelr-ts, float32(mid))
		canvas.CText(tx, ty, ts, fmt.Sprintf("%s (%.2f%%)", d.name, p*100), color)
		canvas.Line(x, y, lx, ly, 0.1, color)
		a1 = a2
	}
}

var pressed bool
var pieNumber int

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
						pieNumber = 0
					}
				case "E": // last slide
					if e.Modifiers == 0 || e.Modifiers == key.ModCtrl {
						pieNumber = ns
					}
				case "B": // back a slide
					if e.Modifiers == 0 || e.Modifiers == key.ModCtrl {
						pieNumber--
					}
				case "F": // forward a slide
					if e.Modifiers == 0 || e.Modifiers == key.ModCtrl {
						pieNumber++
					}
				case "P": // previous slide
					if e.Modifiers == 0 || e.Modifiers == key.ModCtrl {
						pieNumber--
					}
				case "N": // next slide
					if e.Modifiers == 0 || e.Modifiers == key.ModCtrl {
						pieNumber++
					}
				case "^", "⇱": // first slide
					pieNumber = 0
				case "$", "⇲": // last slide
					pieNumber = ns
				case key.NameSpace, "⏎":
					if e.Modifiers == 0 {
						pieNumber++
					} else {
						pieNumber--
					}
				case key.NameRightArrow, key.NamePageDown, key.NameDownArrow, "K":
					pieNumber++
				case key.NameLeftArrow, key.NamePageUp, key.NameUpArrow, "J":
					pieNumber--
				case key.NameEscape, "Q":
					os.Exit(0)
				}
			}

		case pointer.Event:
			switch e.Kind {
			case pointer.Scroll:
				nev++
				if e.Scroll.Y > 0 && nev == 2 {
					pieNumber--
				}
				if e.Scroll.Y == 0 && nev == 2 {
					pieNumber++
				}
			case pointer.Press:
				switch e.Buttons {
				case pointer.ButtonPrimary:
					pieNumber++
				case pointer.ButtonSecondary:
					pieNumber--
				case pointer.ButtonTertiary:
					pieNumber = 0
				}
			}
		}
	}
	event.Op(context, &pressed)
}

// pie is the app
func pie(w *app.Window, files []string) error {
	var err error
	var data [][]piedata
	var title []string
	var nfiles int
	if len(files) == 0 { // if no files are specified, load default data
		data = make([][]piedata, 1)
		title = make([]string, 1)
		title[0] = "Browser Market Share, 2021-09"
		data[0] = []piedata{
			{name: "Chrome", value: 67.17, color: "#2171b5"},
			{name: "Edge", value: 9.33, color: "#4292c6"},
			{name: "Firefox", value: 7.87, color: "#6baed6"},
			{name: "Safari", value: 9.63, color: "#9ecae1"},
			{name: "Other", value: 5.99, color: "#c6dbef"},
		}
	} else { // if you have files, read and load data, skipping bad files
		data = make([][]piedata, len(files))
		title = make([]string, len(files))

		nfiles = 0
		for i := 0; i < len(files); i++ {
			d, t, ferr := readpie(files[i])
			if ferr != nil {
				fmt.Fprintf(os.Stderr, "%v\n", ferr)
				continue
			}
			data[nfiles] = d
			title[nfiles] = t
			nfiles++
		}
	}
	nf := nfiles - 1
	pieNumber = 0
	var (
		piesize float32 = 25
		top     float32 = 95
		bottom  float32 = 100 - top
	)
	for {
		e := w.NextEvent()
		switch e := e.(type) {
		case app.DestroyEvent:
			return err
		case app.FrameEvent:
			w, h := float32(e.Size.X), float32(e.Size.Y)
			canvas := giocanvas.NewCanvas(w, h, app.FrameEvent{})
			if pieNumber >= nf {
				pieNumber = 0
			}
			if pieNumber < 0 {
				pieNumber = nf
			}
			canvas.Background(color.NRGBA{0, 0, 0, 255})
			canvas.CText(50, top, 4, title[pieNumber], color.NRGBA{240, 240, 240, 255})
			canvas.CText(50, bottom, 2, "Source: StatCounter", color.NRGBA{150, 150, 150, 255})
			piechart(canvas, 50, 50, piesize, data[pieNumber])
			kbpointer(e.Source, canvas.Context.Ops, nfiles)
			e.Frame(canvas.Context.Ops)
		}
	}
	return err
}

func main() {
	var cw, ch int

	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.Parse()

	width := float32(cw)
	height := float32(ch)

	go func() {
		w := app.NewWindow(app.Title("pie"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := pie(w, flag.Args()); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

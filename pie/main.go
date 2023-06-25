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
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
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

func animate(canvas *giocanvas.Canvas, duration time.Duration) {
	t := canvas.Context.Now.Add(duration)
	op.InvalidateOp{At: t}.Add(canvas.Context.Ops)
}

// pie is the app
func pie(w *app.Window, width, height float32, duration time.Duration, files []string) error {
	var err error
	var data [][]piedata
	var title []string
	var nfiles int
	if len(files) == 0 { // if no files are specified, load default data
		data = make([][]piedata, 1)
		title = make([]string, 1)
		title[0] = "Desktop Browser Market Share, 2021-09"
		data[0] = []piedata{
			{name: "Chrome", value: 67.17, color: "red"},
			{name: "Edge", value: 9.33, color: "green"},
			{name: "Firefox", value: 7.87, color: "orange"},
			{name: "Safari", value: 9.63, color: "blue"},
			{name: "Other", value: 5.99, color: "gray"},
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

	i := 0
	for {
		e := <-w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			return err
		case system.FrameEvent:
			canvas := giocanvas.NewCanvas(float32(e.Size.X), float32(e.Size.Y), system.FrameEvent{})
			canvas.CText(50, 92, 4, title[i], color.NRGBA{20, 20, 20, 255})
			canvas.CText(50, 5, 2, "Source: StatCounter", color.NRGBA{150, 150, 150, 255})
			piechart(canvas, 50, 50, 25, data[i])
			if len(files) > 0 { // if you have multiple files, cycle through them
				i++
				if i == nfiles {
					i = 0
				}
				animate(canvas, duration)
				// time.Sleep(duration)
			}
			e.Frame(canvas.Context.Ops)
		}
	}
	return err
}

func main() {
	var cw, ch int
	var interval time.Duration
	ddur, _ := time.ParseDuration("1s")
	flag.IntVar(&cw, "width", 1000, "canvas width")
	flag.IntVar(&ch, "height", 1000, "canvas height")
	flag.DurationVar(&interval, "duration", ddur, "animation interval")
	flag.Parse()

	width := float32(cw)
	height := float32(ch)

	go func() {
		w := app.NewWindow(app.Title("pie"), app.Size(unit.Dp(width), unit.Dp(height)))
		if err := pie(w, width, height, interval, flag.Args()); err != nil {
			io.WriteString(os.Stderr, "Cannot create the window\n")
			os.Exit(1)
		}
		os.Exit(0)
	}()
	app.Main()
}

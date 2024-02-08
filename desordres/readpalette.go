package main

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"os"
	"strconv"
	"strings"
)

type spalette map[string][]string
type rgbpalette map[string][]color.NRGBA

var palette = spalette{
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
	"pen-n-paper":            {"#e4dbba", "#a4929a", "#4f3a54", "#260d1c"},
}

func rgb(x uint32) (uint8, uint8, uint8) {
	r := x & 0xff0000 >> 16
	g := x & 0x00ff00 >> 8
	b := x & 0x0000ff
	return uint8(r), uint8(g), uint8(b)
}

func ReadString(r io.Reader) (spalette, error) {
	scanner := bufio.NewScanner(r)
	p := make(spalette)
	for scanner.Scan() {
		args := strings.Fields(scanner.Text())
		l := len(args)
		if l < 2 {
			continue
		}
		name := args[0]
		p[name] = args[1:]
	}
	return p, scanner.Err()
}

func ReadRGB(r io.Reader) (rgbpalette, error) {

	palette, err := ReadString(r)
	if err != nil {
		return nil, err
	}

	rp := make(rgbpalette)
	for name, value := range palette {
		colors := make([]color.NRGBA, len(value))
		i := 0
		for _, c := range value {
			if len(c) != 7 {
				continue // must be #nnnnnn
			}
			x, err := strconv.ParseUint(c[1:], 16, 32)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}
			r, g, b := rgb(uint32(x))
			colors[i] = color.NRGBA{R: r, G: g, B: b, A: 0xff}
			i++
		}
		rp[name] = colors
	}
	return rp, nil
}

func LoadPalette(filename string) (spalette, error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return ReadString(r)
}

func LoadRGBPalette(filename string) (rgbpalette, error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return ReadRGB(r)
}

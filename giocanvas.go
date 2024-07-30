// Package giocanvas is a 2D canvas API built on gio
package giocanvas

import (
	"image/color"

	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"
)

// Canvas defines the Gio canvas
type Canvas struct {
	Width, Height float32
	Theme         *material.Theme
	TextColor     color.NRGBA
	Context       layout.Context
}

// setupCanvas sets up common canvas items
func setupCanvas(width, height float32, e app.FrameEvent, f []font.FontFace) *Canvas {
	canvas := new(Canvas)
	canvas.Width = width
	canvas.Height = height
	canvas.TextColor = color.NRGBA{0, 0, 0, 255}
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(f))
	canvas.Theme = theme
	canvas.Context = app.NewContext(new(op.Ops), e)
	iw, ih := int(width), int(height)
	canvas.Context.Constraints.Min.X = iw
	canvas.Context.Constraints.Min.Y = ih
	canvas.Context.Constraints.Max.X = iw
	canvas.Context.Constraints.Max.Y = ih
	return canvas
}

// NewCanvas initializes a Canvas using the default font set
func NewCanvas(width, height float32, e app.FrameEvent) *Canvas {
	return setupCanvas(width, height, e, gofont.Regular())
}

// NewCanvasFonts initializes a canvas with a specified set of set of fonts
func NewCanvasFonts(width, height float32, fonts []font.FontFace, e app.FrameEvent) *Canvas {
	return setupCanvas(width, height, e, fonts)
}

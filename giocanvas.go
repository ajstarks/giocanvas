// Package giocanvas is a 2D canvas API built on gio
package giocanvas

import (
	"image/color"

	"gioui.org/app"
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

// NewCanvas initializes a Canvas
func NewCanvas(width, height float32, e app.FrameEvent) *Canvas {
	canvas := new(Canvas)
	canvas.Width = width
	canvas.Height = height
	canvas.TextColor = color.NRGBA{0, 0, 0, 255}
	theme := material.NewTheme()
	theme.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(gofont.Regular()))
	canvas.Theme = theme
	canvas.Context = app.NewContext(new(op.Ops), e)
	iw, ih := int(width), int(height)
	canvas.Context.Constraints.Min.X = iw
	canvas.Context.Constraints.Min.Y = ih
	canvas.Context.Constraints.Max.X = iw
	canvas.Context.Constraints.Max.Y = ih
	return canvas
}

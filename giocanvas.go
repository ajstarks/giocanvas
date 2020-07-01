// Package giocanvas is a 2D canvas API built on gio
package giocanvas

import (
	"image/color"

	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
)

// Canvas defines the Gio canvas
type Canvas struct {
	Width, Height float32
	TextColor     color.RGBA
	Context       layout.Context
}

// NewCanvas initializes a Canvas
func NewCanvas(width, height float32, e system.FrameEvent) *Canvas {
	canvas := new(Canvas)
	canvas.Width = width
	canvas.Height = height
	canvas.TextColor = color.RGBA{0, 0, 0, 255}
	canvas.Context = layout.NewContext(new(op.Ops), e)
	iw, ih := int(width), int(height)
	canvas.Context.Constraints.Min.X = iw
	canvas.Context.Constraints.Min.Y = ih
	canvas.Context.Constraints.Max.X = iw
	canvas.Context.Constraints.Max.Y = ih
	return canvas
}

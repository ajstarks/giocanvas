package main

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/unit"
	gc "github.com/ajstarks/giocanvas"
)

func hello(title string, width, height float32) {
	win := app.NewWindow(app.Title(title), app.Size(unit.Px(width), unit.Px(height)))
	for e := range win.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			canvas := gc.NewCanvas(width, height, e)
			canvas.CenterRect(50, 50, 100, 100, gc.ColorLookup("black"))
			canvas.Circle(50, 0, 50, gc.ColorLookup("blue"))
			canvas.TextMid(50, 20, 10, "hello, world", gc.ColorLookup("white"))
			canvas.CenterImage("earth.jpg", 50, 70, 1000, 1000, 30)
			e.Frame(canvas.Context.Ops)
		}
	}
}

func main() {
	go hello("hello", 1000, 1000)
	app.Main()
}

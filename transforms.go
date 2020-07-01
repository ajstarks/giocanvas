package giocanvas

import "gioui.org/op"

// Transformations

// Translate moves current location by (x,y) using percentage-based measures
func (c *Canvas) Translate(x, y float32) op.StackOp {
	x, y = dimen(x, y, c.Width, c.Height)
	return c.AbsTranslate(x, y)
}

// Rotate around (x,y) by angle (radians) using percentage-based measures
func (c *Canvas) Rotate(x, y, angle float32) op.StackOp {
	x, y = dimen(x, y, c.Width, c.Height)
	return c.AbsRotate(x, y, angle)
}

// Scale centered at (x,y) by factor using percentage-based measures
func (c *Canvas) Scale(x, y, factor float32) op.StackOp {
	x, y = dimen(x, y, c.Width, c.Height)
	return c.AbsScale(x, y, factor)
}

// Shear the object centered at (x,y) using x-angle and y-angle (radians) using percentage-based measures
func (c *Canvas) Shear(x, y, ax, ay float32) op.StackOp {
	x, y = dimen(x, y, c.Width, c.Height)
	return c.AbsShear(x, y, ax, ay)
}

// EndTransform ends a transformation
func EndTransform(stack op.StackOp) {
	stack.Pop()
}

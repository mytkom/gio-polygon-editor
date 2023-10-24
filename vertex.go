package main

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

const vertexRadius = 5
const vertexHoverRadius = vertexRadius + 3

type EdgeConstraint int

const(
    None EdgeConstraint = iota
    Vertical
    Horizontal
)

type Vertex struct {
    Point f32.Point
    Hovered bool
    EdgeConstraint EdgeConstraint
    next *Vertex
    previous *Vertex
}

func (v *Vertex) Layout(ops *op.Ops, c color.NRGBA) {
    circle := v.GetEllipse().Op(ops)
    paint.FillShape(ops, c, circle)
}

func (v *Vertex) MoveBy(x float32, y float32, gtx *layout.Context) {
    if !isPointInWindow(v.Point.X + x, v.Point.Y + y, gtx) {
        return 
    }

    v.Point.X += x
    v.Point.Y += y

    if v.EdgeConstraint == Vertical {
        v.next.Point.X += x
    } else if v.EdgeConstraint == Horizontal {
        v.next.Point.Y += y
    }

    if v.previous.EdgeConstraint == Vertical {
        v.previous.Point.X += x
    } else if v.previous.EdgeConstraint == Horizontal {
        v.previous.Point.Y += y
    }
}

func (v *Vertex) IsClicked(point f32.Point) bool {
    xDiff := point.X - v.Point.X
    yDiff := point.Y - v.Point.Y
    if xDiff * xDiff + yDiff * yDiff <= vertexHoverRadius * vertexHoverRadius {
        return true
    } 

    return false
}

func (v* Vertex) GetEllipse() clip.Ellipse {
    return v.getClipEllipse(vertexRadius)
}

func (v* Vertex) GetHoverEllipse() clip.Ellipse {
    return v.getClipEllipse(vertexHoverRadius)
}

func (v *Vertex) getClipEllipse(radius int) clip.Ellipse {
    imagePoint := v.Point.Round()

    return clip.Ellipse{
        Min: image.Pt(imagePoint.X - radius, imagePoint.Y - radius),
        Max: image.Pt(imagePoint.X + radius, imagePoint.Y + radius),
    }   
}




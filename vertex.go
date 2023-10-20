package main

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

const vertexRadius = 5
const vertexHoverRadius = vertexRadius + 10

type Vertex struct {
    Point f32.Point
    Hovered bool
    Selected bool
}

func (v *Vertex) Layout(ops *op.Ops, c color.NRGBA) {
    circle := v.GetEllipse().Op(ops)
    paint.FillShape(ops, c, circle)

    if v.Hovered || v.Selected {
        circle = v.GetHoverEllipse().Op(ops)
        paint.FillShape(ops, HoverizeColor(c), circle)
    }
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


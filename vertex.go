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

type Vertex struct {
    Point f32.Point
    Color color.NRGBA
}

func (v *Vertex) Layout(ops *op.Ops) {
    imagePoint := v.Point.Round()
    circle := clip.Ellipse{
        Min: image.Pt(imagePoint.X - vertexRadius, imagePoint.Y - vertexRadius),
        Max: image.Pt(imagePoint.X + vertexRadius, imagePoint.Y + vertexRadius),
    }.Op(ops)
    paint.FillShape(ops, v.Color, circle)
}

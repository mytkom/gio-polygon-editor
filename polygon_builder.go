package main

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type PolygonBuilder struct {
    vertices []Vertex
    active bool
    PointerTailEndPosition f32.Point
}

func (pb *PolygonBuilder) Layout(ops *op.Ops) {
    if len(pb.vertices) == 0 {
        return
    }

    for _, vertex := range pb.vertices {
        vertex.Layout(ops)
    }

    lastVertex := pb.vertices[0]
    var path clip.Path 

    for i := 1; i < len(pb.vertices); i++ {
        path.Begin(ops)
        path.MoveTo(lastVertex.Point)
        path.LineTo(pb.vertices[i].Point)
        path.Close()

        paint.FillShape(ops, color.NRGBA{R:255, A:255},
        clip.Stroke{
            Path: path.End(),
            Width: 4,
        }.Op()) 

        lastVertex = pb.vertices[i]
    }

    pb.drawCursorTailLine(ops)
}

func (pb *PolygonBuilder) drawCursorTailLine(ops *op.Ops) {
    if !pb.active {
        return
    }

    var path clip.Path 
    path.Begin(ops)
    path.MoveTo(pb.vertices[len(pb.vertices)-1].Point)
    path.LineTo(pb.PointerTailEndPosition)
    path.Close()

    paint.FillShape(ops, color.NRGBA{R:255, A:255},
    clip.Stroke{
        Path: path.End(),
        Width: 4,
    }.Op())
}

func (pb *PolygonBuilder) SetTailEnd(newTailEndPosition f32.Point) {
   pb.PointerTailEndPosition = newTailEndPosition
}

func (pb *PolygonBuilder) AddVertex(pointerPosition f32.Point) {
    vertexDefaultColor := color.NRGBA{R: 255, A: 255}

    if len(pb.vertices) == 0 {
        pb.active = true
    }

    vertex := &Vertex{
        Point: pointerPosition,
        Color: vertexDefaultColor,
    }

    pb.vertices = append(
        pb.vertices, 
        *vertex,
    )
    pb.SetTailEnd(pointerPosition)
}

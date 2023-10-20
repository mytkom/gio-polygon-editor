package main

import (
	"math"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type Edge struct {
    Vertices [2]*Vertex
    EventTag bool
}

func (e *Edge) MoveBy(x float32, y float32) {
    for _, vertex := range e.Vertices {
        vertex.Point.X += x         
        vertex.Point.Y += y
    } 
}

func (e *Edge) GetHoverClipArea(ops *op.Ops) clip.Outline {
    var path clip.Path
    
    p1 := &e.Vertices[0].Point
    p2 := &e.Vertices[1].Point
    a := (p2.X - p1.X) / (p2.Y - p1.Y)
    z := float32(10.0)
    simData := z / float32(math.Sqrt(float64(1.0 + a * a)))

    x := p1.X - simData
    y := a * (x - p1.X) + p1.Y
    v1 := &Vertex{Point: f32.Pt(x, y)}
    v1.Layout(ops, brushColor)

    x = p1.X + simData
    y = a * (x - p1.X) + p1.Y
    v2 := &Vertex{Point: f32.Pt(x, y)}
    v2.Layout(ops, brushColor)

    x = p2.X - simData
    y = a * (x - p2.X) + p2.Y
    v3 := &Vertex{Point: f32.Pt(x, y)}
    v3.Layout(ops, brushColor)

    x = p2.X + simData
    y = a * (x - p2.X) + p2.Y
    v4 := &Vertex{Point: f32.Pt(x, y)}
    v4.Layout(ops, brushColor)

    path.Begin(ops)
    path.MoveTo(v1.Point)
    path.LineTo(v2.Point)
    path.LineTo(v4.Point)
    path.LineTo(v3.Point)

    path.Close()

    outline := clip.Outline{
        Path: path.End(),
    }

    paint.FillShape(ops, brushColor, outline.Op())
    return outline
}

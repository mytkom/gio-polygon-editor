package main

import (
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
)

type Edge struct {
    Vertices [2]*Vertex
}

func (e *Edge) MoveBy(x float32, y float32, gtx *layout.Context) {
    for _, vertex := range e.Vertices {
        if !isPointInWindow(vertex.Point.X + x, vertex.Point.Y + y, gtx) {
            return
        }
    }

    v1 := e.Vertices[0]
    v2 := e.Vertices[1]

    if v2.ConstraintAfter != Vertical && v1.ConstraintBefore != Vertical {
        v1.Point.X += x
        v2.Point.X += x
    }
    if v2.ConstraintAfter != Horizontal && v1.ConstraintBefore != Horizontal {
        v1.Point.Y += y
        v2.Point.Y += y
    }
}

func (e *Edge) SetConstraint(c EdgeConstraint) {
    if c != None &&
        (e.Vertices[0].ConstraintBefore != None ||
         e.Vertices[1].ConstraintAfter != None) {
       return 
    }

    mid := e.GetMiddlePoint()
    if c == Horizontal {
        e.Vertices[0].Point.Y = mid.Y
        e.Vertices[1].Point.Y = mid.Y
    } else if c == Vertical {
        e.Vertices[0].Point.X = mid.X
        e.Vertices[1].Point.X = mid.X
    }

    e.Vertices[0].ConstraintAfter = c
    e.Vertices[1].ConstraintBefore = c
}

func (e *Edge) IsClicked(point *f32.Point) bool {
    selectionOffset := 0.15
    
    distA := calculateDistanceBetweenPoints(e.Vertices[0].Point, *point)
    distB := calculateDistanceBetweenPoints(e.Vertices[1].Point, *point)
    distC := calculateDistanceBetweenPoints(e.Vertices[0].Point, e.Vertices[1].Point)

    if distA + distB <= distC + selectionOffset {
       return true 
    }

    return false
}

func (e *Edge) GetMiddlePoint() f32.Point {
    a := e.Vertices[0].Point
    b := e.Vertices[1].Point
    return f32.Point{X: (a.X + b.X) / 2.0, Y: (a.Y + b.Y) / 2.0}
}

func calculateDistanceBetweenPoints(a f32.Point, b f32.Point) float64 {
    xQuad := math.Pow(float64(a.X - b.X), 2.0)
    yQuad := math.Pow(float64(a.Y - b.Y), 2.0)

    return math.Sqrt(xQuad + yQuad) 
}


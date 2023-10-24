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

    v1.Point.X += x
    v2.Point.X += x
    v1.Point.Y += y
    v2.Point.Y += y

    if v1.previous.EdgeConstraint == Vertical {
        v1.previous.Point.X += x
    } else if v1.previous.EdgeConstraint == Horizontal {
        v1.previous.Point.Y += y
    }

    if v2.EdgeConstraint == Vertical {
        v2.next.Point.X += x
    } else if v2.EdgeConstraint == Horizontal {
        v2.next.Point.Y += y
    } 
}

func (e *Edge) SetConstraint(c EdgeConstraint) {
    if c != None &&
        (e.Vertices[0].previous.EdgeConstraint == c ||
         e.Vertices[1].EdgeConstraint == c) {
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

    e.Vertices[0].EdgeConstraint = c
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


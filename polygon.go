package main

import (
	"encoding/json"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
)

type Polygon struct {
	VerticesHead *Vertex
    VerticesCount int `json:"verticesCount"`
    Color    color.NRGBA `json:"color"`
	Edges    []*Edge
}

func CreatePolygon(head *Vertex, tail *Vertex, count int, color color.NRGBA) {
    tail.next = head
    head.previous = tail

	polygon := &Polygon{
        VerticesHead: head,
        VerticesCount: count,
		Color:    color,
	}
	polygon.CreateEdges()
    polygon.CalculateOffsetVectors()

    if !polygon.IsClockwise() {
        current := head
        for i := 0; i < count - 1; i++ {
            current.previous.next = current.previous.previous
            current.previous.previous = current
            current = current.next
        }
        current.next.next = current.next.previous
        current.next.previous = current
    }

	polygons = append(
		polygons,
		polygon,
	)
}

func (p *Polygon) IsClockwise() bool {
    current := p.VerticesHead
    sum := float32(0.0)
    for i := 0; i < p.VerticesCount; i++ {
        next := current.next
        sum += (next.Point.X - current.Point.X) * (next.Point.Y + current.Point.Y)
        current = next
    }

    if sum < 0 {
        return true
    }
    return false
}

func (p *Polygon) CreateEdges() {
	edges := make([]*Edge, p.VerticesCount)
    var next *Vertex
    current := p.VerticesHead
    for i := 0; i < p.VerticesCount; i++ {
        next = current.next
        edges[i] = &Edge{Vertices: [2]*Vertex{current, next}}
        current = next
	}

	p.Edges = edges
}

func (p *Polygon) IsClicked(point f32.Point) bool {
	isInside := false

    for _, edge := range p.Edges {
		vi, vj := edge.Vertices[0].Point, edge.Vertices[1].Point

		if (vi.Y > point.Y) != (vj.Y > point.Y) &&
			point.X < (vj.X-vi.X)*(point.Y-vi.Y)/(vj.Y-vi.Y)+vi.X {
			isInside = !isInside
		}
	}

	return isInside
}

func (p *Polygon) Layout(gtx *layout.Context) {
	drawPolygonFromVertices(p.VerticesHead, p.VerticesCount, gtx.Ops, &p.Color)

    if offsetPolygonFeatureEnabled {
        p.DrawOffsetPolygon(gtx.Ops)
    }
}

func (p *Polygon) DrawOffsetPolygon(ops *op.Ops) {
    current := p.VerticesHead
    for i := 0; i < p.VerticesCount; i++ {
        vertex := Vertex{Point: current.Point.Add(current.offsetVector.Mul(float32(polygonOffset)))}
        vertex.Layout(ops, offsetColor)
        if current.next == nil {
            break
        }
        applicationPainter.DrawLine(
            current.Point.Add(current.offsetVector.Mul(float32(polygonOffset))),
            current.next.Point.Add(current.next.offsetVector.Mul(float32(polygonOffset))),
            offsetColor,
            ops,
        )
        current = current.next
    }
}

func (p *Polygon) CalculateOffsetVectors() {
    current := p.VerticesHead
    for i := 0; i < p.VerticesCount; i++ {
        sub := current.next.Point.Sub(current.Point)
        n1 := normalizeVector(f32.Point{X: sub.Y, Y: -sub.X})
        sub = current.Point.Sub(current.previous.Point)
        n2 := normalizeVector(f32.Point{X: sub.Y, Y: -sub.X})
        vector := n1.Add(n2)
        divider := float32(math.Sqrt(float64((1.0 + n1.X * n2.X + n1.Y * n2.Y)/2.0)))
        current.offsetVector = normalizeVector(vector).Div(divider)
        current = current.next
    }
}

func normalizeVector(vector f32.Point) f32.Point {
    div := float32(math.Sqrt(float64(vector.X*vector.X + vector.Y*vector.Y)))
    return vector.Div(div)
}

func (p *Polygon) MarshalJSON() ([]byte, error) {
    vertices := make([]*Vertex, 0)
    current := p.VerticesHead
    for i := 0; i < p.VerticesCount; i++ {
        vertices = append(vertices, current)
        current = current.next
    }

    b, e := json.Marshal(vertices)
    if e != nil {
        return nil, e
    }

    return b, nil
}

func (p *Polygon) UnmarshalJSON(b []byte) error {
    var vertices []*Vertex
    e := json.Unmarshal(b, &vertices)
    if e != nil {
        return e
    }

    head := vertices[0]
    tail := vertices[len(vertices) - 1]
    for i := 0; i < len(vertices) - 1; i++ {
        vertices[i].next = vertices[i + 1] 
        vertices[i + 1].previous = vertices[i]
    }

    tail.next = head
    head.previous = tail

    p.VerticesHead = head
    p.VerticesCount = len(vertices)
    p.Color = polygonColor

	p.CreateEdges()
    p.CalculateOffsetVectors()

    if !p.IsClockwise() {
        current := p.VerticesHead
        for i := 0; i < p.VerticesCount - 1; i++ {
            current.previous.next = current.previous.previous
            current.previous.previous = current
            current = current.next
        }
        current.next.next = current.next.previous
        current.next.previous = current
    }

    return nil
}

func (p *Polygon) AppendVertexAfter(v *Vertex, point f32.Point) {
    next := v.next
    newVertex := &Vertex{previous: v, next: next, Point: point}
    v.next = newVertex
    next.previous = newVertex
    p.VerticesCount++
    v.EdgeConstraint = None

    p.CalculateOffsetVectors()
}

func (p *Polygon) DestroyVertex(v *Vertex) {
    prev := v.previous
    next := v.next

    if v == p.VerticesHead {
        p.VerticesHead = next
    }

    prev.next = next
    next.previous = prev
    prev.EdgeConstraint = None

    p.VerticesCount--
    p.CalculateOffsetVectors()
}

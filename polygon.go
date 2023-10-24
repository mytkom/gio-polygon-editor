package main

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
)

type Polygon struct {
	VerticesHead *Vertex
    VerticesCount int
	Color    color.NRGBA
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
}

func (p *Polygon) AppendVertexAfter(v *Vertex, point f32.Point) {
    next := v.next
    newVertex := &Vertex{previous: v, next: next, Point: point}
    v.next = newVertex
    next.previous = newVertex
    p.VerticesCount++
}

func (p *Polygon) DestroyVertex(v *Vertex) {
    prev := v.previous
    next := v.next

    if v == p.VerticesHead {
        p.VerticesHead = next
    }

    prev.next = next
    next.previous = prev

    p.VerticesCount--
}

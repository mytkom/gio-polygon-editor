package main

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type Polygon struct {
	VerticesHead *Vertex
    VerticesTail *Vertex
    VerticesCount int
	Color    color.NRGBA
	Edges    []*Edge
}

func CreatePolygon(head *Vertex, tail *Vertex, count int, color color.NRGBA) {
	polygon := &Polygon{
        VerticesHead: head,
        VerticesTail: tail,
        VerticesCount: count,
		Color:    color,
	}
	polygon.CreateEdges()

	polygons = append(
		polygons,
		polygon,
	)
}

func (p *Polygon) CreateEdges() {
	edges := make([]*Edge, p.VerticesCount)
    var next *Vertex
    current := p.VerticesHead
    i := 0
    for current != nil {
        next = current.next
        edges[i] = &Edge{Vertices: [2]*Vertex{current, next}}
        current = next
        i++
	}
    edges[p.VerticesCount - 1] = &Edge{Vertices: [2]*Vertex{p.VerticesTail, p.VerticesHead}}

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
	drawPolygonFromVertices(p.VerticesHead, gtx.Ops, &p.Color)
}

func (p *Polygon) AppendVertexAfter(v *Vertex, point f32.Point) {
    if v.next == nil {
        next := p.VerticesHead
        newVertex := &Vertex{next: next, Point: point}
        next.previous = newVertex
        p.VerticesHead = newVertex
    } else {
        next := v.next
        newVertex := &Vertex{previous: v, next: next, Point: point}
        v.next = newVertex
        next.previous = newVertex
    }
    p.VerticesCount += 1
}

func (p *Polygon) DestroyVertex(v *Vertex) {
    prev := v.previous
    next := v.next

    if prev == nil {
        p.VerticesHead = next
        next.previous = nil
    } else if next == nil {
        p.VerticesTail = prev
        prev.next = nil
    } else {
        prev.next = next
        next.previous = prev
    }

    p.VerticesCount -= 1
}

func drawPolygonFromVertices(head *Vertex, ops *op.Ops, color *color.NRGBA) {
	path := getPathFromVertices(head, ops, *color)
	path.Close()
	fullPath := path.End()

	paint.FillShape(ops, HoverizeColor(*color),
		clip.Outline{
			Path: fullPath,
		}.Op(),
	)

	paint.FillShape(ops, *color,
		clip.Stroke{
			Path:  fullPath,
			Width: float32(lineWidth),
		}.Op(),
	)
}

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
	Vertices []*Vertex
	Color    color.NRGBA
	Edges    []*Edge
}

func CreatePolygon(vertices []*Vertex, color color.NRGBA) {
	polygon := &Polygon{
		Vertices: vertices,
		Color:    color,
	}
	polygon.CreateEdges()

	polygons = append(
		polygons,
		polygon,
	)
}

func (p *Polygon) CreateEdges() {
	edges := make([]*Edge, len(p.Vertices))
	for i := 0; i < len(p.Vertices); i += 1 {
		j := (i + 1) % len(p.Vertices)
		edges[i] = &Edge{Vertices: [2]*Vertex{p.Vertices[i], p.Vertices[j]}}
	}
	p.Edges = edges
}

func (p *Polygon) IsClicked(point f32.Point) bool {
	numVertices := len(p.Vertices)
	isInside := false

	for i, j := 0, numVertices-1; i < numVertices; i++ {
		vi, vj := p.Vertices[i].Point, p.Vertices[j].Point

		if (vi.Y > point.Y) != (vj.Y > point.Y) &&
			point.X < (vj.X-vi.X)*(point.Y-vi.Y)/(vj.Y-vi.Y)+vi.X {
			isInside = !isInside
		}

		j = i
	}

	return isInside
}

func (p *Polygon) Layout(gtx *layout.Context) {
	drawPolygonFromVertices(p.Vertices, gtx.Ops, &p.Color)
}

func drawPolygonFromVertices(v []*Vertex, ops *op.Ops, color *color.NRGBA) {
	path := getPathFromVertices(v, ops, *color)
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

package main

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

var highLightColor = color.NRGBA{R: 255, G: 252, B: 127, A:128}

type Selectable interface {
    Destroy()
    HighLight(gtx *layout.Context)
    MoveBy(x float32, y float32, gtx *layout.Context)
}

type PolygonVertex struct {
    Polygon *Polygon
    VertexIndex int
}

func (pv *PolygonVertex) Destroy() {
    ret := make([]*Vertex, 0)

    if len(pv.Polygon.Vertices) == 3 {
        pv.Polygon.Destroy()
        return
    }

    ret = append(ret, pv.Polygon.Vertices[:pv.VertexIndex]...)
    pv.Polygon.Vertices = append(ret, pv.Polygon.Vertices[pv.VertexIndex+1:]...)
    pv.Polygon.CreateEdges()
}

func (pv *PolygonVertex) HighLight(gtx *layout.Context) {
    v := pv.Polygon.Vertices[pv.VertexIndex]
    circle := v.GetHoverEllipse().Op(gtx.Ops)
    paint.FillShape(gtx.Ops, highLightColor, circle)
}

func (pv *PolygonVertex) MoveBy(x float32, y float32, gtx *layout.Context) {
    pv.Polygon.Vertices[pv.VertexIndex].MoveBy(x, y, gtx)
}

type PolygonEdge struct {
    Polygon *Polygon
    EdgeIndex int
}

func (pe *PolygonEdge) Destroy() {
    ret := make([]*Vertex, 0)

    if len(pe.Polygon.Edges) == 3 {
        pe.Polygon.Destroy()
        return
    }

    pe.Polygon.Vertices[pe.EdgeIndex + 1 % len(pe.Polygon.Edges)].Point = pe.getEdge().GetMiddlePoint()
    ret = append(ret, pe.Polygon.Vertices[:pe.EdgeIndex]...)

    pe.Polygon.Vertices = append(ret, pe.Polygon.Vertices[pe.EdgeIndex + 1:]...)
    pe.Polygon.CreateEdges()
}

func (pe *PolygonEdge) HighLight(gtx *layout.Context) {
    vertices := pe.getEdge().Vertices
    var path clip.Path 
    path.Begin(gtx.Ops)
    path.MoveTo(vertices[0].Point)
    path.LineTo(vertices[1].Point)
    path.Close()

    paint.FillShape(
        gtx.Ops,
        highLightColor,
        clip.Stroke{Path: path.End(), Width: float32(lineWidth + 1.0)}.Op(),
    )
}

func (pe *PolygonEdge) MoveBy(x float32, y float32, gtx *layout.Context) {
    pe.getEdge().MoveBy(x, y, gtx)
}

func (pe *PolygonEdge) getEdge() *Edge {
    return pe.Polygon.Edges[pe.EdgeIndex] 
}

func (p *Polygon) Destroy() {
    if len(polygons) == 1 {
        polygons = []*Polygon{}
        return
    }
    for i, polygon := range polygons {
        if polygon == p {
            ret := make([]*Polygon, 0)
            ret = append(ret, polygons[:i]...)
            polygons = append(ret, polygons[i + 1:]...)
        }
    }
}

func (p *Polygon) HighLight(gtx *layout.Context) {
    drawPolygonFromVertices(p.Vertices, gtx.Ops, &highLightColor)
}

func (p *Polygon) MoveBy(x float32, y float32, gtx *layout.Context) {
    for _, vertex := range p.Vertices {
        newX := vertex.Point.X + x
        newY := vertex.Point.Y + y
        if !isPointInWindow(newX, newY, gtx) {
            return
        }
    }

	for _, vertex := range p.Vertices {
		vertex.Point.X += x
		vertex.Point.Y += y
	}
}


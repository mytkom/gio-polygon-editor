package main

import (
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

var highLightColor = color.NRGBA{R: 255, G: 32, B: 32, A:180}

type Selectable interface {
    Destroy()
    HighLight(gtx *layout.Context)
    MoveBy(x float32, y float32, gtx *layout.Context)
}

type PolygonVertex struct {
    Polygon *Polygon
    Vertex *Vertex
}

func (pv *PolygonVertex) Destroy() {
    if pv.Polygon.VerticesCount == 3 {
        pv.Polygon.Destroy()
        return
    }

    pv.Polygon.DestroyVertex(pv.Vertex)
    pv.Polygon.CreateEdges()
    pv.Polygon.CalculateOffsetVectors()
}

func (pv *PolygonVertex) HighLight(gtx *layout.Context) {
    circle := pv.Vertex.GetHoverEllipse().Op(gtx.Ops)   
    paint.FillShape(gtx.Ops, highLightColor, circle)
}

func (pv *PolygonVertex) MoveBy(x float32, y float32, gtx *layout.Context) {
    pv.Vertex.MoveBy(x, y, gtx)
    pv.Polygon.CalculateOffsetVectors()
}

type PolygonEdge struct {
    Polygon *Polygon
    EdgeIndex int
}

func (pe *PolygonEdge) Destroy() {
    if len(pe.Polygon.Edges) == 3 {
        pe.Polygon.Destroy()
        return
    }

    e := pe.getEdge()
    v := e.Vertices[0]
    pe.Polygon.DestroyVertex(v)
    v.next.Point = e.GetMiddlePoint()

    pe.Polygon.CreateEdges()
    pe.Polygon.CalculateOffsetVectors()
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
        clip.Stroke{Path: path.End(), Width: float32(edgeHoverWidth)}.Op(),
    )
}

func (pe *PolygonEdge) MoveBy(x float32, y float32, gtx *layout.Context) {
    pe.getEdge().MoveBy(x, y, gtx)
    pe.Polygon.CalculateOffsetVectors()
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
    drawPolygonFromVertices(p.VerticesHead, p.VerticesCount, gtx.Ops, &highLightColor)

    var path clip.Path
    path.Begin(gtx.Ops)
    current := p.VerticesHead
    path.MoveTo(current.Point)
    for i := 1; i < p.VerticesCount; i++ {
        current = current.next
        path.LineTo(current.Point)
    }
    path.Close()
    
    paint.FillShape(
        gtx.Ops,
        HoverizeColor(highLightColor),
        clip.Outline{Path: path.End()}.Op(),
    )
}

func (p *Polygon) MoveBy(x float32, y float32, gtx *layout.Context) {
    vertex := p.VerticesHead

    for i := 0; i < p.VerticesCount; i++ {
        newX := vertex.Point.X + x
        newY := vertex.Point.Y + y
        if !isPointInWindow(newX, newY, gtx) {
            return
        }
        vertex = vertex.next
    }

    for i := 0; i < p.VerticesCount; i++ {
		vertex.Point.X += x
		vertex.Point.Y += y
        vertex = vertex.next
	}
}


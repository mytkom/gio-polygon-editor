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
}

func (pv *PolygonVertex) HighLight(gtx *layout.Context) {
    circle := pv.Vertex.GetHoverEllipse().Op(gtx.Ops)
    paint.FillShape(gtx.Ops, highLightColor, circle)
}

func (pv *PolygonVertex) MoveBy(x float32, y float32, gtx *layout.Context) {
    pv.Vertex.MoveBy(x, y, gtx)
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
    next := v.next
    prev := v.previous
    if next == nil {
        next = pe.Polygon.VerticesHead
        prev.next = nil
    } else if prev == nil {
        pe.Polygon.VerticesHead = next
        next.previous = nil
    } else {
        prev.next = next 
        next.previous = prev
    }
    
    next.Point = e.GetMiddlePoint()
    pe.Polygon.VerticesCount -= 1

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
    drawPolygonFromVertices(p.VerticesHead, gtx.Ops, &highLightColor)
}

func (p *Polygon) MoveBy(x float32, y float32, gtx *layout.Context) {
    vertex := p.VerticesHead

    for vertex != nil {
        newX := vertex.Point.X + x
        newY := vertex.Point.Y + y
        if !isPointInWindow(newX, newY, gtx) {
            return
        }
        vertex = vertex.next
    }

    vertex = p.VerticesHead
    for vertex != nil {
		vertex.Point.X += x
		vertex.Point.Y += y
        vertex = vertex.next
	}
}


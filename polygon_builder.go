package main

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type PolygonBuilder struct {
    vertices []*Vertex
    PointerTailEndPosition f32.Point
    vertexAddEventTag bool
    Color color.NRGBA
}

func (pb *PolygonBuilder) Layout(gtx *layout.Context) {
    if len(pb.vertices) == 0 {
        return
    }

    drawPathFromVertices(pb.vertices, gtx.Ops, &pb.Color)

    pb.drawCursorTailLine(gtx.Ops)
}

func (pb *PolygonBuilder) HandleEvents(gtx *layout.Context) {
    if len(pb.vertices) > 2 {
        pb.handleClosingVertexEvent(gtx) 
    }
    pb.handleAddVertexEvent(gtx)
}

func (pb *PolygonBuilder) handleAddVertexEvent(gtx *layout.Context) {
    if eventsStoppedInFrame {
        return
    } 

    for _, ev := range gtx.Events(&pb.vertexAddEventTag) {
        if x, ok := ev.(pointer.Event); ok {
            switch x.Type {
            case pointer.Press:
                polygonBuilder.addVertex(x.Position)
            case pointer.Move:
                polygonBuilder.setTailEnd(x.Position)
            }
        }
    }
}

func (pb *PolygonBuilder) handleClosingVertexEvent(gtx *layout.Context) {
    if eventsStoppedInFrame {
        return
    }

    v := pb.vertices[0]

    for _, ev := range gtx.Events(&v.Hovered) {
        if x, ok := ev.(pointer.Event); ok {
            switch x.Type {
            case pointer.Press:
                CreatePolygon(pb.vertices, color.NRGBA{R: 255, G: 255, B: 255, A: 255}) 
                pb.vertices = []*Vertex{}
                StopEventsBelow() 
            case pointer.Enter:
                pb.vertices[0].Hovered = true 
            case pointer.Leave:
                pb.vertices[0].Hovered = false
            }
        }
    }
}

func (pb *PolygonBuilder) RegisterEvents(gtx *layout.Context) {
    pb.registerAddVertexEvent(gtx)
    if len(pb.vertices) > 2 {
        pb.registerClosingVertexEvent(gtx)
    }
}

func (pb *PolygonBuilder) registerAddVertexEvent(gtx *layout.Context) {
    pointer.InputOp{
        Tag: &pb.vertexAddEventTag,
        Types: pointer.Press | pointer.Release | pointer.Move,
    }.Add(gtx.Ops)
}

func (pb *PolygonBuilder) registerClosingVertexEvent(gtx *layout.Context) {
    v := pb.vertices[0]

    defer v.GetHoverEllipse().Push(gtx.Ops).Pop()  
    pointer.InputOp{
        Tag: &v.Hovered,
        Types: pointer.Press | pointer.Enter | pointer.Leave,
    }.Add(gtx.Ops)
}

func (pb *PolygonBuilder) drawCursorTailLine(ops *op.Ops) {
    var path clip.Path 
    path.Begin(ops)
    path.MoveTo(pb.vertices[len(pb.vertices)-1].Point)
    path.LineTo(pb.PointerTailEndPosition)
    path.Close()

    paint.FillShape(ops, pb.Color,
    clip.Stroke{
        Path: path.End(),
        Width: float32(lineWidth),
    }.Op())
}

func (pb *PolygonBuilder) setTailEnd(newTailEndPosition f32.Point) {
   pb.PointerTailEndPosition = newTailEndPosition
}

func (pb *PolygonBuilder) addVertex(pointerPosition f32.Point) {
    vertex := &Vertex{
        Point: pointerPosition,
    }

    pb.vertices = append(
        pb.vertices, 
        vertex,
    )
    pb.setTailEnd(pointerPosition)
}

func drawPathFromVertices(v []*Vertex, ops *op.Ops, color *color.NRGBA) {
    path := getPathFromVertices(v, ops, *color)

    paint.FillShape(ops, *color,
    clip.Stroke{
        Path: path.End(),
        Width: float32(lineWidth),
    }.Op()) 
}

func getPathFromVertices(v []*Vertex, ops *op.Ops, color color.NRGBA) clip.Path {
    var path clip.Path

    for _, vertex := range v {
        vertex.Layout(ops, color)
    }

    path.Begin(ops)
    path.MoveTo(v[0].Point)
    for i := 1; i < len(v); i += 1 {
        path.LineTo(v[i].Point)
    }

    return path
}

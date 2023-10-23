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
    verticesHead *Vertex
    verticesTail *Vertex
    vertexCount int
    PointerTailEndPosition f32.Point
    Color color.NRGBA
    Active bool
}

func (pb *PolygonBuilder) Layout(gtx *layout.Context) {
    if pb.vertexCount == 0 {
        return
    }

    drawPathFromVertices(pb.verticesHead, gtx.Ops, &pb.Color)

    pb.drawCursorTailLine(gtx.Ops)
}

func (pb *PolygonBuilder) HandleEvents(gtx *layout.Context) {
    if pb.vertexCount > 2 {
        pb.handleClosingVertexEvent(gtx) 
    }
}

func (pb *PolygonBuilder) handleAddVertexEvent(e *pointer.Event) {
    if eventsStoppedInFrame || !pb.Active {
        return
    } 

    switch e.Type {
    case pointer.Press:
        polygonBuilder.addVertex(e.Position)
    case pointer.Move:
        polygonBuilder.setTailEnd(e.Position)
    }
}

func (pb *PolygonBuilder) handleClosingVertexEvent(gtx *layout.Context) {
    if eventsStoppedInFrame {
        return
    }

    v := pb.verticesHead

    for _, ev := range gtx.Events(&v.Hovered) {
        if x, ok := ev.(pointer.Event); ok {
            switch x.Type {
            case pointer.Press:
                CreatePolygon(pb.verticesHead, pb.verticesTail, pb.vertexCount, color.NRGBA{R: 255, G: 255, B: 255, A: 255}) 
                pb.verticesHead = nil
                pb.verticesTail = nil
                pb.vertexCount = 0
                pb.Active = false
                StopEventsBelow() 
            case pointer.Enter:
                pb.verticesHead.Hovered = true 
            case pointer.Leave:
                pb.verticesHead.Hovered = false
            }
        }
    }
}

func (pb *PolygonBuilder) RegisterEvents(gtx *layout.Context) {
    if !pb.Active {
        return
    }

    if pb.vertexCount > 2 {
        pb.registerClosingVertexEvent(gtx)
    }
}

func (pb *PolygonBuilder) registerClosingVertexEvent(gtx *layout.Context) {
    v := pb.verticesHead

    defer v.GetHoverEllipse().Push(gtx.Ops).Pop()  
    pointer.InputOp{
        Tag: &v.Hovered,
        Types: pointer.Press | pointer.Enter | pointer.Leave,
    }.Add(gtx.Ops)
}

func (pb *PolygonBuilder) drawCursorTailLine(ops *op.Ops) {
    var path clip.Path 
    path.Begin(ops)
    path.MoveTo(pb.verticesTail.Point)
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

    if pb.verticesHead == nil {
        pb.verticesHead = vertex
    } else {
        pb.verticesTail.next = vertex 
        vertex.previous = pb.verticesTail
    }
    pb.vertexCount += 1
    pb.verticesTail = vertex
    pb.setTailEnd(pointerPosition)
}

func drawPathFromVertices(head *Vertex, ops *op.Ops, color *color.NRGBA) {
    path := getPathFromVertices(head, ops, *color)

    paint.FillShape(ops, *color,
    clip.Stroke{
        Path: path.End(),
        Width: float32(lineWidth),
    }.Op()) 
}

func getPathFromVertices(head *Vertex, ops *op.Ops, color color.NRGBA) clip.Path {
    var path clip.Path

    current := head

    for current != nil {
        current.Layout(ops, color)
        current = current.next
    }

    current = head

    path.Begin(ops)
    path.MoveTo(current.Point)
    current = current.next
    for current != nil {
        path.LineTo(current.Point)
        current = current.next
    }

    return path
}

package main

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
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

    drawPolygonFromVertices(pb.verticesHead, pb.vertexCount, gtx.Ops, &pb.Color)
    pb.drawCursorTailLine(gtx.Ops)
}

func (pb *PolygonBuilder) HandleEvents(x *pointer.Event) {
    if pb.vertexCount > 2 {
        pb.handleClosingVertexEvent(x) 
    }

    pb.handleAddVertexEvent(x)
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

func (pb *PolygonBuilder) handleClosingVertexEvent(x *pointer.Event) {
    if eventsStoppedInFrame ||
    x.Type != pointer.Press ||
    !pb.verticesHead.IsClicked(x.Position) {
        return
    }

    CreatePolygon(
        pb.verticesHead,
        pb.verticesTail,
        pb.vertexCount,
        polygonColor,
    ) 
    pb.verticesHead = nil
    pb.verticesTail = nil
    pb.vertexCount = 0
    pb.Active = false
    StopEventsBelow() 
}

func (pb *PolygonBuilder) drawCursorTailLine(ops *op.Ops) {
    applicationPainter.DrawLine(pb.verticesTail.Point, pb.PointerTailEndPosition, pb.Color, ops)
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

func drawPolygonFromVertices(head *Vertex, count int, ops *op.Ops, color *color.NRGBA) {
    current := head
    for i := 0; i < count; i++ {
        current.Layout(ops, *color)
        if current.next == nil {
            break
        }
        applicationPainter.DrawLine(current.Point, current.next.Point, *color, ops)
        middlePoint := current.Point.Add(current.next.Point).Div(2.0)
        halfWidth := float32(edgeHoverWidth)
        if current.EdgeConstraint == Vertical {
            applicationPainter.DrawLine(
                f32.Point{X: middlePoint.X - halfWidth, Y: middlePoint.Y},
                f32.Point{X: middlePoint.X + halfWidth, Y: middlePoint.Y},
                constraintColor,
                ops,
            )
        } else if current.EdgeConstraint == Horizontal {
            applicationPainter.DrawLine(
                f32.Point{X: middlePoint.X, Y: middlePoint.Y - halfWidth},
                f32.Point{X: middlePoint.X, Y: middlePoint.Y + halfWidth},
                constraintColor,
                ops,
            )
        }
        current = current.next
    }
}

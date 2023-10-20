package main

import (
	"image/color"

	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type Polygon struct {
    Vertices []*Vertex
    selectedVertexId int
    Color color.NRGBA
    Edges []*Edge

    // pointer.Drag event helpers
    dragID pointer.ID
    dragPosition f32.Point
    drag bool
}

func CreatePolygon(vertices []*Vertex, color color.NRGBA) {
    var edges []*Edge
    for i := 0; i < len(vertices); i += 1 {
        j := (i + 1) % len(vertices)
        edges = append(edges, &Edge{Vertices: [2]*Vertex{vertices[i], vertices[j]}})
    }

    polygons = append(
        polygons,
        &Polygon{
            Vertices: vertices,
            Color: color,
            Edges: edges,
            selectedVertexId: -1,
        },
    )
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
            Path: fullPath,
            Width: float32(lineWidth),
        }.Op(),
    ) 
}

func (p *Polygon) HandleEvents(gtx *layout.Context) {
    if eventsStoppedInFrame {
        return
    }

    for _, e := range gtx.Events(p) {
        if x, ok := e.(key.Event); ok {
            switch x.Name {
            case "D":
                    if len(p.Vertices) <= 3 {
                        break 
                    }

                    vId := p.selectedVertexId
                    if p.Vertices[vId].Selected {
                        p.Vertices = append(p.Vertices[:vId], p.Vertices[vId+1:]...)
                    }
            } 
        }
    }

    for i, vertex := range p.Vertices {
        for _, e := range gtx.Events(vertex) {
            if x, ok := e.(pointer.Event); ok {
                switch x.Type {
                case pointer.Drag:
                    vertex.Selected = false
                    p.selectedVertexId = -1
                    vertex.Point = x.Position
                    StopEventsBelow()
                case pointer.Press:

                    if vertex.Selected {
                        vertex.Selected = false                      
                        p.selectedVertexId = -1
                    } else {
                        if p.selectedVertexId >= 0 {
                            p.Vertices[p.selectedVertexId].Selected = false

                        }
                        vertex.Selected = true
                        p.selectedVertexId = i
                    }
                    StopEventsBelow()
                case pointer.Enter:
                    vertex.Hovered = true
                case pointer.Leave:
                    vertex.Hovered = false
                }
            } 
        }
    }
}

func (p *Polygon) HandleEdgeEvents(gtx *layout.Context) {
    for _, edge := range p.Edges {
        for _, e := range gtx.Events(&edge.EventTag) {
            if x, ok := e.(pointer.Event); ok {
                switch x.Type {
                case pointer.Drag:
                    if p.dragID != x.PointerID {
                        break
                    }

                    edge.MoveBy(
                        x.Position.X - p.dragPosition.X,
                        x.Position.Y - p.dragPosition.Y,
                    )
                    p.dragPosition = x.Position
                    StopEventsBelow()
                case pointer.Press:
                    StopEventsBelow() 
                    if p.drag {
                        break
                    }

                    p.dragID = x.PointerID
                    p.dragPosition = x.Position
                case pointer.Release:
                    fallthrough
                case pointer.Enter:
                    edge.EventTag = true
                case pointer.Leave:
                    edge.EventTag = false
                }
            } 
        }
    }
}

func (p *Polygon) RegisterEvents(gtx *layout.Context) {
    key.InputOp{
        Tag: p,
        Keys: "d|D",
    }.Add(gtx.Ops)

    var area clip.Stack
    for _, vertex := range p.Vertices {
        area = vertex.GetHoverEllipse().Push(gtx.Ops)
        pointer.InputOp{
            Tag: vertex,
            Types: pointer.Drag | pointer.Press | pointer.Release | pointer.Enter | pointer.Leave,
        }.Add(gtx.Ops)
        area.Pop()
    }

}

func (p *Polygon) RegisterEdgeEvents(gtx *layout.Context) {
    var area clip.Stack
    for _, edge := range p.Edges {
        area = edge.GetHoverClipArea(gtx.Ops).Op().Push(gtx.Ops)
        pointer.InputOp{
            Tag: &edge.EventTag,
            Types: pointer.Press | pointer.Drag | pointer.Release, 
            Grab: p.drag,
        }.Add(gtx.Ops)
        area.Pop()
    } 
}

package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
)

var brushColor = color.NRGBA{R: 255, A: 255}
var backgroundColor = color.NRGBA{A:255}
var lineWidth = 4
var eventsStoppedInFrame = false

var polygonBuilder *PolygonBuilder
var polygons []*Polygon
var globalEventTag bool

var selected Selectable
var selectedDragId pointer.ID
var selectedDragPosition f32.Point
var selectedEdge *PolygonEdge

func main() {
    go func() {
        w := app.NewWindow(
            app.Title("Hexitor"),
            app.Size(1920, 1080),
        )
        err := run(w)
        if err != nil {
            log.Fatal(err)
        }
        os.Exit(0)
    }()
    app.Main()
}

func run(w *app.Window) error {
    var ops op.Ops
    polygonBuilder = &PolygonBuilder{Color: brushColor}
    for {
        e := <-w.Events()
        switch e := e.(type) {
        case system.DestroyEvent:
            return e.Err
        case system.FrameEvent:
            gtx := layout.NewContext(&ops, e)

            handleFrameEvents(&gtx)

            e.Frame(gtx.Ops)
            eventsStoppedInFrame = false
        }
    }

}

func handleFrameEvents(gtx *layout.Context) {
    drawBackground(gtx.Ops)
    polygonBuilder.HandleEvents(gtx)
    handleGlobalEvents(gtx)

    polygonBuilder.Layout(gtx)
    registerGlobalEvents(gtx)
    drawPolygons(gtx)
    if selected != nil {
        selected.HighLight(gtx)
    }

    polygonBuilder.RegisterEvents(gtx)
}

func drawBackground(ops *op.Ops) {
    paint.ColorOp{Color: backgroundColor}.Add(ops)
    paint.PaintOp{}.Add(ops)
}

func drawPolygons(gtx *layout.Context) {
    for _, polygon := range polygons {
        polygon.Layout(gtx)
    }
}

func StopEventsBelow() {
    eventsStoppedInFrame = true 
}

func handleGlobalEvents(gtx *layout.Context) {
    for _, e := range gtx.Events(&globalEventTag) {
        if x, ok := e.(key.Event); ok {
            if x.State != key.Press {
                break
            }

            switch x.Name {
            case "A":
                if selectedEdge != nil {
                    e := selectedEdge.getEdge().GetMiddlePoint()
                    prev := selectedEdge.getEdge().Vertices[0]
                    polygon := selectedEdge.Polygon
                    polygon.AppendVertexAfter(prev, e)
                    polygon.CreateEdges()
                    selected = nil
                    selectedEdge = nil
                }
            case "C":
                if selectedEdge != nil {
                   selectedEdge.getEdge().SetConstraint(None) 
                }
            case "D":
                if selected != nil {
                    selected.Destroy()
                    selected = nil
                    StopEventsBelow()
                }
            case "N":
                polygonBuilder.Active = true
            case "H":
                if selectedEdge != nil {
                    selectedEdge.getEdge().SetConstraint(Horizontal)
                }
            case "V":
                if selectedEdge != nil {
                    selectedEdge.getEdge().SetConstraint(Vertical)
                }
            }
        }
        if x, ok := e.(pointer.Event); ok {
            // handle PolygonBuilder global Events
            polygonBuilder.handleAddVertexEvent(&x)

            if eventsStoppedInFrame {
                break
            }

            // handle Events of selected object
            switch x.Type {
            case pointer.Press:
                selected = nil
                selectedEdge = nil
                for _, polygon := range polygons {
                    // Handle vertex click
                    vertex := polygon.VerticesHead
                    for i := 0; vertex != nil; i += 1 {
                        if vertex.IsClicked(x.Position) {
                            selected = &PolygonVertex{Polygon: polygon, Vertex: vertex}
                            selectedDragPosition = x.Position
                            selectedDragId = x.PointerID
                            StopEventsBelow()
                            return
                        }
                        vertex = vertex.next
                    }

                    // Handle edge click
                    for i, edge := range polygon.Edges {
                        if edge.IsClicked(&x.Position) {
                            pe := &PolygonEdge{Polygon: polygon, EdgeIndex: i}
                            selected = pe
                            selectedEdge = pe
                            selectedDragPosition = x.Position
                            selectedDragId = x.PointerID
                            StopEventsBelow()
                            return
                        }
                    }

                    // Handle polygon click
                    if polygon.IsClicked(x.Position) {
                        selected = polygon
                        selectedDragPosition = x.Position
                        selectedDragId = x.PointerID
                        StopEventsBelow()
                        return
                    }
                }
            case pointer.Drag:
                if selected != nil && selectedDragId == x.PointerID {
                    dp := selectedDragPosition
                    pos := x.Position
                    selected.MoveBy(pos.X - dp.X, pos.Y - dp.Y, gtx) 
                    selectedDragPosition = pos
                    StopEventsBelow()
                }
            }
        }
    }
}

func registerGlobalEvents(gtx *layout.Context) {
    key.InputOp{
        Tag: &globalEventTag,
        Keys: "A|C|D|N|H|V",
    }.Add(gtx.Ops)

    pointer.InputOp{
        Tag: &globalEventTag,
        Types: pointer.Press | pointer.Release | pointer.Drag | pointer.Move,
    }.Add(gtx.Ops)
}


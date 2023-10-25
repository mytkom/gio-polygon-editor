package main

import (
	"gk1-project1/painter"
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
	"gioui.org/widget"
	"gioui.org/widget/material"
)

var polygonColor = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
var builderColor = color.NRGBA{R: 255, A: 255}
var constraintColor = color.NRGBA{G: 255, A: 255}
var offsetColor = color.NRGBA{R: 178, G: 255, B:255, A: 255}
var backgroundColor = color.NRGBA{A:255}
var applicationPainter painter.Painter = &painter.GioPainter{}
var eventsStoppedInFrame = false
var offsetPolygonFeatureEnabled = false
var polygonOffset = 15

var polygonBuilder *PolygonBuilder
var polygons []*Polygon
var globalEventTag bool

var hovered Selectable
var selected Selectable
var selectedDragId pointer.ID
var selectedDragPosition f32.Point
var selectedEdge *PolygonEdge

var painterRadioButton widget.Enum =  widget.Enum{Value: "Gio"}

func main() {
    go func() {
        w := app.NewWindow(
            app.Title("Hexitor"),
            app.Size(800, 600),
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
    th := material.NewTheme()
    polygonBuilder = &PolygonBuilder{Color: builderColor}
    loadDefaultScene()
    for {
        e := <-w.Events()
        switch e := e.(type) {
        case system.DestroyEvent:
            return e.Err
        case system.FrameEvent:
            gtx := layout.NewContext(&ops, e)

            handleFrameEvents(&gtx, th)

            e.Frame(gtx.Ops)
            eventsStoppedInFrame = false
        }
    }

}

func handleFrameEvents(gtx *layout.Context, th *material.Theme) {
    drawBackground(gtx.Ops)
    handleEvents(gtx)
    polygonBuilder.Layout(gtx)
    registerEvents(gtx)
    DrawControlPanel(gtx, th)
    drawPolygons(gtx)
    if hovered != nil {
        hovered.HighLight(gtx)
    }

    if selected != nil {
        selected.HighLight(gtx)
    }
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

func handleEvents(gtx *layout.Context) {
    if painterRadioButton.Changed() {
        if painterRadioButton.Value == "Gio" {
            applicationPainter = &painter.GioPainter{}
        } else {
            applicationPainter = &painter.BresenhamPainter{}
        }
    }

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
            case "P":
                if applicationPainter.Type() == painter.Gio {
                    applicationPainter = &painter.BresenhamPainter{}
                } else {
                    applicationPainter = &painter.GioPainter{}
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
            case "O":
                if offsetPolygonFeatureEnabled {
                    offsetPolygonFeatureEnabled = false
                } else {
                    offsetPolygonFeatureEnabled = true
                }
            case "+":
                if offsetPolygonFeatureEnabled {
                   polygonOffset++ 
                }
            case "-":
                if offsetPolygonFeatureEnabled && polygonOffset > 1 {
                    polygonOffset--
                }
            }   
        }
        if x, ok := e.(pointer.Event); ok {
            // handle PolygonBuilder global Events
            polygonBuilder.HandleEvents(&x)

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
                    for i := 0; i < polygon.VerticesCount; i++ {
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
            case pointer.Move:
                hovered = nil
                for _, polygon := range polygons {
                    // Handle vertex hover
                    vertex := polygon.VerticesHead
                    for i := 0; i < polygon.VerticesCount; i++ {
                        if vertex.IsClicked(x.Position) {
                            hovered = &PolygonVertex{Polygon: polygon, Vertex: vertex}
                            return
                        }
                        vertex = vertex.next
                    }

                    // Handle edge hover
                    for i, edge := range polygon.Edges {
                        if edge.IsClicked(&x.Position) {
                            pe := &PolygonEdge{Polygon: polygon, EdgeIndex: i}
                            hovered = pe
                            return
                        }
                    }

                    // Handle polygon hover
                    if polygon.IsClicked(x.Position) {
                        hovered = polygon
                        return
                    }
                }

            }
        }
    }
}

func registerEvents(gtx *layout.Context) {
    key.InputOp{
        Tag: &globalEventTag,
        Keys: "A|C|D|P|N|H|V|O|+|-",
    }.Add(gtx.Ops)

    pointer.InputOp{
        Tag: &globalEventTag,
        Types: pointer.Press | pointer.Release | pointer.Drag | pointer.Move,
    }.Add(gtx.Ops)
}

func DrawControlPanel(gtx *layout.Context, th *material.Theme) {
    layout.Flex{Axis: layout.Axis(Vertical), Spacing: layout.SpaceStart}.Layout(
        *gtx,
        layout.Rigid(
            func(gtx layout.Context) layout.Dimensions {
                return layout.Flex{Spacing: layout.SpaceStart}.Layout(gtx,
                layout.Rigid(
                    func(gtx layout.Context) layout.Dimensions {
                        btn := material.RadioButton(th, &painterRadioButton, "Gio", "Gio")
                        btn.Color = polygonColor
                        return btn.Layout(gtx)
                    },
                ),
                layout.Rigid(
                    func(gtx layout.Context) layout.Dimensions {
                        btn := material.RadioButton(th, &painterRadioButton, "Bresenham", "Bresenham")
                        btn.Color = polygonColor
                        return btn.Layout(gtx)
                    },
                ),
            )
            },
        ),
    )
}

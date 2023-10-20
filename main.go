package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

var brushColor = color.NRGBA{R: 255, A: 255}
var backgroundColor = color.NRGBA{A:255}
var polygonBuilder *PolygonBuilder
var lineWidth = 4
var polygons []*Polygon
var selectedPolygon *Polygon
var eventsStoppedInFrame = false

func main() {
    go func() {
        w := app.NewWindow(
            app.Title("Hexitor"),
            app.Size(unit.Dp(500), unit.Dp(500)),
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
    handleSelectPolygonEvent(gtx)
    for _, polygon := range polygons {
        polygon.HandleEvents(gtx)
    }
    polygonBuilder.HandleEvents(gtx)
    polygonBuilder.Layout(gtx)
    polygonBuilder.RegisterEvents(gtx)
    registerSelectPolygonEvent(gtx)
    drawPolygons(gtx)

    if selectedPolygon != nil {
        hoverSelectedPolygon(gtx)
    }
}

func drawBackground(ops *op.Ops) {
    paint.ColorOp{Color: backgroundColor}.Add(ops)
    paint.PaintOp{}.Add(ops)
}

func drawPolygons(gtx *layout.Context) {
    for _, polygon := range polygons {
        polygon.Layout(gtx)
        polygon.RegisterEvents(gtx)
    }
}

func StopEventsBelow() {
    eventsStoppedInFrame = true 
}

func handleSelectPolygonEvent(gtx *layout.Context) {
    for _, polygon := range polygons {
        for _, e := range gtx.Events(polygon) {
            if x, ok := e.(pointer.Event); ok {
                switch x.Type {
                case pointer.Press:
                    selectedPolygon = polygon
                    StopEventsBelow()
                }
            }
        }
    }
}

func registerSelectPolygonEvent(gtx *layout.Context) {
    var area clip.Stack
    for _, polygon := range polygons {
        path := getPathFromVertices(polygon.Vertices, gtx.Ops, color.NRGBA{})
        path.Close()
        area = clip.Outline{Path: path.End()}.Op().Push(gtx.Ops)
        pointer.InputOp{
            Tag: polygon,
            Types: pointer.Press,
        }.Add(gtx.Ops)
        area.Pop()
    }
}

func hoverSelectedPolygon(gtx *layout.Context) {
    drawPolygonFromVertices(selectedPolygon.Vertices, gtx.Ops, &color.NRGBA{R: 255, G: 252, B: 127})
}

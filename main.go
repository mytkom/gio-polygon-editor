package main

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

var windowInputTag = false
var polygonBuilder PolygonBuilder

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
    polygonBuilder = PolygonBuilder{active: false}
    for {
        e := <-w.Events()
        switch e := e.(type) {
        case system.DestroyEvent:
            return e.Err
        case system.FrameEvent:
            gtx := layout.NewContext(&ops, e)

            drawBackground(&ops)
            handleFrameEvents(&ops, e.Queue)

            e.Frame(gtx.Ops)
        }
    }
}

func drawBackground(ops *op.Ops) {
    backgroundColor := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
    paint.ColorOp{Color: backgroundColor}.Add(ops)
    paint.PaintOp{}.Add(ops)
}

func handleFrameEvents(ops *op.Ops, q event.Queue) {
    for _, ev := range q.Events(&windowInputTag) {
        if x, ok := ev.(pointer.Event); ok {
            switch x.Type {
            case pointer.Press:
                polygonBuilder.AddVertex(x.Position)
            case pointer.Move:
                polygonBuilder.SetTailEnd(x.Position)
            }
        }
    }

    polygonBuilder.Layout(ops)

    pointer.InputOp{
        Tag:   &windowInputTag,
        Types: pointer.Press | pointer.Release | pointer.Move,
    }.Add(ops)
}

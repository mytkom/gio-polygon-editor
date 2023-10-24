package painter

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

type PainterType int

const(
    Gio PainterType = iota
    Bresenham
)

type Painter interface {
    Type() PainterType
    DrawLine(a f32.Point, b f32.Point, color color.NRGBA, ops *op.Ops) 
}

type GioPainter struct {}

func (gp *GioPainter) DrawLine(a f32.Point, b f32.Point, color color.NRGBA, ops *op.Ops) {
    var path clip.Path
    path.Begin(ops)
    path.MoveTo(a)
    path.LineTo(b)
    paint.FillShape(ops, color,
        clip.Stroke{Path: path.End(), Width: float32(1)}.Op(),
    )
}

func (gp *GioPainter) Type() PainterType {
    return Gio
}

type BresenhamPainter struct {}

func (bp *BresenhamPainter) DrawLine(a f32.Point, b f32.Point, color color.NRGBA, ops *op.Ops) {
    s := a.Round()
    d := b.Round()
    dx := math.Abs(float64(d.X - s.X))
    dy := -math.Abs(float64(d.Y - s.Y))
    sx := 1.0
    if s.X > d.X {
        sx = -1.0
    }
    sy := 1.0
    if s.Y > d.Y {
        sy = -1.0
    }
    err := dx + dy

    x := int(s.X)
    y := int(s.Y)

    for {
        offset := op.Offset(image.Pt(x, y)).Push(ops)
        rect := clip.Rect{
            Min: image.Pt(0, 0),
            Max: image.Pt(1, 1),
        }.Push(ops)
        paint.ColorOp{Color: color}.Add(ops)
        paint.PaintOp{}.Add(ops)
        rect.Pop()
        offset.Pop()

        if x == int(d.X) && y == int(d.Y) {
            break
        }

        e2 := 2 * err
        if e2 >= dy {
            err += dy
            x += int(sx)
        }
        if e2 <= dx {
            err += dx
            y += int(sy)
        }
    }
}

func (bp *BresenhamPainter) Type() PainterType {
    return Bresenham
}

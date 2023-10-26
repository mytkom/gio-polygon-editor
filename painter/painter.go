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


type Painter interface {
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

type WUPainter struct {}

func (bp *WUPainter) DrawLine(a f32.Point, b f32.Point, color color.NRGBA, ops *op.Ops) {
    drawWuLine(ops, a.X, a.Y, b.X, b.Y, color)
}

// https://en.wikipedia.org/wiki/Xiaolin_Wu%27s_line_algorithm
func drawWuLine(ops *op.Ops, x0, y0, x1, y1 float32, color color.NRGBA) {
    dx := float32(math.Abs(float64(x1 - x0)))
    dy := float32(math.Abs(float64(y1 - y0)))

    steep := dy > dx
    if steep {
        x0, y0 = y0, x0
        x1, y1 = y1, x1
    }
    if x0 > x1 {
        x0, x1 = x1, x0
        y0, y1 = y1, y0
    }

    dx = float32(math.Abs(float64(x1 - x0)))
    dy = float32(math.Abs(float64(y1 - y0)))

    gradient := dy / dx
    if dx == 0 {
        gradient = 1.0
    }

    if y0 > y1 {
        y0, y1 = y1, y0
    }

    xend := round(x0)
    yend := y0 - gradient*(float32(xend)-x0)
    xgap := rfpart(x0 + 0.5)
    xpxl1 := xend
    ypxl1 := int(yend)

    if steep {
        plot(ops, ypxl1, xpxl1, rfpart(yend)*xgap, color)
        plot(ops, ypxl1+1, xpxl1, fpart(yend)*xgap, color)
    } else {
        plot(ops, xpxl1, ypxl1, rfpart(yend)*xgap, color)
        plot(ops, xpxl1, ypxl1+1, fpart(yend)*xgap, color)
    }

    intery := yend + gradient

    xend = round(x1)
    yend = y1 + gradient*(float32(xend)-x1)
    xgap = fpart(x1 + 0.5)
    xpxl2 := xend
    ypxl2 := int(yend)

    if steep {
        plot(ops, ypxl2, xpxl2, rfpart(yend)*xgap, color)
        plot(ops, ypxl2+1, xpxl2, fpart(yend)*xgap, color)
    } else {
        plot(ops, xpxl2, ypxl2, rfpart(yend)*xgap, color)
        plot(ops, xpxl2, ypxl2+1, fpart(yend)*xgap, color)
    }

    if steep {
        for x := xpxl1 + 1; x < xpxl2; x++ {
            plot(ops, int(intery), x, rfpart(intery), color)
            plot(ops, int(intery)+1, x, fpart(intery), color)
            intery += gradient
        }
    } else {
        for x := xpxl1 + 1; x < xpxl2; x++ {
            plot(ops, x, int(intery), rfpart(intery), color)
            plot(ops, x, int(intery)+1, fpart(intery), color)
            intery += gradient
        }
    }
}

func plot(ops *op.Ops, x, y int, c float32, col color.NRGBA) {
    cr := col.R
    cg := col.G
    cb := col.B
    ca := col.A

    alpha := float32(ca) * c
    if alpha > 0 {
        offset := op.Offset(image.Pt(x, y)).Push(ops)
        rect := clip.Rect{
            Min: image.Pt(0, 0),
            Max: image.Pt(1, 1),
        }.Push(ops)
        paint.ColorOp{Color: color.NRGBA{
            R: cr,
            G: cg,
            B: cb,
            A: uint8(alpha),
        }}.Add(ops)
        paint.PaintOp{}.Add(ops)
        rect.Pop()
        offset.Pop()
    }
}

func round(x float32) int {
    return int(x + 0.5)
}

func fpart(x float32) float32 {
    return x - float32(math.Floor(float64(x)))
}

func rfpart(x float32) float32 {
    return 1 - fpart(x)
}






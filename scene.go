package main

import (
	"gioui.org/f32"
)

func loadDefaultScene() {
    polygonA := []*Vertex{
        {Point: f32.Point{X: 100, Y: 100}, EdgeConstraint: Horizontal},
        {Point: f32.Point{X: 200, Y: 100}, EdgeConstraint: Vertical},
        {Point: f32.Point{X: 200, Y: 200}, EdgeConstraint: Horizontal},
        {Point: f32.Point{X: 300, Y: 200}},
        {Point: f32.Point{X: 350, Y: 250}, EdgeConstraint: Horizontal},
        {Point: f32.Point{X: 100, Y: 250}, EdgeConstraint: Vertical},
    }
    polygonB := []*Vertex{
        {Point: f32.Point{X: 400, Y: 400}},
        {Point: f32.Point{X: 420, Y: 420}},
        {Point: f32.Point{X: 420, Y: 460}},
        {Point: f32.Point{X: 380, Y: 500}, EdgeConstraint: Vertical},
        {Point: f32.Point{X: 380, Y: 450}},
        {Point: f32.Point{X: 360, Y: 400}},
    }

    loadPolygon(polygonA)
    loadPolygon(polygonB)
}

func loadPolygon(polygon []*Vertex) {
    head := polygon[0]
    tail := polygon[len(polygon) - 1]
    for i := 0; i < len(polygon) - 1; i++ {
        polygon[i].next = polygon[i + 1] 
        polygon[i + 1].previous = polygon[i]
    }

    CreatePolygon(head, tail, len(polygon), polygonColor)
}

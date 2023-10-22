package main

import "gioui.org/layout"

func isPointInWindow(x float32, y float32, gtx *layout.Context) bool {
    boundWidth := float32(10.0)
    max := gtx.Constraints.Max
    if x > float32(max.X) - boundWidth || x < boundWidth ||
       y > float32(max.Y) - boundWidth || y < boundWidth {
        return false 
    }

    return true
}



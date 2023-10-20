package main

import "image/color"

var hoverAlphaValue = uint8(32)

func HoverizeColor(c color.NRGBA) color.NRGBA {
    return color.NRGBA{
        R: c.R,
        G: c.G,
        B: c.B,
        A: hoverAlphaValue,
    }
}

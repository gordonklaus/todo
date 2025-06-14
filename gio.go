package main

import (
	"image/color"

	"gioui.org/layout"
)

type C = layout.Context
type D = layout.Dimensions

func mulAlpha(c color.NRGBA, alpha uint8) color.NRGBA {
	c.A = uint8(uint32(c.A) * uint32(alpha) / 0xFF)
	return c
}

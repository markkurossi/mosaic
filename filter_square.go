//
// Copyright (c) 2021-2022 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"image"
	"image/color"
)

type FilterSquare struct {
	Width int
}

func (f *FilterSquare) Init(hist Histogram) {
}

func (f *FilterSquare) Tiled() bool {
	return true
}

func (f *FilterSquare) Transform(input image.Image, output *image.NRGBA) {
	bounds := input.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var c color.NRGBA
			var border bool
			if y-bounds.Min.Y <= f.Width || bounds.Max.Y-y <= f.Width ||
				x-bounds.Min.X <= f.Width || bounds.Max.X-x <= f.Width {
				border = true
			}
			if border {
				r, g, b, a := input.At(x, y).RGBA()
				c = color.NRGBA{
					R: uint8(r >> 8),
					G: uint8(g >> 8),
					B: uint8(b >> 8),
					A: uint8(a >> 8),
				}
			} else {
				c.A = 255
			}
			output.Set(x, y, c)
		}
	}
}

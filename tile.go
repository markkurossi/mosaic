//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"image"
	"image/color"
)

type Tile struct {
	base   image.Image
	bounds image.Rectangle
}

func (t *Tile) ColorModel() color.Model {
	return t.base.ColorModel()
}

func (t *Tile) Bounds() image.Rectangle {
	return t.bounds
}

func (t *Tile) At(x, y int) color.Color {
	return t.base.At(x, y)
}

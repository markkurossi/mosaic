//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"image"
)

type Histogram []uint64

func NewHistogram(img image.Image) Histogram {
	bounds := img.Bounds()

	hist := make([]uint64, 65536)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			hist[Index16(img.At(x, y).RGBA())]++
		}
	}
	return hist
}

func Index16(r, g, b, a uint32) (idx uint16) {
	idx = uint16(r >> 12)
	idx <<= 4
	idx |= uint16(g >> 12)
	idx <<= 4
	idx |= uint16(b >> 12)
	idx <<= 4
	idx |= uint16(a >> 12)
	return
}

func Index12(r, g, b, a uint32) (idx uint16) {
	idx = uint16(r >> 13)
	idx <<= 3
	idx |= uint16(g >> 13)
	idx <<= 3
	idx |= uint16(b >> 13)
	idx <<= 3
	idx |= uint16(a >> 13)
	return
}

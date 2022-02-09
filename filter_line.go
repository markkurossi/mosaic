//
// Copyright (c) 2022 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"image"
	"image/color"
)

type FilterLine struct {
	Width int
}

func (f *FilterLine) Init(hist Histogram) {
}

func (f *FilterLine) Tiled() bool {
	return false
}

func (f *FilterLine) Transform(input image.Image, output *image.NRGBA) {
	bounds := input.Bounds()

	w := f.Width / 2
	h := f.Width

	numValues := (bounds.Max.X-bounds.Min.X)/w + 2

	values := make([]float64, numValues)

	for y := bounds.Min.Y; y+h <= bounds.Max.Y; y += h {
		i := 1
		for x := bounds.Min.X; x+w < bounds.Max.X; x += w {
			values[i] = countRegion(input, image.Rectangle{
				Min: image.Point{
					X: x,
					Y: y,
				},
				Max: image.Point{
					X: x + w,
					Y: y + h,
				},
			})
			i++
		}

		i = 1
		for x := bounds.Min.X; x+w < bounds.Max.X; x += w {
			v := uint8(values[i])
			paintRegion(output, image.Rectangle{
				Min: image.Point{
					X: x,
					Y: y,
				},
				Max: image.Point{
					X: x + w,
					Y: y + h,
				},
			}, color.NRGBA{
				R: v,
				G: v,
				B: v,
				A: uint8(0xff),
			})
			i++
		}
	}
}

func countRegion(input image.Image, r image.Rectangle) (result float64) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			r, g, b, _ := input.At(x, y).RGBA()

			result += float64(r >> 8)
			result += float64(g >> 8)
			result += float64(b >> 8)
		}
	}
	return result / float64((r.Max.Y-r.Min.Y)*(r.Max.X-r.Min.X)*3)
}

func paintRegion(output *image.NRGBA, r image.Rectangle, c color.NRGBA) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			output.Set(x, y, c)
		}
	}
}

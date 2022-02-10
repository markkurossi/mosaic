//
// Copyright (c) 2022 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"image"
	"image/color"
	"math"

	"github.com/fogleman/gg"
)

type FilterLine struct {
	Width  int
	Height int
}

func (f *FilterLine) Init(hist Histogram) {
}

func (f *FilterLine) Tiled() bool {
	return false
}

func (f *FilterLine) Transform(input image.Image, output *image.NRGBA) {
	bounds := input.Bounds()

	paintRegion(output, bounds, color.NRGBA{
		R: 0xff,
		G: 0xff,
		B: 0xff,
		A: 0xff,
	})

	gc := gg.NewContextForImage(output)
	gc.SetColor(color.NRGBA{
		A: 0xff,
	})
	gc.SetLineWidth(4)

	w := f.Width / 2
	h := f.Height

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
			from := values[i-1]
			to := values[i]
			f.line(gc, from, to, x, y, i == 1)
			i++
		}
		gc.Stroke()
	}

	gc.SavePNG(",out.png")
}

func (f *FilterLine) line(gc *gg.Context, from, to float64, x, y int,
	first bool) {

	steps := f.Width / 2

	for i := 0; i < steps; i++ {
		deg := float64((x+i)%f.Width) / float64(f.Width) * math.Pi * 2
		val := float64(steps-i)/float64(steps)*from +
			float64(i)/float64(steps)*to
		val = 255 - val
		d := int(math.Sin(deg) * val / 255.0 * float64(f.Height/2))

		if first {
			gc.MoveTo(float64(x+i), float64(y+f.Width/2+d))
		} else {
			gc.LineTo(float64(x+i), float64(y+f.Width/2+d))
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

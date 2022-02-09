//
// Copyright (c) 2021-2022 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

var palettes = [][]color.NRGBA{
	{
		{0xCB, 0x99, 0x7E, 0xFF},
		{0xDD, 0xBE, 0xA9, 0xFF},
		{0xFF, 0xE8, 0xD6, 0xFF},
		{0xB7, 0xB7, 0xA4, 0xFF},
		{0xA5, 0xA5, 0x8D, 0xFF},
		{0x6B, 0x70, 0x5C, 0xFF},
	},
	{
		{255, 205, 178, 255},
		{255, 180, 162, 255},
		{229, 152, 155, 255},
		{181, 131, 141, 255},
		{109, 104, 117, 255},
	},
	{
		{239, 71, 111, 255},
		{255, 209, 102, 255},
		{6, 214, 160, 255},
		{17, 138, 178, 255},
		{7, 59, 76, 255},
	},
	{
		{254, 197, 187, 255},
		{252, 213, 206, 255},
		{250, 225, 221, 255},
		{248, 237, 235, 255},
		{232, 232, 228, 255},
		{216, 226, 220, 255},
		{236, 228, 219, 255},
		{255, 229, 217, 255},
		{255, 215, 186, 255},
		{254, 200, 154, 255},
	},
	{
		{0x00, 0x30, 0x50, 0xFF},
		{0x70, 0x96, 0xA0, 0xFF},
		{0xB0, 0xB7, 0xA7, 0xFF},
		{0xFA, 0xE3, 0xAD, 0xFF},
		{0xDA, 0x14, 0x15, 0xFF},
	},
}

type PaletteType int

const (
	TypeClosest PaletteType = iota
	TypeDistribution
)

type FilterPalette struct {
	Palette []color.NRGBA
	Type    PaletteType

	index     []uint16
	bestIndex float64
}

func (f *FilterPalette) Init(hist Histogram) {
	f.index = make([]uint16, len(f.Palette))
	f.bestIndex = math.MaxFloat64

	var total uint64
	for _, count := range hist {
		total += uint64(count)
	}

	perBucket := total / uint64(len(f.Palette))
	fmt.Printf("total: %v, len(palette): %v, perBucket: %v\n",
		total, len(f.Palette), perBucket)

	f.minimize(hist, 0, 0, make([]uint16, len(f.Palette)),
		make([]uint64, len(f.Palette)), perBucket)

	for idx, count := range f.index {
		fmt.Printf("%d:\t%d\t%v\n", idx, count, perBucket)
	}
	fmt.Printf("badness: %v\n", f.bestIndex)
}

func (f *FilterPalette) minimize(hist Histogram, idx, pos int, index []uint16,
	counts []uint64, perBucket uint64) {

	if pos >= len(counts) {
		for i, idx := range index {
			fmt.Printf("%d\t%d\t%d\n", i, idx, counts[i])
		}
		copy(f.index, index)
		return
	}

	limit := len(hist) - (len(index) - pos)

	for i := idx; i < limit; i++ {
		if hist[i] == 0 {
			continue
		}

		index[pos] = uint16(i)
		counts[pos] += hist[i]

		if counts[pos] >= perBucket {
			fmt.Printf("%d: case 0: %v\n", pos, counts[pos])
			f.minimize(hist, i+1, pos+1, index, counts, perBucket)
			return
		} else if i+1 < limit && counts[pos]+hist[i+1] >= perBucket {
			if perBucket-counts[pos] <= counts[pos]+hist[i+1]-perBucket {
				fmt.Printf("%d: case 1: %v\n", pos, counts[pos])
				f.minimize(hist, i+1, pos+1, index, counts, perBucket)
				return
			} else {
				index[pos] = uint16(i + 1)
				counts[pos] += hist[i+1]
				fmt.Printf("%d: case 2: %v\n", pos, counts[pos])
				f.minimize(hist, i+2, pos+1, index, counts, perBucket)
				return
			}
		}
	}

	index[pos] = uint16(limit)
	fmt.Printf("%d: case 3: %v\n", pos, counts[pos])
	f.minimize(hist, limit, pos+1, index, counts, perBucket)
}

func (f *FilterPalette) Tiled() bool {
	return true
}

func (f *FilterPalette) Transform(input image.Image, output *image.NRGBA) {
	switch f.Type {
	case TypeClosest:
		f.fClosest(input, output)

	case TypeDistribution:
		f.fDistribution(input, output)
	}
}

func (f *FilterPalette) fClosest(input image.Image, output *image.NRGBA) {
	var rs, gs, bs, as uint32

	bounds := input.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := input.At(x, y).RGBA()
			rs += r >> 8
			gs += g >> 8
			bs += b >> 8
			as += a >> 8
		}
	}
	count := uint32((bounds.Max.X - bounds.Min.X) *
		(bounds.Max.Y - bounds.Min.Y))

	c := color.NRGBA{
		R: uint8(rs / count),
		G: uint8(gs / count),
		B: uint8(bs / count),
		A: uint8(as / count),
	}

	var bestDelta int
	var bestColor color.NRGBA

	for idx, color := range f.Palette {
		d := (diff(color.R, c.R) + diff(color.G, c.G) + diff(color.B, c.B)) / 3
		if idx == 0 || d < bestDelta {
			bestDelta = d
			bestColor = color
		}
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			output.Set(x, y, bestColor)
		}
	}
}

func (f *FilterPalette) fDistribution(input image.Image, output *image.NRGBA) {
	var rs, gs, bs, as uint32

	bounds := input.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := input.At(x, y).RGBA()
			rs += r >> 8
			gs += g >> 8
			bs += b >> 8
			as += a >> 8
		}
	}
	count := uint32((bounds.Max.X - bounds.Min.X) *
		(bounds.Max.Y - bounds.Min.Y))

	c := color.NRGBA{
		R: uint8(rs / count),
		G: uint8(gs / count),
		B: uint8(bs / count),
		A: uint8(as / count),
	}

	idx := uint16(c.R >> 2)
	idx <<= 4
	idx |= uint16(c.G >> 2)
	idx <<= 4
	idx |= uint16(c.B >> 2)
	idx <<= 4
	idx |= uint16(c.A >> 2)

	var color color.NRGBA

	for i := 0; i < len(f.index); i++ {
		if i+1 >= len(f.index) || f.index[i+1] > idx {
			color = f.Palette[i]
			break
		}
	}
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			output.Set(x, y, color)
		}
	}

}

func diff(a, b uint8) int {
	if a > b {
		return int(a - b)
	}
	return int(b - a)
}

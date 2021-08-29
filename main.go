//
// Copyright (c) 2021 Markku Rossi
//
// All rights reserved.
//

package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"

	_ "image/jpeg"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	for _, arg := range flag.Args() {
		err := processFile(arg)
		if err != nil {
			log.Fatalf("failed to process file '%s': %s\n", arg, err)
		}
	}
}

func printPalette(idx int, palette []color.NRGBA) {
	fmt.Printf("%d\t", idx)
	fmt.Printf("\033[38;2;255;82;197mHello")
	for _, c := range palette {
		fmt.Printf("\033[48;2;%d;%d;%dm*\033[m", c.R, c.G, c.B)
	}
	fmt.Println()
}

func processFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	m, _, err := image.Decode(f)
	if err != nil {
		return err
	}
	bounds := m.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y

	fmt.Printf("%s: %dx%d\n", path, width, height)

	output := image.NewNRGBA(image.Rectangle{
		Max: image.Point{
			X: width,
			Y: height,
		},
	})

	tile := 1

	for y := bounds.Min.Y; y < bounds.Max.Y; y += tile {
		for x := bounds.Min.X; x < bounds.Max.X; x += tile {

			maxX := x + tile
			if maxX > bounds.Max.X {
				maxX = bounds.Max.X
			}
			maxY := y + tile
			if maxY > bounds.Max.Y {
				maxY = bounds.Max.Y
			}

			filter(&Tile{
				base: m,
				bounds: image.Rectangle{
					Min: image.Point{
						X: x,
						Y: y,
					},
					Max: image.Point{
						X: maxX,
						Y: maxY,
					},
				},
			}, output)
		}
	}

	name := fmt.Sprintf("%s.mosaic.png", path)
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()

	return png.Encode(out, output)
}

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

func filter(input image.Image, output *image.NRGBA) {
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

	for idx, color := range palettes[4] {
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

func diff(a, b uint8) int {
	if a > b {
		return int(a - b)
	}
	return int(b - a)
}

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

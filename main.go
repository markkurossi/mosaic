//
// Copyright (c) 2021-2022 Markku Rossi
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

type Filter interface {
	Init(hist Histogram)
	Tiled() bool
	Transform(input image.Image, output *image.NRGBA)
}

var (
	_ Filter = &FilterPalette{}
	_ Filter = &FilterLine{}
	_ Filter = &FilterSquare{}
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	var f Filter

	switch 2 {
	case 0:
		f = &FilterPalette{
			Palette: palettes[4],
			//Type:    TypeDistribution,
		}

	case 1:
		f = &FilterSquare{}

	case 2:
		f = &FilterLine{
			Width:  30,
			Height: 30,
		}
	}

	for _, arg := range flag.Args() {
		err := processFile(arg, f)
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

func processFile(path string, filter Filter) error {
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

	hist := NewHistogram(m)
	var nonZero int
	for _, count := range hist {
		if count > 0 {
			nonZero++
			// fmt.Printf("%d:\t%d\n", idx, count)
		}
	}
	fmt.Printf("%s: #nonZero: %d\n", path, nonZero)
	filter.Init(hist)

	output := image.NewNRGBA(image.Rectangle{
		Max: image.Point{
			X: width,
			Y: height,
		},
	})

	if filter.Tiled() {
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

				filter.Transform(&Tile{
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
	} else {
		filter.Transform(&Tile{
			base:   m,
			bounds: bounds,
		}, output)
	}

	name := fmt.Sprintf("%s.mosaic.png", path)
	out, err := os.Create(name)
	if err != nil {
		return err
	}
	defer out.Close()

	return png.Encode(out, output)
}

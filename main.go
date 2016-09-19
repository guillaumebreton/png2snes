package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%v\n", os.Args)
		fmt.Println("usage : png2snes filepath.png")
		fmt.Println("output two file in the current directory : filepath.pic and filepath.clr")
		os.Exit(1)
	}
	// Load file
	infile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Fail to open file %s\n", os.Args[1])
		os.Exit(2)
	}
	defer infile.Close()
	src, _, err := image.Decode(infile)
	if err != nil {
		fmt.Printf("Fail to read file %s : %s\n", os.Args[1], err)
		os.Exit(3)
	}

	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	if w%8 != 0 || h%8 != 0 {
		fmt.Printf("Invalid image size  (%d x %d) the image size should be a multiple of 8", w, h)
		os.Exit(4)
	}
	palette := NewPalette()
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			bits := to15Bits(src.At(x, y))
			idx := palette.Add(bits)
			print(idx, " ")
		}
		println("")
	}
}

// A Palette store the distinct colors in a 15 bits LSB RGB values
type Palette struct {
	Colors []uint16
}

func NewPalette() *Palette {
	return &Palette{make([]uint16, 0)}
}
func (p *Palette) Add(c uint16) int {
	index := p.Index(c)
	if index == -1 {
		p.Colors = append(p.Colors, c)
		return len(p.Colors) - 1
	}
	return index
}

func (p *Palette) Index(bits uint16) int {
	for k, v := range p.Colors {
		if v == bits {
			return k
		}
	}
	return -1
}

func to15Bits(c color.Color) uint16 {
	rgba := c.(color.RGBA)
	var r uint16 = (uint16(rgba.R) & 0xF8) >> 19
	var g uint16 = (uint16(rgba.G) & 0xF8) >> 6
	var b uint16 = (uint16(rgba.B) & 0xF8) << 7
	return b | g | r
}

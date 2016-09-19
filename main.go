package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
)

var bitplanesNumer int = 4

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
	//get palette and generate the color map
	colors := make([][]int, h)
	palette := NewPalette()
	for y := 0; y < h; y++ {
		colors[y] = make([]int, w)
		for x := 0; x < w; x++ {
			bits := To15Bits(src.At(x, y))
			idx := palette.Add(bits)
			colors[y][x] = idx
		}
	}
	PrintColors(colors)
	println("---------- Bit plane -------------")
	//get all tiles data
	for x := 0; x < w; x += 8 {
		for y := 0; y < h; y += 8 {
			bp := GetTileBitplanes(colors, x, y, bitplanesNumer)
			PrintBitplane(bp)
		}
	}
	println("---------- END Bit plane -------------")
}

// GetTile data
// Bpn bitplane number
func GetTileBitplanes(colors [][]int, x, y, bpn int) []byte {
	bitplanes := make([]byte, 0)
	for b := 0; b < bpn; b += 2 {
		// for each color, extract the b and b+1 bit plane layer
		mask1 := 0x1 << uint(b)
		mask2 := 0x1 << (uint(b) + 1)
		for j := y; j < y+8; j++ {
			var bitplane1 byte = 0
			var bitplane2 byte = 0
			for i := x; i < x+8; i++ {
				cIdx := colors[j][i]
				shift := uint(i - x)
				var bit1 int = (cIdx & mask1) << (7 - shift - uint(b))
				var bit2 int = (cIdx & mask2) << (7 - shift - uint(b) - 1)
				bitplane1 = bitplane1 | byte(bit1)
				bitplane2 = bitplane2 | byte(bit2)
			}
			bitplanes = append(bitplanes, bitplane1)
			bitplanes = append(bitplanes, bitplane2)
		}
	}
	return bitplanes
}

func PrintBitplane(bp []byte) {
	for k, v := range bp {
		fmt.Printf("%02d - %08b\n", k+1, v)
	}
}

func PrintColors(colors [][]int) {
	for y := 0; y < len(colors); y++ {
		row := colors[y]
		for x := 0; x < len(row); x++ {
			print(colors[y][x], " ")
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

// To15Bits return a MSB 15bits color base on a RGBA color
func To15Bits(c color.Color) uint16 {
	MASK := 0xF8 // to keep the first n bytes
	rgba := c.(color.RGBA)
	var r int = (int(rgba.R) & MASK) >> 3
	var g int = (int(rgba.G) & MASK) << 2
	var b int = (int(rgba.B) & MASK) << 7
	return uint16(b | g | r)
}

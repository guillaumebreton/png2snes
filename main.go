package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"strings"
)

var in = flag.String("in", "", "the input tile set")
var clr = flag.String("out", ".", "The output director")

func main() {
	flag.Parse()
	if *in == "" {
		fmt.Println("Input file cannot be empty")
		os.Exit(1)
	}
	m, err := LoadMap(*in)
	if err != nil {
		fmt.Printf("Failed to load map : %s\n", err.Error())
		os.Exit(2)
	}
	m.Print()

}

func test(in, clr, pic *string) {

	if *in == "" {
		fmt.Println("Input file cannot be empty")
		os.Exit(1)
	}

	idx := strings.LastIndex(*in, ".")
	inPath := *in
	if idx != -1 {
		file := *in
		inPath = file[:idx]
	}

	clrFilePath := *clr
	if clrFilePath == "" {
		clrFilePath = inPath + ".clr"
	}

	picFilePath := *pic
	if picFilePath == "" {
		picFilePath = inPath + ".pic"
	}
	// Load file
	infile, err := os.Open(*in)
	if err != nil {
		fmt.Printf("Fail to open file %s\n", *in)
		os.Exit(2)
	}
	defer infile.Close()
	src, _, err := image.Decode(infile)
	if err != nil {
		fmt.Printf("Fail to read file %s : %s\n", *in, err)
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
	// write the CLR file
	err = palette.Write(clrFilePath)
	if err != nil {
		fmt.Printf("Fail to write %s file : %v\n", clrFilePath, err)
		os.Exit(5)
	}

	//get all tiles data
	bps := NewBitPlanes(4)
	for x := 0; x < w; x += 8 {
		for y := 0; y < h; y += 8 {
			bps.Add(colors, x, y)
		}
	}

	// write the PIC file
	err = bps.Write(picFilePath)
	if err != nil {
		fmt.Printf("Fail to write %s.clr file : %v\n", picFilePath, err)
		os.Exit(5)
	}
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

func GetWriter(filepath string) (*bufio.Writer, *os.File, error) {
	file, err := os.OpenFile(filepath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, nil, err
	}
	return bufio.NewWriter(file), file, err
}

type Bitplanes struct {
	bitplanes       []byte
	bitplanesNumber int
}

func NewBitPlanes(bitplanesNumber int) *Bitplanes {
	return &Bitplanes{make([]byte, 0), bitplanesNumber}
}

func (bps *Bitplanes) Add(colors [][]int, x, y int) {
	for b := 0; b < bps.bitplanesNumber; b += 2 {
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
			bps.bitplanes = append(bps.bitplanes, bitplane1)
			bps.bitplanes = append(bps.bitplanes, bitplane2)
		}
	}
}

func (b *Bitplanes) Write(filepath string) error {
	w, file, err := GetWriter(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	err = binary.Write(w, binary.LittleEndian, b.bitplanes)
	w.Flush()
	return err
}

// A Palette store the distinct colors in a 15 bits LSB RGB values
type Palette struct {
	Colors []uint16
}

func NewPalette() *Palette {
	return &Palette{make([]uint16, 0)}
}

func (p *Palette) Write(filepath string) error {
	w, file, err := GetWriter(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, v := range p.Colors {
		err := binary.Write(w, binary.LittleEndian, v)
		if err != nil {
			return err
		}
	}
	w.Flush()
	return nil
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

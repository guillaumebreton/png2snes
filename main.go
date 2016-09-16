package main

import (
	"fmt"
	"image"
	// "image/color"
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
	// Get the filename
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
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			color := src.At(x, y)
			r, g, b, _ := color.RGBA()
			println(r, g, b)
		}
	}
}

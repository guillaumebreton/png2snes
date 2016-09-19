package main

import (
	"fmt"
	"image/color"
	"testing"
)

func TestTo15BitsWhite(t *testing.T) {
	c := color.RGBA{0xFF, 0xFF, 0xFF, 255}
	if To15Bits(c) != 32767 {
		fmt.Printf("Invalid value : %d\n", To15Bits(c))
		t.Fail()
	}
}

//1 * 1024+2*32+2
func TestTo15BitsMisc(t *testing.T) {
	c := color.RGBA{0x11, 0x10, 0x0E, 255}
	if To15Bits(c) != 1090 {
		fmt.Printf("Invalid value : %d\n", To15Bits(c))
		t.Fail()
	}
}

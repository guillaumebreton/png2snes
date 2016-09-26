package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Map struct {
	Version     string    `xml:"version,attr"`
	Orientation string    `xml:"orientation,attr"`
	Renderorder string    `xml:"renderorder,attr"`
	Width       int       `xml:"width,attr"`
	Height      int       `xml:"height,attr"`
	Tilesets    []Tileset `xml:"tileset"`
	Layers      []Layer   `xml:"layer"`
	// properties, tileset, layer, objectgroup, imagelayer
}

func (m *Map) Print() {
	fmt.Printf("+ Map (%d x %d) version=%s renderorder=%s orientation=%s\n", m.Width, m.Height, m.Version, m.Orientation, m.Renderorder)
	// Generates palettes and bitplanes for each tileset
	fmt.Println("+ Tilesets")
	for k, v := range m.Tilesets {
		fmt.Printf("  - %d - name=%s", k, v.Name)
		if v.Image.Source != "" {
			fmt.Printf(" source=%s", v.Image.Source)
		}
		fmt.Println("")
	}
	// Generates layer for each tileset
	fmt.Println("+ Layers")
	for k, v := range m.Layers {
		// TODO Check that the layer does not contains more than one tile map
		fmt.Printf("  - %d - %s\n", k, v.Name)
	}
}

type Tileset struct {
	Name      string `xml:"name,attr"`
	Firstgid  int    `xml:"firstgid,attr"`
	Tilecount int    `xml:tilecount,attr`
	Image     Image  `xml:"image"`

	// Other props are useless for us
}

type Image struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

type Layer struct {
	Name string `xml:"name,attr"`
	// Other props are usless for us
}

func NewMap() *Map {
	return &Map{}
}

func LoadMap(filepath string) (*Map, error) {
	r, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	m := NewMap()
	err = xml.Unmarshal(b, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

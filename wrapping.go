// Wrapping demonstrates
// a repeating texture.
/*
NewWrapping is most useful with "seamless" (aka "tileable")
textures. Callers must be aware non-power-of-two images
are scaled up to the next power of two, so the best practice
is to use power-of-two textures.
*/
package main

import (
	"image"
	_ "image/png"
	"log"

	"github.com/j7b/two/geo"

	"github.com/j7b/two"
	"github.com/j7b/two/doc/_examples/assets"
	"github.com/j7b/two/glc"
)

const (
	width  = 640
	height = 480
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type sprite interface {
	two.Sprite
	geo.Source // to manipulate the sample area with floats
}

type base struct {
	width, height float64
	tex           two.Texture
	s             sprite
}

func (b *base) init() {
	glc.ClearColor(0, 0, 255, 255)
	f, err := assets.Open("trees.png")
	check(err)
	img, _, err := image.Decode(f)
	check(err)
	b.tex, err = two.NewWrapping(img)
	check(err)
	b.s = b.tex.NewSprite(int(b.width), int(b.height)).(sprite)
	b.s.SetLoc(b.width/2, b.height/2)
}

func (b *base) Update(f float64) {
	if b.tex == nil {
		b.init()
		return
	}
	b.s.SetSource(b.s.Source().Translate(4*f, f))
}

func main() {
	log.Fatal(two.Run(&base{width: width, height: height}, nil, width, height))
}

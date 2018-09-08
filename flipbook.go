// Flipbook demonstrates
// texture image replacement.
/*
Many platforms can replace all the textures
the platform supports at every update (at least as
implemented here; DMA transfers should complete in
well less than 16ms). Others are less capable, particularly
the WebGL environment and some SOC platforms.
Images of type *image.RGBA have a relatively fast path,
(textures are stored premultipled) but other types have conversion overhead.
*/
package main

import (
	"image"
	"image/draw"
	"image/gif"
	"log"

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

type bananas []two.Sprite

type base struct {
	delay         uint64
	width, height float64
	inc           float64
	frames        uint64
	gif           *gif.GIF
	images        []image.Image
	counter       int
	len           int
	tex           two.Texture
	bananas
}

func (b *base) init() {
	b.delay = 6
	glc.ClearColor(128, 255, 0, 255)
	f, err := assets.Open("bana.gif")
	check(err)
	b.gif, err = gif.DecodeAll(f)
	check(err)
	b.len = len(b.gif.Image)
	b.images = make([]image.Image, b.len)
	for n, i := range b.gif.Image {
		dx, dy := i.Rect.Dx(), i.Rect.Dy()
		bounds := i.Bounds()
		rgba := image.NewNRGBA(image.Rect(0, 0, dx, dy))
		draw.Draw(rgba, bounds, i, image.Point{0, 0}, draw.Src)
		b.images[n] = rgba
	}
	b.tex, err = two.NewTexture(b.images[0])
	check(err)
	s := b.tex.NewSprite(int(b.width)/2, int(b.height)/2)
	s.SetLoc(b.width/2, b.height/4)
	b.bananas = append(b.bananas, s)
	s = b.tex.NewSprite(int(b.width)/2, int(b.height)/2)
	s.SetLoc(b.width/2, b.height-b.height/4)
	b.bananas = append(b.bananas, s)
	b.inc = 1
}

func (b *base) Update(float64) {
	b.frames++
	if b.tex == nil {
		b.init()
		return
	}
	x, y := b.bananas[0].Loc()
	x += b.inc
	switch {
	case x > 640-160:
		x = 640 - 160
		b.inc = -1
	case x < 160:
		x = 160
		b.inc = 1
	}
	b.bananas[0].SetLoc(x, y)
	_, y = b.bananas[1].Loc()
	b.bananas[1].SetLoc(x, y)
	if b.frames%b.delay == 0 {
		b.counter++
		if b.counter == b.len {
			b.counter = 0
		}
		b.tex.ReplaceImage(b.images[b.counter])
	}
}

func main() {
	log.Fatal(two.Run(&base{width: width, height: height}, nil, width, height))
}

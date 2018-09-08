// Rotate demonstrates
// rotating a uniform Sprite
/*
Usually also demonstrates monitor ghosting.
*/
package main

import (
	"image/color"
	"log"
	"math"

	"github.com/j7b/two"
	"github.com/j7b/two/glc"
	"github.com/j7b/two/window"
	"github.com/j7b/two/xtra"
)

func main() {
	const (
		width, height = 4, 4
		piece         = math.Pi / 30
	)
	glc.ClearColor(0, 0, 255, 0)
	u := two.NewUniform(color.RGBA{255, 0, 0, 255})
	s := u.NewSprite(width/2, height/2)
	s.SetLoc(width/2, height/2)
	Updater := func(f float64) {
		s.SetRot(s.Rot() + piece*f)
	}
	log.Fatal(two.Run(xtra.Func(Updater), window.Standard(384, 384), width, height))
}

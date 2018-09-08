// Orbit demonstrates integration with game dynamics.
/*
Inspired by https://github.com/jakecoffman/cp/blob/master/examples/planet/planet.go
*/
package main

import (
	"image"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/j7b/two"
	"github.com/j7b/two/doc/_examples/assets"
	. "github.com/j7b/two/doc/_examples/assets/_internal/cp"
	"github.com/j7b/two/glc"
	"github.com/j7b/two/xtra/atlas"
)

const (
	width  = 800
	height = 600
)

const gravityStrength = 5.0e6

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

type planet struct {
	width, height float64
	tex           two.Texture
	m             atlas.Map
	space         *Space
	planetBody    *Body
}

func (p *planet) update(b *Body) {
	if s, ok := b.UserData.(two.Sprite); ok {
		pos := b.Position()
		// s.SetLoc(float64(pos.X)+p.width/2, float64(pos.Y)+p.height/2)
		s.SetLoc(pos.X, pos.Y)
		s.SetRot(float64(b.Rotation().ToAngle()))
	}
}

func (p *planet) planetGravityVelocity(body *Body, gravity Vector, damping, dt float64) {
	pos := body.Position()
	sqdist := pos.LengthSq()
	g := pos.Mult(-gravityStrength / (sqdist * math.Sqrt(sqdist)))
	body.UpdateVelocity(g, damping, dt)
}

func (p *planet) addBox() {
	size := 10.0
	mass := 1.0

	verts := []Vector{
		{-size, -size},
		{-size, size},
		{size, size},
		{size, -size},
	}

	radius := Vector{size, size}.Length()
	pos := p.randPos(radius)

	body := p.space.AddBody(NewBody(mass, MomentForPoly(mass, len(verts), verts, Vector{}, 0)))
	body.SetVelocityUpdateFunc(p.planetGravityVelocity)
	body.SetPosition(pos)

	r := pos.Length()
	v := math.Sqrt(gravityStrength/r) / r
	body.SetVelocityVector(pos.Perp().Mult(v))

	body.SetAngularVelocity(v)
	body.SetAngle(math.Atan2(pos.Y, pos.X))

	sp := p.tex.NewSprite(20, 20)
	for n := range p.m {
		if n == "parrot.png" {
			continue
		}
		sp.SetSrc(p.m[n].Rect())
		if rand.Intn(2) == 0 {
			break
		}
	}
	body.UserData = sp

	shape := p.space.AddShape(NewPolyShape(body, 4, verts, NewTransformIdentity(), 0))
	shape.SetElasticity(0)
	shape.SetFriction(0.7)
}

func (p *planet) init() {
	f, err := assets.Open("kenney.nl/Spritesheet/round_nodetails_outline.png")
	check(err)
	img, _, err := image.Decode(f)
	check(err)
	p.tex, err = two.NewTexture(img)
	check(err)
	f, err = assets.Open("kenney.nl/Spritesheet/round_nodetails_outline.xml")
	check(err)
	p.m, err = atlas.Load(f)
	check(err)
	p.space = NewSpace()
	p.space.Iterations = 2

	p.planetBody = p.space.AddBody(NewKinematicBody())
	p.planetBody.SetAngularVelocity(0.2)
	sp := p.tex.NewSprite(140, 140)
	sp.SetSrc(p.m["parrot.png"].Rect())
	p.planetBody.UserData = sp

	for i := 0; i < 60; i++ {
		p.addBox()
	}

	shape := p.space.AddShape(NewCircle(p.planetBody, 70, Vector{}))
	shape.SetElasticity(1)
	shape.SetFriction(1)
}

func (p *planet) Update(time float64) {
	if p.planetBody == nil {
		p.init()
		return
	}
	p.space.EachBody(p.update)
	p.space.Step(time / 60)
}

func (p *planet) randPos(radius float64) Vector {
	var v Vector
	for {
		v = Vector{rand.Float64()*(float64(p.width)-2*radius) - (float64(p.width)/2 - radius), rand.Float64()*(float64(p.height)-2*radius) - (float64(p.height)/2 - radius)}
		if v.Length() >= 85 {
			return v
		}
	}
}

func main() {
	glc.ClearColor(0, 0, 0, 255)
	glc.Projection(mgl32.Ortho2D(width/2, -width/2, height/2, -height/2))
	log.Fatal(two.Run(&planet{width: width, height: height}, nil, width, height))
}

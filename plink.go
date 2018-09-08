// Plink demonstrates
// integration with game dynamics
// and event handling.
/*
Clicking on one of the guys freezes him.

Inspired by https://github.com/jakecoffman/cp/blob/master/examples/plink/plink.go
*/
package main

import (
	"image"
	"image/color"
	_ "image/png" // PNG ok golint?
	"log"
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/j7b/two/events/input/mouse"
	"github.com/j7b/two/geo"
	"github.com/j7b/two/glc"
	"github.com/j7b/two/window"

	"math/rand"

	"github.com/j7b/two"
	"github.com/j7b/two/doc/_examples/assets"
	"github.com/j7b/two/xtra/atlas"

	"github.com/j7b/two/doc/_examples/assets/_internal/cp"
	. "github.com/j7b/two/doc/_examples/assets/_internal/cp"
)

const (
	width  = 800
	height = 600
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var GRABBABLE_MASK_BIT uint = 1 << 31

var GrabFilter cp.ShapeFilter = cp.ShapeFilter{
	cp.NO_GROUP, GRABBABLE_MASK_BIT, GRABBABLE_MASK_BIT,
}
var NotGrabbableFilter cp.ShapeFilter = cp.ShapeFilter{
	cp.NO_GROUP, ^GRABBABLE_MASK_BIT, ^GRABBABLE_MASK_BIT,
}

var (
	pentagonMass   = 0.0
	pentagonMoment = 0.0
)

const numVerts = 5

const (
	_ = iota
	PinCollision
	GuyCollision
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type plink struct {
	width, height float64
	*Space
	pentagonMass   float64
	pentagonMoment float64
}

func (p *plink) Clicked(pos mouse.Position) {
	x, y := pos.Pos()
	nearest := p.Space.PointQueryNearest(Vector{X: float64(x), Y: float64(y)}, 0, GrabFilter)
	if nearest.Shape != nil {
		body := nearest.Shape.Body()
		if body.GetType() == BODY_STATIC {
			body.SetType(BODY_DYNAMIC)
			body.SetMass(p.pentagonMass)
			body.SetMoment(p.pentagonMoment)
			if sp, ok := body.UserData.(two.Sprite); ok {
				sp.SetHSVA(0, 0, 0, 0)
			}
		} else if body.GetType() == BODY_DYNAMIC {
			body.SetType(BODY_STATIC)
			if sp, ok := body.UserData.(two.Sprite); ok {
				sp.SetHSVA(0, -128, 127, 0)
			}
		}
	}
}

func (p *plink) Update(f float64) {
	if p.Space == nil {
		p.init()
		return
	}
	p.Space.EachBody(func(body *Body) {
		pos := body.Position()
		if pos.Y < -260 || math.Abs(pos.X) > 340 {
			x := rand.Float64()*640 - 320
			body.SetPosition(Vector{x, 260})
		}
		if sp, ok := body.UserData.(two.Sprite); ok {
			sp.SetLoc(float64(pos.X), float64(pos.Y))
			sp.SetRot(float64(body.Rotation().ToAngle()))
		}
	})
	p.Space.Step(f / 60.0)
}

func (p *plink) init() {
	glc.ClearColor(0, 255, 255, 255)
	width, height := float32(p.width), float32(p.height)
	glc.Projection(mgl32.Ortho2D(-width/2, width/2, -height/2, height/2))
	f, err := assets.Open("kenney.nl/Spritesheet/round_nodetails_outline.png")
	check(err)
	img, _, err := image.Decode(f)
	check(err)
	tex, err := two.NewTexture(img)
	check(err)
	f, err = assets.Open("kenney.nl/Spritesheet/round_nodetails_outline.xml")
	check(err)
	m, err := atlas.Load(f)
	check(err)
	uni := two.NewUniform(color.White)
	check(err)
	_, _ = tex, m

	space := NewSpace()
	p.Space = space
	space.Iterations = 2
	space.SetGravity(Vector{0, -100})

	var body *Body
	var shape *Shape

	tris := []Vector{
		{-6, -6},
		{0, 5},
		{6, -6},
	}

	for i := 0; i < 9; i++ {
		for j := 0; j < 6; j++ {
			stagger := (j % 2) * 40
			offset := Vector{float64(i*80 - 320 + stagger), float64(j*70 - 240)}
			shape = space.AddShape(NewPolyShape(space.StaticBody, 3, tris, NewTransformTranslate(offset), 0))
			shape.SetElasticity(1)
			shape.SetFriction(0.15)
			shape.SetFilter(NotGrabbableFilter)
			shape.SetCollisionType(PinCollision)
			s := uni.NewSprite(16, 16)
			if b, ok := s.(geo.Boxed); ok {
				q := b.Box()
				q.A.X += 8
				q.D.X -= 8
				b.SetBox(q)
			}
			s.SetLoc(float64(offset.X), float64(offset.Y))
			shape.UserData = s
		}
	}

	verts := []Vector{}
	for i := 0; i < numVerts; i++ {
		angle := -2.0 * math.Pi * float64(i) / numVerts
		verts = append(verts, Vector{20 * math.Cos(angle), 20 * math.Sin(angle)})
	}

	pentagonMass = 1.0
	p.pentagonMass = pentagonMass
	pentagonMoment = MomentForPoly(1, numVerts, verts, Vector{}, 0)
	p.pentagonMoment = pentagonMoment

	for i := 0; i < 40; i++ {
		body = space.AddBody(NewBody(pentagonMass, pentagonMoment))
		x := rand.Float64()*640 - 320
		body.SetPosition(Vector{x, 350})

		shape = space.AddShape(NewPolyShape(body, numVerts, verts, NewTransformIdentity(), 0))
		shape.SetElasticity(0)
		shape.SetFriction(0.4)
		shape.SetCollisionType(GuyCollision)
		sp := tex.NewSprite(32, 32)
		for n := range m {
			q := m[n].Rect()
			q.Max, q.Min = q.Min, q.Max
			sp.SetSrc(q)
			if rand.Intn(2) == 0 {
				break
			}
		}
		sp.SetLoc(float64(x), 350)
		body.UserData = sp
	}

	handler := space.NewCollisionHandler(PinCollision, GuyCollision)
	handler.BeginFunc = collide
}

func collide(arb *Arbiter, space *Space, _ interface{}) bool {
	s, _ := arb.Shapes()
	if sp, ok := s.UserData.(two.Sprite); ok {
		_, _, v, _ := sp.HSVA()
		if v > -100 {
			v -= 4
			sp.SetHSVA(0, 0, v, 0)
		}
	}
	return true
}

func main() {
	log.Fatal(two.Run(&plink{width: width, height: height}, window.Standard(width, height), 640, 480))
}

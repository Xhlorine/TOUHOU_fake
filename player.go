package main

import (
	"fmt"
	"math"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var (
	// The player-controled role
	Player *Character

	// Picture of the judging point
	point pixel.Picture

	// Picture of Reimu
	reimu pixel.Picture
)

func init() {
	var err error
	point, err = loadPNG(".\\pics\\point.png")
	if err != nil {
		panic(err)
	}
	reimu, err = loadPNG(".\\pics\\reimu.png")
	if err != nil {
		panic(err)
	}
	SpawnPlayer()
}
func SpawnPlayer() {
	Player = NewPlayer(5, 2, 0.6, pixel.V(WindowRange.W()/2, WindowRange.H()/4.0))
}

func UpdatePlayer(win *pixelgl.Window) {
	speed := Player.speedFast
	// Deal with movements
	if win.Pressed(pixelgl.KeyLeftShift) {
		speed = Player.speedSlow
	}
	if win.Pressed(pixelgl.KeyA) {
		speed = Player.speedSlowest
	}
	if win.Pressed(pixelgl.KeyUp) {
		Player.Velocity.Y += 1
	}
	if win.Pressed(pixelgl.KeyDown) {
		Player.Velocity.Y -= 1
	}
	if win.Pressed(pixelgl.KeyLeft) {
		Player.Velocity.X -= 1
	}
	if win.Pressed(pixelgl.KeyRight) {
		Player.Velocity.X += 1
	}
	if !Player.Velocity.Eq(pixel.ZV) {
		Player.Velocity = Player.Velocity.Unit().Scaled(speed)
	}
	// TODO
	if !MoveRange.Contains(Player.Position.Add(pixel.V(Player.Velocity.X, 0))) {
		Player.Velocity.X = 0
	}
	if !MoveRange.Contains(Player.Position.Add(pixel.V(0, Player.Velocity.Y))) {
		Player.Velocity.Y = 0
	}
	Player.Position = Player.Position.Add(Player.Velocity)
	Player.Velocity = pixel.ZV
	var hasHit bool = false
	for _, e := range emitters {
		for _, b := range e.bullets {
			if b.Pattern != none && !hasHit && isHit(b.Position, Player.Position, BulletList[b.Pattern].HitBoxRadius*b.Scale, Player.hitboxRadius) {
				if !b.Keepable {
					b.isDeleted = true
				}
				HitCount++
				// Avoid multi-hit in one frame
				hasHit = true
			}
			if b.Pattern != none && !b.grazed && isHit(b.Position, Player.Position, BulletList[b.Pattern].HitBoxRadius*b.Scale+1, Player.hitboxRadius+20) {
				b.grazed = true
				GrazeCount++
			}
		}
	}
	win.SetTitle("Hit: " + fmt.Sprint(HitCount) + ", Graze: " + fmt.Sprint(GrazeCount))
	Player.ChDraw(win)
}

type Character struct {
	age             int
	Position        pixel.Vec
	speedFast       float64
	speedSlow       float64
	speedSlowest    float64
	Velocity        pixel.Vec
	hitboxRadius    float64
	CharacterSprite *pixel.Sprite
	HitBoxSprite    *pixel.Sprite
}

func (p *Character) HitBox() pixel.Circle {
	return pixel.C(p.Position, p.hitboxRadius)
}

func (p *Character) HBDraw(t pixel.Target) {
	p.HitBoxSprite.Draw(t, pixel.IM.Scaled(pixel.ZV, 0.05*p.hitboxRadius).Rotated(pixel.ZV, math.Pi*float64(CurrentFrame)/60).Moved(p.Position))
}

func XOR(a, b bool) bool {
	if (a && !b) || (b && !a) {
		return true
	}
	return false
}

func (p *Character) ChDraw(win *pixelgl.Window) {
	p.age++
	basic := pixel.R(0, 208, 32, 256)
	ux := pixel.V(32, 0)
	uy := pixel.V(0, -48)
	if !XOR(win.Pressed(pixelgl.KeyLeft), win.Pressed(pixelgl.KeyRight)) {
		p.CharacterSprite = pixel.NewSprite(p.CharacterSprite.Picture(), basic.Moved(ux.Scaled(float64((p.age/6)%8))))
	}
	if win.JustPressed(pixelgl.KeyLeft) || win.JustPressed(pixelgl.KeyRight) {
		p.age = 0
	}
	if win.Pressed(pixelgl.KeyLeft) {
		x := 4
		if p.age < 16 {
			x = (p.age / 2) % 8
		} else {
			x += (p.age / 7) % 4
		}
		p.CharacterSprite = pixel.NewSprite(p.CharacterSprite.Picture(), basic.Moved(uy.Add(ux.Scaled(float64(x)))))
	}
	if win.Pressed(pixelgl.KeyRight) {
		x := 4
		if p.age < 16 {
			x = (p.age / 2) % 8
		} else {
			x += (p.age / 7) % 4
		}
		p.CharacterSprite = pixel.NewSprite(p.CharacterSprite.Picture(), basic.Moved(uy.Scaled(2).Add(ux.Scaled(float64(x)))))
	}
	p.CharacterSprite.Draw(win, pixel.IM.Moved(p.Position))
}

func NewPlayer(speedF float64, speedS float64, speedSest float64, position pixel.Vec) *Character {
	p := Character{
		Position:     position,
		speedFast:    speedF,
		speedSlow:    speedS,
		speedSlowest: speedSest,
		hitboxRadius: 3.0,
	}
	p.HitBoxSprite = pixel.NewSprite(point, point.Bounds())
	p.CharacterSprite = pixel.NewSprite(reimu, reimu.Bounds())
	return &p
}

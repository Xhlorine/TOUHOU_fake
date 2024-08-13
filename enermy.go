package main

import (
	"encoding/json"
	"io"
	"math"
	"math/rand"
	"os"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

const (
	tinyFairy_blue = iota
	tinyFairy_red
	tinyFairy_green
	tinyFairy_yellow

	smallFairy_blue
	smallFairy_red

	middleFairy_orange
	middleFairy_cyan

	bigFairy
)

var (
	EnemyPic  pixel.Picture
	EnemyList map[int]EnemyInfo
	enemys    []*Enemy
)

func init() {
	EnemyList = make(map[int]EnemyInfo)
	var err error
	EnemyPic, err = loadPNG(".\\pics\\enemy.png")
	if err != nil {
		panic(err)
	}

	file, err := os.Open("enemyList.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &EnemyList)
	if err != nil {
		panic(err)
	}

	a := tinyFairy_blue + tinyFairy_green + tinyFairy_red + tinyFairy_yellow + smallFairy_blue + smallFairy_red
	a += middleFairy_cyan + middleFairy_orange + bigFairy

	/*
		os.Remove("enemyList.json")
		f, err := os.Create("enemyList.json")
		if err != nil {
			panic(err)
		}
		d, err := json.Marshal(EnemyList)
		if err != nil {
			panic(err)
		}
		f.Write(d)
	*/
}

func UpdadeEnemy(win *pixelgl.Window) {
	for _, enemy := range enemys {
		enemy.Update()
		if !DeleteRange.Contains(enemy.Position) {
			enemy.Delete()
		}
		if enemy.isDeleted {
			for i, en := range enemys {
				if en == enemy {
					if i == len(enemys)-1 {
						enemys = enemys[:i]
					} else {
						enemys = append(enemys[:i], enemys[i+1:]...)
					}
					break
				}
			}
		}
	}
	for _, enemy := range enemys {
		enemy.Draw(win)
	}

}

type RigidBody struct {
	Age          int
	Position     pixel.Vec
	Velocity     pixel.Vec
	Acceleration pixel.Vec
}

func MoveToPos(start int, pos pixel.Vec, duration float64) func(*Enemy) {
	return func(r *Enemy) {
		if duration <= 0 {
			if r.Age == start {
				r.Position = pos
			}
		} else {
			if r.Age == start {
				r.Acceleration = r.Position.To(pos).Scaled(4.5 / math.Pow(duration, 2))
			} else if r.Age == start+int(duration/3) {
				r.Acceleration = pixel.ZV
			} else if r.Age == start+int(duration*2/3) {
				r.Acceleration = r.Position.To(pos).Scaled(-18 / math.Pow(duration, 2))
			} else if r.Age == start+int(duration) {
				r.Acceleration = pixel.ZV
				r.Velocity = pixel.ZV
			}
		}
	}
}

func MoveTowardsPos(start int, pos pixel.Vec, targetSpeed float64) func(*Enemy) {
	return func(r *Enemy) {
		if r.Age == start {
			r.Velocity = pixel.ZV
			r.Acceleration = r.Position.To(pos).Unit().Scaled(0.03)
		}
		if r.Velocity.Len() >= targetSpeed {
			r.Acceleration = pixel.ZV
		}
	}
}

func MoveTowards(start int, angle float64, targetSpeed float64) func(*Enemy) {
	return func(r *Enemy) {
		if r.Age == start {
			r.Velocity = pixel.ZV
			r.Acceleration = AtoV(targetSpeed, angle).Scaled(1.0 / 60.0)
		}
		if r.Velocity.Len() >= targetSpeed {
			r.Acceleration = pixel.ZV
		}
	}
}

func MoveRand(start int, rectangle pixel.Rect, duration float64, speedLimit float64) func(*Enemy) {
	return func(r *Enemy) {
		if r.Age < start || r.Age > start+int(duration) {
			return
		}
		if r.Velocity.Len() > speedLimit {
			r.Velocity = r.Velocity.Unit().Scaled(speedLimit * 0.9)
		}
		if (r.Age-start)%50 == 0 || !rectangle.Contains(r.Position) {
			r.Acceleration = AtoV((rand.Float64()+1)/49, r.Position.To(rectangle.Center()).Angle()+rand.Float64()-0.5)
		}
		if r.Age == start+int(duration) {
			r.Acceleration = pixel.ZV
		}
	}
}

type EnemyInfo struct {
	Name                string
	Bounds              []pixel.Rect
	DefaultHitboxRadius float64
}

type EnemyConfig struct {
	RigidBody
	Hp                  int
	Pattern             int
	EmitterForcedFollow bool
	Events              []func(e *Enemy)
}

func (e *EnemyConfig) AddEvent(ev func(*Enemy)) {
	e.Events = append(e.Events, ev)
}

type Enemy struct {
	RigidBody
	Hp int
	// Never change this unless you know what you're doing
	age                 int
	Pattern             int
	EmitterForcedFollow bool
	Events              []func(e *Enemy)
	sprite              *pixel.Sprite
	background          *pixel.Sprite
	emitter             *Emitter
	isDeleted           bool

	lastDir float64
}

func NewEnemy(cfg *EnemyConfig) *Enemy {
	e := Enemy{
		age:                 0,
		Hp:                  cfg.Hp,
		RigidBody:           cfg.RigidBody,
		Pattern:             cfg.Pattern,
		EmitterForcedFollow: cfg.EmitterForcedFollow,
	}
	e.lastDir = e.Velocity.X
	e.Events = cfg.Events
	e.sprite = pixel.NewSprite(EnemyPic, EnemyList[e.Pattern].Bounds[0])
	e.background = nil
	enemys = append(enemys, &e)
	return &e
}

func (e *Enemy) JustRedirect() bool {
	reop := func(e *Enemy) {
		e.lastDir = e.Velocity.X
	}
	defer reop(e)
	if e.lastDir*e.Velocity.X <= 0 && math.Abs(e.lastDir)+math.Abs(e.Velocity.X) != 0 {
		return true
	}
	return false
}

func (e *Enemy) Draw(t pixel.Target) {
	var x int
	HorizonalFlip := pixel.Matrix([]float64{-1, 0, 0, 1, 0, 0})
	if e.background != nil {
		e.background.Draw(t, pixel.IM.Moved(e.Position))
	}
	if e.JustRedirect() {
		e.age = 0
	}
	if e.Velocity.X == 0 {
		e.sprite.Set(EnemyPic, EnemyList[e.Pattern].Bounds[(e.age/6)%4])
		e.sprite.Draw(t, pixel.IM.Moved(e.Position))
	} else {
		if e.age < 24 {
			x = e.age/6 + 4
		} else {
			x = (e.age/6)%4 + 8
		}
		e.sprite.Set(EnemyPic, EnemyList[e.Pattern].Bounds[x])
		if e.Velocity.X > 0 {
			e.sprite.Draw(t, pixel.IM.Moved(e.Position))
		} else {
			e.sprite.Draw(t, HorizonalFlip.Moved(e.Position))
		}
	}
}

func (e *Enemy) Update() {
	e.Position = e.Position.Add(e.Velocity)
	e.Velocity = e.Velocity.Add(e.Acceleration)
	if e.EmitterForcedFollow && e.emitter != nil {
		e.emitter.Position = e.Position
	}
	for _, fun := range e.Events {
		fun(e)
	}
	if e.emitter == nil || e.emitter.isDeleted {
		e.emitter = nil
	}
	e.age++
	e.Age++
}

func (e *Enemy) SetEmitter(cfg *EmitterConfig) {
	e.emitter = NewEmitter(cfg)
	e.emitter.Position = e.Position
	e.emitter.Register()
}

func (e *Enemy) Delete() {
	e.isDeleted = true
	if e.emitter != nil {
		e.emitter.Delete()
	}
}

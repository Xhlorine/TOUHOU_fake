package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

var level int
var info [5]float64

func run() {
	cfg := pixelgl.WindowConfig{
		Bounds: ScaledRect(WindowRange, 1.2),
		VSync:  true,
	}

	// for {
	fmt.Println("请选择关卡：（1~29，输入0退出游戏）")
	fmt.Scanf("%d\n", &level)
	switch level {
	case 1:
		info = [5]float64{15, 5, 4, 0.15, 1.6}
	case 2:
		info = [5]float64{15, 6, 5, 0.28, 1.6}
	case 3:
		info = [5]float64{15, 5, 4, 0.8, 1.2}
	case 4:
		info = [5]float64{8, 9, 3, 0.05, 1.3}
	case 5:
		info = [5]float64{80, 17, 9, 1.3, 1.3}
	case 6:
		info = [5]float64{30, 20, 2, 0.6, 1}
	case 7:
		info = [5]float64{50, 16, 12, 0.6, 1.5}
	case 8:
		info = [5]float64{180, 36, 9, 1.2, 1.1}
	case 9:
		info = [5]float64{2, 2, 1, 0, 1.2}
	case 10:
		info = [5]float64{3, 3, 2, -2.4, 1.2}
	case 0:
		os.Exit(0)
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	win.SetMatrix(pixel.IM.Scaled(WindowRange.Center(), 1.2))
	sec := time.Tick(Tick)
	SpawnPlayer()
	if level >= 1 && level <= 10 {
		emitCFG := EmitterConfig{
			Pattern:   needle_blue_light,
			Interval:  uint32(info[0]),
			Lines:     uint32(info[1]),
			Layers:    uint32(info[2]),
			DeltaBV:   info[3],
			BVelocity: info[4],
			Radius:    8,
		}
		emitCFG.AddEvent(
			func(e *Emitter) {
				e.EmitAngle += 0.9
			},
		)
		emitCFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 100 {
					b.SetVelocityA(b.Velocity.Len(), b.Velocity.Angle()+math.Pi/6)
				}
			},
		)

		enemyCFG := EnemyConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleTop,
			},
			EmitterForcedFollow: true,
			Pattern:             bigFairy,
		}
		enemyCFG.AddEvent(MoveToPos(0, Window_MiddleUp, 100))
		//enemyCFG.AddEvent(MoveRand(100, WindowR_Up, 2400, 1))
		enemyCFG.AddEvent(MoveTowardsPos(2501, Window_MiddleTop, 1.2))
		en := NewEnemy(&enemyCFG)
		en.SetEmitter(&emitCFG)
		/*en.emitter.NewEmitterAttached(&sonCFG, false)*/
		en2 := NewEnemy(&enemyCFG)
		emitCFG.Pattern = needle_red_light
		emitCFG.events = []Event{
			func(e *Emitter) {
				e.EmitAngle -= 0.9
			},
		}
		emitCFG.bulletEvents = []BulletEvent{
			func(e *Emitter, b *Bullet) {
				if b.Age == 100 {
					b.SetVelocityA(b.Velocity.Len(), b.Velocity.Angle()-math.Pi/6)
				}
			},
		}
		en2.SetEmitter(&emitCFG)
	}
	switch level {
	case 11:
		emitCFG := EmitterConfig{
			Pattern:   needle_blue_light,
			Interval:  1,
			Lines:     2,
			Layers:    4,
			DeltaBV:   0.1,
			BVelocity: 1.7,
		}
		emitCFG.AddEvent(
			func(e *Emitter) {
				e.Position = Window_MiddleUp //.Add(pixel.V(1, 0).Scaled(60).Rotated(float64(e.Age) / 20))
				e.EmitAngle += float64(e.Age) / 465
			},
		)
		emit := NewEmitter(&emitCFG)
		emit.Register()
	case 12:
		emitCFG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Pattern:   smallstar_cyan_light,
			Interval:  4,
			Lines:     1,
			BVelocity: 2,
			Range:     math.Pi / 12,
			EmitAngle: -math.Pi / 2,
			Rotation:  1,
		}
		emitCFG.AddEvent(
			func(e *Emitter) {
				e.Position = pixel.V(rand.Float64()*WindowRange.W(), WindowRange.H()*1.1)
			},
		)
		emitCFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Velocity.Y > -0.2 {
					b.SetVelocityA(2.3, -math.Pi/2)
				}
				if b.Age == 0 {
					b.flag = b.Velocity.Angle()
					b.SetAcceleration(pixel.V((rand.Float64()-0.5)*0.01, 0.03))
				}
			},
		)
		emit := NewEmitter(&emitCFG)
		emit.Register()
	case 13:
		emitCFG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Pattern:   scale_red_light,
			Interval:  150,
			Lines:     144,
			Layers:    2,
			DeltaBV:   0.5,
			BVelocity: 1,
		}
		emitCFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				r := func(d float64) float64 {
					return d / math.Pi * 180
				}
				mod := func(a float64, b float64) float64 {
					for a >= b {
						a -= b
					}
					return a
				}
				if b.Age == 0 {
					b.SetVelocity(b.Velocity.Scaled(1.5 - math.Abs(float64(mod(r(b.Velocity.Angle()+math.Pi-e.Position.To(Player.Position).Angle()+15), 30)-15)/30)))
				}
			},
		)
		emit := NewEmitter(&emitCFG)
		emit.Register()
	case 14:
		baseCGG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleTop,
			},
			Pattern:   knife_red,
			Interval:  100,
			Lines:     1,
			Layers:    1,
			BVelocity: 3,
		}
		baseCGG.AddEvent(
			func(e *Emitter) {
				if e.Age%int(e.Interval) == 1 {
					v := func() pixel.Vec {
						u := Window_MiddleBottum.Rotated(math.Pi * rand.Float64())
						for WindowRange.Contains(Window_Middle.Add(u)) {
							u = u.Scaled(1.1)
						}
						return u
					}
					e.Position = Window_Middle.Add(v())
				}
				e.EmitAngle = e.Position.To(Player.Position).Angle() + (rand.Float64()-0.5)*0.1
			},
		)

		sonCFG := EmitterConfig{
			Pattern:       needle_cyan_dark,
			Interval:      20,
			Lines:         2,
			Layers:        1,
			BVelocity:     0.01,
			Radius:        8,
			ProtectRadius: 32,
		}
		sonCFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 50 {
					b.SetAccelerationA(0.01, b.Velocity.Angle())
				}
				if b.Age == 220 {
					b.SetAcceleration(pixel.ZV)
				}
			},
		)

		another := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Pattern:   circle_orange,
			Lines:     32,
			Interval:  100,
			BVelocity: 1.5,
		}
		another.AddEvent(
			func(e *Emitter) {
				e.EmitAngle = rand.Float64() * 2 * math.Pi
			},
		)

		emit := NewEmitter(&baseCGG)
		emit.NewEmitterAttached(&sonCFG, true)
		emit.Register()
		emit2 := NewEmitter(&another)
		emit2.Register()
	case 15:
		baseCGG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Pattern:   knife_red,
			Interval:  85,
			Lines:     2,
			Layers:    1,
			BVelocity: -3.5,
			Radius:    320,
			EmitAngle: math.Pi * 0.5,
		}
		baseCGG.AddEvent(
			func(e *Emitter) {
				e.Radius = 300 + rand.Int31n(32)
			},
		)

		sonCFG := EmitterConfig{
			Pattern:       needle_cyan_dark,
			Interval:      8,
			Lines:         2,
			Layers:        1,
			BVelocity:     1.2,
			DeltaBV:       0.3,
			Radius:        8,
			ProtectRadius: 32,
		}

		another := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Pattern:   circle_orange,
			Lines:     64,
			Interval:  70,
			BVelocity: 1.5,
		}
		another.AddEvent(
			func(e *Emitter) {
				e.EmitAngle = rand.Float64() * 2 * math.Pi
			},
		)

		emit := NewEmitter(&baseCGG)
		emit.NewEmitterAttached(&sonCFG, true)
		emit.Register()
		emit2 := NewEmitter(&another)
		emit2.Register()
	case 16:
		emitCFG := EmitterConfig{
			Pattern:   arrow_blue,
			Lines:     24,
			Interval:  180,
			DeltaBV:   0.3,
			Layers:    4,
			BVelocity: -0.8,
			Radius:    200,
		}
		emitCFG.AddEvent(
			func(e *Emitter) {
				e.Position = Player.Position
				e.EmitAngle = rand.Float64()
			},
		)
		emitCFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 40 {
					b.SetVelocity(b.Velocity.Scaled(2))
				}
				if b.Age == 90 {
					b.SetVelocityA(1.8, b.Velocity.Angle())
				}
			},
		)
		emit := NewEmitter(&emitCFG)
		emit.Register()
	case 17:
		baseCFG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Pattern:  none,
			Lines:    1,
			Layers:   1,
			Interval: 210,
		}
		targetedCFG := EmitterConfig{
			Life:      90,
			Pattern:   needle_red_light,
			Lines:     2,
			Layers:    3,
			Interval:  1,
			BVelocity: 1.2,
			DeltaBV:   1.1,
			Radius:    80,
		}
		targetedCFG.AddEvent(
			func(e *Emitter) {
				e.EmitAngle += 0.21
				x := float64(e.Age)
				e.Radius = int32(100*(1-math.Exp(-x)) + 16*math.Exp(-x/200)*math.Sin(x/2))
			},
		)
		targetedCFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 0 {
					b.flag = b.Velocity.Angle()
					b.SetVelocityA(b.Velocity.Len(), b.Position.To(Player.Position).Angle())
				}
				if b.Velocity.Len() < 2 && b.Age == 90 {
					b.ChangePattern(needle_blue_light)
					b.SetVelocityA(b.Velocity.Len()-0.4 /*rand.Float64()*2*math.Pi*/, b.flag)
				}

			},
		)
		base := NewEmitter(&baseCFG)
		base.NewEmitterAttached(&targetedCFG, false)
		base.Register()
	case 18:
		emitCFG := EmitterConfig{
			Pattern: grain_blue_light,
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Interval:  5,
			Lines:     2,
			BVelocity: 1,
			Layers:    6,
			DeltaBV:   2,
			Radius:    16,
		}
		emitCFG.AddEvent(
			func(e *Emitter) {
				e.EmitAngle += math.Pi * 2 / 3 / float64(e.Interval) / 8
			},
		)
		emitCFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 0 {
					b.SetAccelerationA(-0.06, b.Velocity.Angle())
					b.flag = b.Velocity.Angle()
				} else if b.Age == 60 {
					b.SetVelocityA(3.4, b.Position.To(Player.Position).Angle())
					b.SetAccelerationA(-0.05, b.Velocity.Angle())
					b.ChangePattern(grain_black)
				} else if b.Age == 120 {
					b.SetVelocityA(2.6, b.Position.To(e.Position).Angle()+math.Pi)
					b.SetAccelerationA(-0.04, b.Velocity.Angle())
					b.ChangePattern(grain_cyan_light)
				} else if b.Age == 190 {
					b.SetAcceleration(pixel.ZV)
					rand.Seed(int64(e.Age + 1))
					b.SetVelocityA(1.6, b.flag+0.0*math.Pi)
					b.ChangePattern(grain_red_light)
				}
				if math.Abs(b.Velocity.Len()) < 0.05 && b.Acceleration.Len() != 0 {
					b.SetAcceleration(pixel.ZV)
				}
			},
		)
		enmCFG := EnemyConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Pattern: bigFairy,
		}
		enm := NewEnemy(&enmCFG)
		enm.SetEmitter(&emitCFG)
	case 19:
		emitCFG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp.Sub(Window_MiddleLeftest.Scaled(0.2)),
			},
			Pattern:   needle_blue_light,
			Interval:  210,
			Lines:     60,
			Layers:    20,
			BVelocity: 5,
		}
		emitCFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Velocity.Len() < 0.02 && b.Acceleration.Len() > 0.001 {
					b.SetAcceleration(pixel.ZV)
					b.SetVelocity(pixel.ZV)
				}
				if b.Age == 0 {
					b.flag = b.Velocity.Angle()
					b.SetAccelerationA(0.08, b.Velocity.Angle()+math.Pi)
				} else if b.Age == 70 {
					angle := math.Pi*float64(2*b.Layer+2)/float64(e.Layers) + e.EmitAngle
					b.SetVelocityA(3, b.Position.To(e.Position.Add(AtoV(1000, angle))).Angle()+math.Pi/2)
					b.SetAccelerationA(0.05, b.Velocity.Angle()+math.Pi)
				} else if b.Age == 130 {
					b.SetVelocityA(3, b.Position.To(e.Position.Add(e.Position.To(Player.Position).Rotated(2*math.Pi*float64(b.Layer)/float64(e.Layers)))).Angle()+(rand.Float64()-0.5)*0.12)
					b.SetAcceleration(pixel.ZV)
				}
			},
		)
		emit := NewEmitter(&emitCFG)
		emit.Register()
	case 20:
		emitCFG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp.Sub(Window_MiddleLeftest.Scaled(0.2)),
			},
			Pattern:   needle_red_light,
			Interval:  210,
			Lines:     60,
			Layers:    18,
			BVelocity: 5,
		}
		emitCFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Velocity.Len() < 0.02 && b.Acceleration.Len() > 0.001 {
					b.SetAcceleration(pixel.ZV)
					b.SetVelocity(pixel.ZV)
				}
				if b.Age == 0 {
					b.flag = b.Velocity.Angle()
					b.SetAccelerationA(0.08, b.Velocity.Angle()+math.Pi)
				} else if b.Age == 70 {
					angle := math.Pi*float64(2*b.Layer+1)/float64(e.Layers) + e.EmitAngle
					b.SetVelocityA(3, b.Position.To(e.Position.Add(AtoV(400, angle))).Angle()+math.Pi/2)
					b.SetAccelerationA(0.05, b.Velocity.Angle()+math.Pi)
				} else if b.Age == 130 {
					b.SetVelocityA(2, b.Position.To(e.Position.Add(e.Position.To(Player.Position).Rotated(2*math.Pi*float64(b.Layer-4)/float64(e.Layers)))).Angle()+(rand.Float64()-0.5)*0.18)
					b.SetAcceleration(pixel.ZV)
				}
			},
		)
		emit := NewEmitter(&emitCFG)
		emit.Register()
	case 21:
		emit1CFG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Pattern:   middleball_red,
			Interval:  60,
			Lines:     28,
			BVelocity: 0.8,
		}
		emit1CFG.AddEvent(
			func(e *Emitter) {
				if e.Age%int(e.Interval) == 0 {
					e.EmitAngle += math.Pi / float64(e.Lines)
				}
			},
		)
		emit2CFG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Pattern:   grain_blue_light,
			Interval:  90,
			Lines:     50,
			Layers:    2,
			BVelocity: 1.4,
		}
		emit2CFG.AddEvent(
			func(e *Emitter) {
				if e.Age%int(e.Interval) == 0 {
					e.EmitAngle = rand.Float64()
					e.Position = Window_MiddleUp.Add(pixel.V((rand.Float64()-0.5)*400, (rand.Float64()-0.5)*40))
				}
			},
		)
		emit2CFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 60 {
					b.SetVelocityA(b.Velocity.Len(), b.Velocity.Angle()+(float64(b.Layer)-0.5)*math.Pi/1.5)
				}
			},
		)
		emit3CFG := EmitterConfig{
			RigidBody: emit1CFG.RigidBody,
			Pattern:   scale_purple_light,
			Interval:  12,
			Lines:     5,
			BVelocity: 1.9,
			Range:     math.Pi,
			EmitAngle: math.Pi / 2,
		}
		emit3CFG.AddEvent(
			func(e *Emitter) {
				if e.Age%60 == 0 {
					e.Position = Window_MiddleUp.Add(pixel.V((rand.Float64()-0.5)*400, (rand.Float64()-0.5)*40))
					e.EmitAngle = e.Position.To(Player.Position).Angle()
				}
			},
		)
		emit1 := NewEmitter(&emit1CFG)
		emit1.Register()
		emit2 := NewEmitter(&emit2CFG)
		emit2.Register()
		emit3 := NewEmitter(&emit3CFG)
		emit3.Register()
	case 22:
		emit2CFG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Pattern:   grain_blue_light,
			Interval:  90,
			Lines:     70,
			Layers:    4,
			BVelocity: 1.2,
		}
		emit2CFG.AddEvent(
			func(e *Emitter) {
				if e.Age%int(e.Interval) == 0 {
					e.EmitAngle += float64(e.Age) / 80
					//e.Position = Window_MiddleUp.Add(pixel.V((rand.Float64()-0.5)*400, (rand.Float64()-0.5)*40))
				}
			},
		)
		emit2CFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 60 {
					b.SetVelocityA(b.Velocity.Len(), b.Velocity.Angle()+(float64(b.Layer)-1.5)*math.Pi*1.5)
				}
			},
		)
		emit := NewEmitter(&emit2CFG)
		emit.Register()
	case 23:
		var sides float64 = 4
		emittCFG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Pattern:   none,
			Interval:  160,
			Lines:     uint32(sides),
			Layers:    1,
			BVelocity: 0.6,
			Radius:    80,
		}
		emittCFG.AddEvent(
			func(e *Emitter) {
				if e.Age%int(e.Interval) == 0 {
					e.EmitAngle += math.Pi * 2 / sides / 3
				}
			},
		)
		emitCFG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Life:      45,
			Pattern:   needle_red_light,
			Lines:     70,
			Layers:    4,
			DeltaBV:   0.16,
			Interval:  100,
			BVelocity: 1.2,
		}
		emitCFG.AddEvent(
			func(e *Emitter) {
				e.EmitAngle = rand.Float64()
			},
		)
		emitCFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 0 {
					b.SetVelocityA(b.Velocity.Len()/math.Cos(math.Pi/sides-Mod(b.Velocity.Rotated(e.Velocity.Angle()).Angle()+math.Pi, math.Pi/sides*2)), b.Velocity.Angle())
				}
				// if e.Age%150 == 0 {
				// 	if b.flag == 0 {
				// 		b.flag = 1
				// 		b.SetVelocityA(1.5, b.Velocity.Angle()+math.Pi)
				// 	}
				// }
			},
		)
		emit := NewEmitter(&emittCFG)
		emit.NewEmitterAttached(&emitCFG, true)
		emit.Register()
	case 24:
		emittCFG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Pattern:   none,
			Interval:  100,
			Lines:     8,
			Layers:    1,
			BVelocity: 0.6,
			Radius:    64,
		}
		emittCFG.AddEvent(
			func(e *Emitter) {
				e.EmitAngle += 0.53
			},
		)
		emit2CFG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Pattern:   grain_blue_light,
			Life:      100,
			Interval:  110,
			Lines:     32,
			Layers:    2,
			BVelocity: 1.2,
			Radius:    64,
		}
		emit2CFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 60 {
					b.SetVelocityA(b.Velocity.Len(), b.Velocity.Angle()+(float64(b.Layer)-0.5)*math.Pi*0.1)
				}
			},
		)
		emit := NewEmitter(&emittCFG)
		emit.NewEmitterAttached(&emit2CFG, true)
		emit.Register()
	case 25:
		baseCFG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Lines:    1,
			Interval: 5,
			Layers:   1,
			Pattern:  middleball_red,
			Radius:   60,
		}
		baseCFG.AddEvent(func(e *Emitter) { e.EmitAngle += 0.06335 })
		baseCFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 0 {
					b.flag = e.EmitAngle
				} else if b.Age == 20 {
					b.SetVelocityA(float64(b.Layer)/float64(e.Layers)*0.5+1, e.EmitAngle+math.Pi*float64(b.Layer))
				}
				if b.Velocity.Len() < 0.1 {
					b.Velocity = pixel.ZV
					b.Acceleration = pixel.ZV
				}
			},
		)
		sonCFG := EmitterConfig{
			Pattern:   smallball_orange,
			Interval:  200,
			Life:      50,
			Lines:     4,
			Layers:    3,
			BVelocity: 0,
			DeltaBV:   8,
			Range:     math.Pi * 2,
		}
		sonCFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 0 {
					b.SetAccelerationA(b.Velocity.Len()/20, b.Velocity.Angle()+math.Pi)
				} else if b.Age == 20 {
					b.flag = b.Position.To(b.belongsTo.Position).Len()
				} else if b.Age > 20 {
					b.Position = b.belongsTo.Position.Add(AtoV(b.flag, float64(b.Num)/float64(sonCFG.Lines)*sonCFG.Range+float64(b.Age)/120))
				}
				if b.Velocity.Len() < 0.1 && b.Acceleration.Len() != 0 {
					b.Velocity = pixel.ZV
					b.Acceleration = pixel.ZV
				}
			},
		)
		emit := NewEmitter(&baseCFG)
		emit.NewEmitterAttached(&sonCFG, true)
		emit.Register()
	case 26:
		emit1CFG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Pattern:   needle_blue_light,
			Interval:  9,
			Lines:     8,
			Layers:    2,
			BVelocity: 1.5,
			DeltaBV:   0.02,
			Radius:    60,
		}
		emit1CFG.AddEvent(func(e *Emitter) { e.EmitAngle -= 0.025 })
		emit1CFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 0 {
					b.SetVelocityA(b.Velocity.Len(), b.Velocity.Angle()+float64(e.Age)/60)
				}
			},
		)
		emit2CFG := emit1CFG
		emit2CFG.events = nil
		emit2CFG.bulletEvents = nil
		emit2CFG.Pattern = needle_red_light
		emit2CFG.AddEvent(func(e *Emitter) { e.EmitAngle += 0.025 })
		emit2CFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 0 {
					b.SetVelocityA(b.Velocity.Len(), b.Velocity.Angle()-float64(e.Age)/60)
				}
			},
		)
		emit1 := NewEmitter(&emit1CFG)
		emit1.Register()
		emit2 := NewEmitter(&emit2CFG)
		emit2.Register()
	case 27:
		emitCFG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Pattern:   ofuda_black,
			Interval:  60,
			Lines:     40,
			Layers:    16,
			BVelocity: 1.9,
			DeltaBV:   -0.7,
			Radius:    0,
		}
		emitCFG.AddEvent(
			func(e *Emitter) {
				e.EmitAngle += 0.0352
				e.Position = Window_MiddleUp.Add(pixel.V((rand.Float64()-0.5)*200, (rand.Float64()-0.5)*100))
			},
		)
		emitCFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 0 {
					switch b.Layer {
					case 0:
						b.ChangePattern(ofuda_red_light)
					case 1:
						b.ChangePattern(ofuda_orange)
					case 2:
						b.ChangePattern(ofuda_yellow_light)
					case 3:
						b.ChangePattern(ofuda_green_light)
					case 4:
						b.ChangePattern(ofuda_cyan_dark)
					case 5:
						b.ChangePattern(ofuda_blue_light)
					case 6:
						b.ChangePattern(ofuda_purple_light)
					}
					b.ChangePattern(ofuda_black + b.Layer)
				}
			},
		)
		emit := NewEmitter(&emitCFG)
		emit.Register()
	case 28:
		emitCFG := EmitterConfig{
			Pattern:   needle_red_light,
			Lines:     4,
			Layers:    1,
			BVelocity: 1.2,
			DeltaBV:   0.05,
			Interval:  3,
		}
		emitCFG.AddEvent(func(e *Emitter) {
			e.Position = Window_Middle.Add(Window_MiddleLeftest.Scaled(1.2).Rotated(math.Pi * float64(e.Age) / 150))
			e.EmitAngle += float64(int(e.Age/150)) / 37
		})
		emitCFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 0 {
					b.ChangePattern(b.Pattern + 2*b.Num)
				}
			},
		)
		NewEmitter(&emitCFG).Register()
	case 29:
		baseCFG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Interval:  40,
			Lines:     32,
			BVelocity: 1,
			Pattern:   circle_red_light,
		}
		baseCFG.AddEvent(
			func(e *Emitter) {
				if e.Age%int(e.Interval) == 0 {
					e.EmitAngle += math.Pi / float64(e.Lines)
				}
			},
		)
		NewEmitter(&baseCFG).Register()
		followCFG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Player.Position,
			},
			Pattern:   scale_blue_light,
			Interval:  5,
			Lines:     3,
			Layers:    1,
			BVelocity: -1.2,
			DeltaBV:   0.2,
			Radius:    256,
		}
		followCFG.AddEvent(
			func(e *Emitter) {
				e.Position = Player.Position
				e.EmitAngle += 0.007
			},
		)
		followCFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 0 {
					b.ChangePattern(b.Pattern + 2*b.Num)
				}
			},
		)
		NewEmitter(&followCFG).Register()
	case 30:
		emitCFG := EmitterConfig{
			RigidBody: RigidBody{
				Position: Player.Position,
			},
			Pattern:   scale_red_light,
			Interval:  6,
			BVelocity: 0.1,
			Radius:    256,
			Lines:     7,
		}
		emitCFG.AddEvent(
			func(e *Emitter) {
				e.Position = Player.Position
				e.EmitAngle += rand.Float64()*0.002 + 0.005
				if e.Age%250 == 150 {
					e.Interval = 500
				} else if e.Age%250 == 0 {
					e.Interval = 6
				}
			},
		)
		emitCFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 0 {
					b.ChangePattern(b.Pattern + b.Num)
				} else if b.Age == 50 {
					b.SetVelocityA(1.3, b.Velocity.Angle()+math.Pi-float64(e.Age%200)/600-0.1)
				}
			},
		)
		NewEmitter(&emitCFG).Register()
	case 31:
		cfg := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Pattern:   middleball_red,
			Interval:  5,
			Lines:     8,
			Layers:    1,
			BVelocity: 1.2,
			DeltaBV:   0.3,
		}
		cfg.AddEvent(
			func(e *Emitter) {
				e.EmitAngle += math.Pi / 180
				e.Position = Window_MiddleUp.Add(pixel.V(1, 0).Scaled(120 * math.Cos(float64(e.Age)/30)).Rotated(float64(e.Age) / 45))
			},
		)
		cfg.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 0 {
					b.ChangePattern(middleball_black + b.Num)
				} else if b.Age == 30 {
					b.Velocity = b.Velocity.Rotated(float64(e.Age) / 120)

				}
			},
		)
		NewEmitter(&cfg).Register()
	case 32:
		cfg := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Pattern:   needle_blue_dark,
			Lines:     4,
			Layers:    5,
			BVelocity: 0.8,
			DeltaBV:   0.5,
			Interval:  15,
		}
		flag := 1.0
		cfg.AddEvent(
			func(e *Emitter) {
				if e.Age%600 == 0 {
					flag = -flag
				}
				if e.Age%int(e.Interval) == 0 {
					e.EmitAngle -= math.Pi * 2 / float64(e.Interval) / 6 * flag
				}
			},
		)
		cfg.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 0 {
					b.flag = flag
					b.ChangePattern(b.Pattern + int(flag))
				}
				if b.Age == 150 {
					b.SetVelocity(b.Velocity.Rotated(math.Pi / float64(e.Layers) * float64(b.Layer) * b.flag))
				}
			},
		)
		NewEmitter(&cfg).Register()
		cfg.events = nil
		cfg.bulletEvents = nil
		cfg.AddEvent(
			func(e *Emitter) {
				if e.Age%600 == 0 {
					flag = -flag
				}
				if e.Age%int(e.Interval) == 0 {
					e.EmitAngle += math.Pi * 2 / float64(e.Interval) / 6 * flag
				}
			},
		)
		cfg.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 0 {
					b.flag = -flag
					b.ChangePattern(b.Pattern + int(-flag))
				}
				if b.Age == 150 {
					b.SetVelocity(b.Velocity.Rotated(math.Pi / float64(e.Layers) * float64(b.Layer) * b.flag))
				}
			},
		)
		NewEmitter(&cfg).Register()
	case 33:
		cfg := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_Middle,
			},
			Interval:  700,
			Lines:     12,
			Layers:    9,
			BVelocity: 0.01,
			DeltaBV:   8.1,
			Pattern:   bigstar_red,
			Rotation:  -2,
		}
		cfg.AddEvent(func(e *Emitter) {
			e.EmitAngle = e.Position.To(Player.Position).Angle()
			if e.Flag == 0 {
				e.Flag = 1
			}
			if e.Age%int(e.Interval) == 0 {
				e.Flag = -e.Flag
			}
		})
		cfg.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 0 {
					if rand.Float32() < 0.6 {
						b.isDeleted = true
					}
					if e.Flag == -1 {
						b.ChangePattern(bigstar_blue)
					}
				} else if b.Age == 40 {
					b.SetVelocity(b.Velocity.Scaled(0.01))
				} else if b.Age == 120 {
					b.SetVelocityA(1.1*e.Flag, Window_Middle.To(b.Position).Angle()+math.Pi/2)
				} else if b.Age > 120 && b.Age < 650 {
					b.SetAccelerationA(math.Pow(b.Velocity.Len(), 2)/Window_Middle.To(b.Position).Len()/1.01, b.Position.To(Window_Middle).Angle())
				} else if b.Age == 650 {
					b.SetAccelerationA(0.01, Window_Middle.To(b.Position).Angle())
				}
			},
		)
		ccfg := EmitterConfig{
			RigidBody: cfg.RigidBody,
			Pattern:   grain_yellow_light,
			Interval:  2,
			Lines:     6,
			BVelocity: 1.8,
			Radius:    int32(WindowRange.W() / 1.8),
		}
		ccfg.AddEvent(
			func(e *Emitter) {
				if e.Age%int(e.Interval) == 0 {
					e.EmitAngle -= math.Pi / 60
				}
			},
		)
		ccfg.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 0 {
					b.SetVelocity(b.Velocity.Rotated(math.Pi * (0.6 + 0.1*math.Sin(float64(e.Age)/60))))
				}
			},
		)
		cccfg := EmitterConfig{
			RigidBody: cfg.RigidBody,
			Pattern:   smallstar_blue_light,
			Interval:  3,
			Lines:     2,
			BVelocity: 1.2,
			Rotation:  4,
		}
		cccfg.AddEvent(func(e *Emitter) { e.EmitAngle = rand.Float64() * math.Pi * 2.01 })
		cccfg.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 0 {
					b.ChangePattern(smallstar_black + rand.Intn(16))
				}
			},
		)
		NewEmitter(&cccfg).Register()
		NewEmitter(&ccfg).Register()
		NewEmitter(&cfg).Register()
	case 34:
		emitCFG := EmitterConfig{
			Pattern: grain_blue_light,
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Interval:  8,
			Lines:     12,
			BVelocity: 0.8,
			Layers:    1,
			DeltaBV:   0.4,
			Radius:    16,
		}
		emitCFG.AddEvent(
			func(e *Emitter) {
				e.EmitAngle += math.Pi * 2 / 3 / float64(e.Interval) / 8
			},
		)
		emitCFG.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 0 {
					b.SetAccelerationA(-0.06, b.Velocity.Angle())
					b.flag = b.Velocity.Angle()
				} else if b.Age == 60 {
					b.SetVelocityA(3.4 /*b.Position.To(Player.Position).Angle()*/, e.EmitAngle*2)
					b.SetAccelerationA(-0.025, b.Velocity.Angle())
					b.ChangePattern(grain_black)
				} else if b.Age == 180 {
					b.SetVelocityA(2.9, b.Position.To(Player.Position).Angle())
					b.SetAccelerationA(-0.02, b.Velocity.Angle())
					b.ChangePattern(grain_cyan_light)
				} else if b.Age == 300 {
					b.SetVelocityA(2.8, b.Position.To(Player.Position).Angle())
					b.SetAccelerationA(-0.02, b.Velocity.Angle())
					b.ChangePattern(grain_orange)
				} else if b.Age == 420 {
					b.ChangePattern(grain_black)
					b.SetVelocityA(1.3, b.Position.To(e.Position.Add(e.Position.To(b.Position).Scaled(3))).Angle())
					b.SetAcceleration(pixel.ZV)
				}
				if math.Abs(b.Velocity.Len()) < 0.05 && b.Acceleration.Len() != 0 {
					b.SetAcceleration(pixel.ZV)
				}
			},
		)
		enmCFG := EnemyConfig{
			RigidBody: RigidBody{
				Position: Window_MiddleUp,
			},
			Pattern: bigFairy,
		}
		enm := NewEnemy(&enmCFG)
		enm.SetEmitter(&emitCFG)
	case 35:
		base := EmitterConfig{
			RigidBody: RigidBody{
				Position: Window_Middle,
			},
			Interval: 400,
			Pattern:  bigstar_red,
			Lines:    1,
			Keepable: true,
		}
		base.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Velocity.Len() < 0.2 {
					if b.Age < 320 {
						b.SetVelocityA(2.4, b.Position.To(Player.Position).Angle())
						b.SetAccelerationA(0.02, b.Velocity.Angle()+math.Pi)
					} else {
						b.SetVelocityA(1.7, b.Position.To(Player.Position).Angle())
					}

				}
			},
		)
		son := EmitterConfig{
			Pattern:   ofuda_blue_light,
			Interval:  12,
			Lines:     2,
			BVelocity: 0.08,
			Keepable:  false,
		}
		son.AddEvent(
			func(e *Emitter) {
				e.EmitAngle = e.Position.To(Player.Position).Angle() + math.Pi/2
			},
		)
		son.AddBulletEvent(
			func(e *Emitter, b *Bullet) {
				if b.Age == 120 {
					b.Velocity = b.Velocity.Scaled(20)
				}
			},
		)
		e := NewEmitter(&base)
		e.NewEmitterAttached(&son, false)
		e.Register()
	}
	// bullet := pixel.NewSprite(FullPic, FullPic.Bounds())
	// player := pixel.NewSprite(reimu, reimu.Bounds())
	// enermy := pixel.NewSprite(EnemyPic, EnemyPic.Bounds())
	for !win.Closed() {
		win.Clear(color.Black)
		// bullet.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.5).Moved(Window_UpLeft))
		// player.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.5).Moved(Window_UpRight))
		// enermy.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.5).Moved(Window_DownLeft))
		Update(win)
		win.Update()
		if win.JustPressed(pixelgl.KeyEscape) {
			STGClear(win)
			win.SetClosed(true)
			win.Destroy()
		}
		// } else if win.JustPressed(pixelgl.KeySpace) {
		// 	savePictureAsPNG(bullet.Picture(), "F:/Code/Go/TOUHOU/output/bullet.png")
		// 	savePictureAsPNG(Player.CharacterSprite.Picture(), "F:/Code/Go/TOUHOU/output/reimu.png")
		// 	savePictureAsPNG(enermy.Picture(), "F:/Code/Go/TOUHOU/output/enermy.png")
		// }
		<-sec
	}
	// }
}

// func savePictureAsPNG(pic pixel.Picture, path string) error {
// 	file, err := os.Create(path)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	err = png.Encode(file, pic.(*pixel.PictureData).Image())
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func main() {

	pixelgl.Run(run)

}

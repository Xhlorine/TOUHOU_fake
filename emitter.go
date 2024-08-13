package main

import (
	"encoding/json"
	"image"
	_ "image/png"
	"io"
	"math"
	"os"
	"time"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/pixelgl"
)

const (
	laser_black = iota
	laser_red_dark
	laser_red_light
	laser_purple_dark
	laser_purple_light
	laser_blue_dark
	laser_blue_light
	laser_cyan_dark
	laser_cyan_light
	laser_green_dark
	laser_green_light
	laser_chartreuse
	laser_yellow_dark
	laser_yellow_light
	laser_orange
	laser_white

	scale_black
	scale_red_dark
	scale_red_light
	scale_purple_dark
	scale_purple_light
	scale_blue_dark
	scale_blue_light
	scale_cyan_dark
	scale_cyan_light
	scale_green_dark
	scale_green_light
	scale_chartreuse
	scale_yellow_dark
	scale_yellow_light
	scale_orange
	scale_white

	circle_black
	circle_red_dark
	circle_red_light
	circle_purple_dark
	circle_purple_light
	circle_blue_dark
	circle_blue_light
	circle_cyan_dark
	circle_cyan_light
	circle_green_dark
	circle_green_light
	circle_chartreuse
	circle_yellow_dark
	circle_yellow_light
	circle_orange
	circle_white

	smallball_black
	smallball_red_dark
	smallball_red_light
	smallball_purple_dark
	smallball_purple_light
	smallball_blue_dark
	smallball_blue_light
	smallball_cyan_dark
	smallball_cyan_light
	smallball_green_dark
	smallball_green_light
	smallball_chartreuse
	smallball_yellow_dark
	smallball_yellow_light
	smallball_orange
	smallball_white

	grain_black
	grain_red_dark
	grain_red_light
	grain_purple_dark
	grain_purple_light
	grain_blue_dark
	grain_blue_light
	grain_cyan_dark
	grain_cyan_light
	grain_green_dark
	grain_green_light
	grain_chartreuse
	grain_yellow_dark
	grain_yellow_light
	grain_orange
	grain_white

	chain_black
	chain_red_dark
	chain_red_light
	chain_purple_dark
	chain_purple_light
	chain_blue_dark
	chain_blue_light
	chain_cyan_dark
	chain_cyan_light
	chain_green_dark
	chain_green_light
	chain_chartreuse
	chain_yellow_dark
	chain_yellow_light
	chain_orange
	chain_white

	needle_black
	needle_red_dark
	needle_red_light
	needle_purple_dark
	needle_purple_light
	needle_blue_dark
	needle_blue_light
	needle_cyan_dark
	needle_cyan_light
	needle_green_dark
	needle_green_light
	needle_chartreuse
	needle_yellow_dark
	needle_yellow_light
	needle_orange
	needle_white

	ofuda_black
	ofuda_red_dark
	ofuda_red_light
	ofuda_purple_dark
	ofuda_purple_light
	ofuda_blue_dark
	ofuda_blue_light
	ofuda_cyan_dark
	ofuda_cyan_light
	ofuda_green_dark
	ofuda_green_light
	ofuda_chartreuse
	ofuda_yellow_dark
	ofuda_yellow_light
	ofuda_orange
	ofuda_white

	bullet_black
	bullet_red_dark
	bullet_red_light
	bullet_purple_dark
	bullet_purple_light
	bullet_blue_dark
	bullet_blue_light
	bullet_cyan_dark
	bullet_cyan_light
	bullet_green_dark
	bullet_green_light
	bullet_chartreuse
	bullet_yellow_dark
	bullet_yellow_light
	bullet_orange
	bullet_white

	blackgrain_black
	blackgrain_red_dark
	blackgrain_red_light
	blackgrain_purple_dark
	blackgrain_purple_light
	blackgrain_blue_dark
	blackgrain_blue_light
	blackgrain_cyan_dark
	blackgrain_cyan_light
	blackgrain_green_dark
	blackgrain_green_light
	blackgrain_chartreuse
	blackgrain_yellow_dark
	blackgrain_yellow_light
	blackgrain_orange
	blackgrain_white

	smallstar_black
	smallstar_red_dark
	smallstar_red_light
	smallstar_purple_dark
	smallstar_purple_light
	smallstar_blue_dark
	smallstar_blue_light
	smallstar_cyan_dark
	smallstar_cyan_light
	smallstar_green_dark
	smallstar_green_light
	smallstar_chartreuse
	smallstar_yellow_dark
	smallstar_yellow_light
	smallstar_orange
	smallstar_white

	drip_black
	drip_red_dark
	drip_red_light
	drip_purple_dark
	drip_purple_light
	drip_blue_dark
	drip_blue_light
	drip_cyan_dark
	drip_cyan_light
	drip_green_dark
	drip_green_light
	drip_chartreuse
	drip_yellow_dark
	drip_yellow_light
	drip_orange
	drip_white

	heart_black
	heart_red
	heart_purple
	heart_blue
	heart_cyan
	heart_green
	heart_yellow
	heart_white

	arrow_black
	arrow_red
	arrow_purple
	arrow_blue
	arrow_cyan
	arrow_green
	arrow_yellow
	arrow_white

	bigstar_black
	bigstar_red
	bigstar_purple
	bigstar_blue
	bigstar_cyan
	bigstar_green
	bigstar_yellow
	bigstar_white

	middleball_black
	middleball_red
	middleball_purple
	middleball_blue
	middleball_cyan
	middleball_green
	middleball_yellow
	middleball_white

	butterfly_black
	butterfly_red
	butterfly_purple
	butterfly_blue
	butterfly_cyan
	butterfly_green
	butterfly_yellow
	butterfly_white

	knife_black
	knife_red
	knife_purple
	knife_blue
	knife_cyan
	knife_green
	knife_yellow
	knife_white

	ellipse_black
	ellipse_red
	ellipse_purple
	ellipse_blue
	ellipse_cyan
	ellipse_green
	ellipse_yellow
	ellipse_white

	bigball_red
	bigball_blue
	bigball_green
	bigball_yellow

	none
)

var (
	// The picture of the HitBox
	Pic pixel.Picture

	FullPic pixel.Picture

	// Namely the fps. Never change it unless you know where it's used
	Fps int64 = 60

	// The duration of a frame. Never change it unless you know where it's used
	Tick time.Duration = time.Second / time.Duration(Fps)

	// The collection of all registered emitters. Invisible outside the Package.
	emitters []*Emitter = nil

	// The collection of all the bullets.

	// Frames since "Update" was first called
	CurrentFrame = 0

	// The max number of bullets
	BulletLimit = 5000

	// Current number of bullets
	CurrentBullet = 0

	// The times (actually frames) hit the bullets
	HitCount = 0

	GrazeCount = 0

	BulletList map[int]BulletInfo

	WindowRange = pixel.R(0, 0, 576, 672)
	MoveRange   = ScaledRect(WindowRange, 0.95)
	DeleteRange = ScaledRect(WindowRange, 1.8)

	Window_Middle      = WindowRange.Center()
	Window_MiddleUp    = pixel.V(WindowRange.W()*0.5, WindowRange.H()*0.75)
	Window_MiddleDown  = pixel.V(WindowRange.W()*0.5, WindowRange.H()*0.25)
	Window_MiddleLeft  = pixel.V(WindowRange.W()*0.25, WindowRange.H()*0.5)
	Window_MiddleRight = pixel.V(WindowRange.W()*0.75, WindowRange.H()*0.5)
	Window_UpLeft      = pixel.V(WindowRange.W()*0.25, WindowRange.H()*0.75)
	Window_UpRight     = pixel.V(WindowRange.W()*0.75, WindowRange.H()*0.75)
	Window_DownLeft    = pixel.V(WindowRange.W()*0.25, WindowRange.H()*0.25)
	Window_DownRight   = pixel.V(WindowRange.W()*0.75, WindowRange.H()*0.25)

	Window_MiddleTop      = pixel.V(WindowRange.W()*0.5, WindowRange.H())
	Window_MiddleBottum   = pixel.V(WindowRange.W()*0.5, 0)
	Window_MiddleLeftest  = pixel.V(0, WindowRange.H()*0.5)
	Window_MiddleRightest = pixel.V(WindowRange.W(), WindowRange.H()*0.5)

	WindowR_Up = pixel.Rect{
		Min: Window_MiddleLeftest,
		Max: WindowRange.Max,
	}
)

// An action of an entity
type Event func(e *Emitter)

// An action of a bullet, especially emitted by the emitter
type BulletEvent func(*Emitter, *Bullet)

// The full illustration of a picture
type PicInfo struct {
	Img  image.Image
	Name string
}

type BulletInfo struct {
	Name         string
	Bounds       pixel.Rect
	HitBoxRadius float64
}

// Initialize the picture, create the player
func init() {
	var err error
	BulletList = make(map[int]BulletInfo)

	FullPic, err = loadPNG(".\\pics\\bullets.png")
	if err != nil {
		panic(err)
	}

	file, err := os.Open("bulletList.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &BulletList)
	if err != nil {
		panic(err)
	}
	// No use, just stop warnings
	CurrentBullet *= scale_black + circle_black + smallball_black + laser_black + grain_black + chain_black + needle_black + ofuda_black
	CurrentBullet *= bullet_black + blackgrain_black + smallstar_black + drip_black + heart_black + arrow_black + middleball_black
	CurrentBullet *= butterfly_black + knife_black + ellipse_black + bigball_blue + bigstar_black + none
	/*os.Remove("bulletList.json")
	f, err := os.Create("bulletList.json")
	if err != nil {
		panic(err)
	}
	d, err := json.Marshal(BulletList)
	if err != nil {
		panic(err)
	}
	f.Write(d)*/
}

// You'd better know what you are changing...
func ChangeWindowRange(rect pixel.Rect) {
	WindowRange = rect
	MoveRange = ScaledRect(WindowRange, 0.95)
	DeleteRange = ScaledRect(WindowRange, 1.5)

	Window_Middle = WindowRange.Center()
	Window_MiddleUp = pixel.V(WindowRange.W()*0.5, WindowRange.H()*0.75)
	Window_MiddleDown = pixel.V(WindowRange.W()*0.5, WindowRange.H()*0.25)
	Window_MiddleLeft = pixel.V(WindowRange.W()*0.25, WindowRange.H()*0.5)
	Window_MiddleRight = pixel.V(WindowRange.W()*0.75, WindowRange.H()*0.5)
	Window_UpLeft = pixel.V(WindowRange.W()*0.25, WindowRange.H()*0.75)
	Window_UpRight = pixel.V(WindowRange.W()*0.75, WindowRange.H()*0.75)
	Window_DownLeft = pixel.V(WindowRange.W()*0.25, WindowRange.H()*0.25)
	Window_DownRight = pixel.V(WindowRange.W()*0.75, WindowRange.H()*0.25)

	Window_MiddleTop = pixel.V(WindowRange.W()*0.5, WindowRange.H())
	Window_MiddleBottum = pixel.V(WindowRange.W()*0.5, 0)
	Window_MiddleLeftest = pixel.V(0, WindowRange.H()*0.5)
	Window_MiddleRightest = pixel.V(WindowRange.W(), WindowRange.H()*0.5)
}

// Read the give picture and return its "pixel.Picture" form
func loadPNG(path string) (pixel.Picture, error) {
	// Open the file and read it
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // Remerber to close the file

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	// Covert the image and return
	return pixel.PictureDataFromImage(img), nil
}

func Mod(devidend float64, devisor float64) float64 {
	return devidend - float64(int32(devidend/devisor))*devisor
}

// A closure.
// Used to update all the emitters and bullets, which otherwise stay still
// Update the player
func Update(win *pixelgl.Window) {
	CurrentFrame++
	for _, e := range emitters {
		e.Update()
		e.Spawn()

		if e.isDeleted && e.bullets == nil {
			for i, em := range emitters {
				if em == e {
					if i == len(emitters)-1 {
						emitters = emitters[:i]
					} else {
						emitters = append(emitters[:i], emitters[i+1:]...)
					}
					break
				}
			}
		}
	}

	UpdadeEnemy(win)
	UpdatePlayer(win)
	for _, e := range emitters {
		e.Draw(win)
	}
	if win.Pressed(pixelgl.KeyLeftShift) || win.Pressed(pixelgl.KeyA) {
		Player.HBDraw(win)
	}
}

func STGClear(win *pixelgl.Window) {
	for _, e := range emitters {
		e.bullets = nil
		e.Delete()
	}
	for _, e := range enemys {
		e.isDeleted = true
	}

	Update(win)
	emitters = nil
	CurrentFrame = 0
	BulletLimit = 5000
	CurrentBullet = 0
	HitCount = 0
	Player = nil
}

// Change a vector from "mod + angle" to "X + Y"
func AtoV(mod, angle float64) pixel.Vec {
	return pixel.V(mod*math.Cos(angle), mod*math.Sin(angle))
}

func isHit(p1, p2 pixel.Vec, r1, r2 float64) bool {
	dist := func(u, v pixel.Vec) float64 {
		return u.To(v).Len()
	}
	if math.Pow(dist(p1, p2), 2) < math.Pow(r1, 2)+math.Pow(r2, 2) {
		return true
	} else {
		return false
	}
}

func ScaledRect(rect pixel.Rect, t float64) pixel.Rect {
	min := rect.Center().Add(rect.Center().To(rect.Min).Scaled(t))
	max := rect.Center().Add(rect.Center().To(rect.Max).Scaled(t))
	return pixel.R(min.X, min.Y, max.X, max.Y)
}

// Used to illustrate a bullet.
type Bullet struct {
	RigidBody
	Life   int
	Sprite *pixel.Sprite
	// The scale of the picture and hitbox, times of the picture
	Scale float64

	// Only use it when the bullet doesn't come from an emitter
	BulletEvents []BulletEvent

	belongsTo *Bullet
	grazed    bool

	// Never Change it directly, use ChangePattern() instead
	Pattern   int
	flag      float64
	Layer     int
	Num       int
	isDeleted bool
	Keepable  bool
	Rotation  int
}

// Create a bullet
func NewBullet(pt int, belong *Bullet, rot int, keep bool) *Bullet {
	var b Bullet = Bullet{
		Scale: 1,
		RigidBody: RigidBody{
			Age: 0,
		},
		Pattern:  pt,
		flag:     0,
		grazed:   false,
		Keepable: keep,
	}
	b.belongsTo = belong
	b.Rotation = rot
	b.Sprite = pixel.NewSprite(FullPic, BulletList[b.Pattern].Bounds)
	CurrentBullet++
	return &b
}

func (b *Bullet) ChangePattern(s int) {
	b.Pattern = s
	b.Sprite = pixel.NewSprite(FullPic, BulletList[b.Pattern].Bounds)
}

func (b *Bullet) SetVelocity(v pixel.Vec) {
	b.Velocity = v
}

func (b *Bullet) SetVelocityA(speed, angle float64) {
	b.Velocity = pixel.V(speed*math.Cos(angle), speed*math.Sin(angle))
}

func (b *Bullet) SetAcceleration(acc pixel.Vec) {
	b.Acceleration = acc
}

func (b *Bullet) SetAccelerationA(acc, angle float64) {
	b.Acceleration = pixel.V(acc*math.Cos(angle), acc*math.Sin(angle))
}

// Return the Hitbox of a bullet, simply a circle
func (b *Bullet) HitBox() pixel.Circle {
	return pixel.C(b.Position, 30*b.Scale)
}

// Update the position and velocity
// For more complex calculation, please use functions
func (b *Bullet) Update() {
	b.Position = b.Position.Add(b.Velocity.Scaled(1))
	b.Velocity = b.Velocity.Add(b.Acceleration.Scaled(1))
	if b.Life > 0 {
		if b.Age > b.Life {
			b.isDeleted = true
		}
	}
	b.Age++
}

// Return the matrix illustrating the scale and position of the sprite
func (b *Bullet) matrix() pixel.Matrix {
	if b.Rotation == 0 {
		return pixel.IM.Scaled(pixel.ZV, b.Scale).Rotated(pixel.ZV, b.Velocity.Angle()-math.Pi/2).Moved(b.Position)
	} else {
		return pixel.IM.Scaled(pixel.ZV, b.Scale).Rotated(pixel.ZV, math.Pi*float64(b.Age)/180*float64(b.Rotation)).Moved(b.Position)
	}
}

// Draw the bullet on a given target
func (b *Bullet) Draw(t pixel.Target) {
	if BulletList[b.Pattern].Name != "none" {
		b.Sprite.Draw(t, b.matrix())
	}

}

type EmitterConfig struct {
	RigidBody
	IsAttaching   bool
	BirthFrame    int
	Life          int
	Lines         uint32
	Interval      uint32
	Radius        int32
	EmitAngle     float64
	Range         float64
	Layers        uint32
	Pattern       int
	events        []Event
	bulletEvents  []BulletEvent
	BLife         int
	BVelocity     float64
	DeltaBV       float64
	ProtectRadius float64
	Rotation      int
	Keepable      bool
}

func (ec *EmitterConfig) AddEvent(ev Event) *EmitterConfig {
	ec.events = append(ec.events, ev)
	return ec
}

func (ec *EmitterConfig) AddBulletEvent(bev BulletEvent) *EmitterConfig {
	ec.bulletEvents = append(ec.bulletEvents, bev)
	return ec
}

type Emitter struct {
	RigidBody
	EmitAngle float64
	// The lifespan of the emitter. 0 and negative numbers for no passive delete
	Life int

	// The number of bullets emitted in one go
	Lines uint32

	// The interval of two emissions. Uint: Tick
	Interval uint32

	// The initial speed of the bullets. Unit: Pixel per Tick
	BVelocity float64
	BLife     int

	// The initial acceleration of the bullets. Unit: Pixel per Square Tick
	BAcceleration pixel.Vec
	Rotation      int

	// The distance of the bullets emitted from the emitter.
	Radius        int32
	ProtectRadius float64

	Layers  uint32
	DeltaBV float64

	// The range of the bullets emitted
	Range float64

	isDeleted  bool
	Keepable   bool
	following  *Emitter
	isRelative bool

	Flag float64

	Pattern      int
	bullets      []*Bullet
	container    *pixel.Batch
	events       []Event
	bulletEvents []BulletEvent
}

// Create a new emitter
// When it comes to "lines", "interval", negative number for default
// The others would be assigned directly
func NewEmitter(cfg *EmitterConfig) *Emitter {
	e := Emitter{
		Lines:         cfg.Lines,
		Interval:      cfg.Interval,
		Radius:        cfg.Radius,
		BVelocity:     cfg.BVelocity,
		RigidBody:     cfg.RigidBody,
		Pattern:       cfg.Pattern,
		Range:         cfg.Range,
		Life:          cfg.Life,
		Layers:        cfg.Layers,
		BLife:         cfg.BLife,
		DeltaBV:       cfg.DeltaBV,
		ProtectRadius: cfg.ProtectRadius,
		EmitAngle:     cfg.EmitAngle,
		Rotation:      cfg.Rotation,
		Keepable:      cfg.Keepable,
	}
	e.container = pixel.NewBatch(&pixel.TrianglesData{}, FullPic)
	if e.Lines <= 0 {
		e.Lines = 0
	}
	if e.Interval <= 0 {
		e.Interval = 60
	}
	if e.Range <= 0 {
		e.Range = 2 * math.Pi
	}
	if e.Layers <= 0 {
		e.Layers = 1
	}
	e.events = cfg.events
	e.bulletEvents = cfg.bulletEvents
	return &e
}

// After calling this function, resettings of position, velocity, acceleration over the emitter are useless
func (t *Emitter) NewEmitterAttached(cfg *EmitterConfig, angleRelative bool) {
	e := NewEmitter(cfg)
	e.following = t
	e.isRelative = angleRelative
	e.Register()
}

// Emit the bullets once
// Never call it manually
func (e *Emitter) Spawn() {
	if e.isDeleted {
		return
	}
	underProtection := func() bool {
		if e.ProtectRadius <= 0 {
			return false
		} else if pixel.C(e.Position, e.ProtectRadius).Contains(Player.Position) {
			return true
		} else {
			return false
		}
	}
	spawnOnce := func(e *Emitter, belong *Bullet) {
		for i := 0; i < int(e.Lines); i++ {
			angle := e.Range*(float64(2*i+1)/float64(2*e.Lines)-0.5) + e.EmitAngle
			pos := e.Position.Add(AtoV(float64(e.Radius), angle))
			for j := 0; j < int(e.Layers); j++ {
				if CurrentBullet >= BulletLimit {
					continue
				}
				b := NewBullet(e.Pattern, belong, e.Rotation, e.Keepable)
				b.Position = pos
				b.Layer = j
				b.Num = i
				if e.Layers == 1 {
					b.SetVelocityA(e.BVelocity, angle)
				} else {
					b.SetVelocityA(e.BVelocity+e.DeltaBV*(float64(j)/float64(e.Layers-1)), angle)
				}
				b.SetAcceleration(e.BAcceleration)
				e.bullets = append(e.bullets, b)
				for _, fun := range e.bulletEvents {
					fun(e, b)
				}
			}

		}
	}
	if e.following == nil {
		if CurrentFrame%int(e.Interval) == 0 && !underProtection() {
			spawnOnce(e, nil)
		}
	} else {
		for _, bu := range e.following.bullets {
			e.Position = bu.Position
			e.Velocity = bu.Velocity
			e.Acceleration = bu.Acceleration
			e.Age = bu.Age
			if e.isRelative {
				e.EmitAngle = bu.Velocity.Angle()
			}
			if !(e.Life > 0 && e.Age > e.Life) && bu.Age%int(e.Interval) == 0 && !underProtection() {
				spawnOnce(e, bu)
			}
		}
	}

}

func (e *Emitter) SetVelocity(v pixel.Vec) {
	e.Velocity = v
}

func (e *Emitter) SetVelocityA(speed, angle float64) {
	e.Velocity = pixel.V(speed*math.Cos(angle), speed*math.Sin(angle))
}

func (e *Emitter) SetAcceleration(acc pixel.Vec) {
	e.Acceleration = acc
}

func (e *Emitter) SetAccelerationA(acc, angle float64) {
	e.Acceleration = pixel.V(acc*math.Cos(angle), acc*math.Sin(angle))
}

// Update according to the Velocity, Acceleration, and events
// Sequece: Calculation of position, velocity; delete the bullets outside the window; implement the events
// Never call it manually
func (e *Emitter) Update() {
	e.Position = e.Position.Add(e.Velocity.Scaled(1))
	e.Velocity = e.Velocity.Add(e.Acceleration.Scaled(1))
	if e.following == nil {
		e.Age++
		if e.Life > 0 && e.Age > e.Life {
			e.isDeleted = true
		}
	} else {
		if e.following.bullets == nil && e.following.isDeleted {
			e.isDeleted = true
			e.following = nil
		}
	}

	for _, fun := range e.events {
		fun(e)
	}
	// Basic Update
	for i, b := range e.bullets {
		b.Update()
		if !DeleteRange.Contains(b.Position) {
			b.belongsTo = nil
			e.bullets[i] = nil
		}
	}
	// Recalculate the EmitAngle
	if e.EmitAngle > 0 {
		for e.EmitAngle > math.Pi*2 {
			e.EmitAngle -= math.Pi * 2
		}
	} else {
		for e.EmitAngle < -math.Pi*2 {
			e.EmitAngle += 2 * math.Pi
		}
	}
	// Delete bullets outside the window
	for i := 0; i < len(e.bullets); i++ {
		if e.bullets[i] != nil {
			if !e.bullets[i].isDeleted {
				continue
			}
		}
		if i == len(e.bullets)-1 {
			e.bullets = e.bullets[:i]
		} else {
			e.bullets = append(e.bullets[:i], e.bullets[(i+1):]...)
		}
		i--
		CurrentBullet--
	}
	// Implement the events
	for _, fun := range e.bulletEvents {
		for _, b := range e.bullets {
			if e.following != nil {
				fun(e.following, b)
			} else {
				fun(e, b)
			}
		}
	}
}

// // Add events of the emitter
// func (e *Emitter) AddEvent(event Event) {
// 	e.events = append(e.events, event)
// }

// // Add events of the bullets
// func (e *Emitter) AddBulletEvent(bEvent BulletEvent) {
// 	e.bulletEvents = append(e.bulletEvents, bEvent)
// }

// Draw all the bullets from one emitter
func (e *Emitter) Draw(t pixel.Target) {
	e.container.Clear()
	for _, b := range e.bullets {
		b.Draw(e.container)
	}
	e.container.Draw(t)
}

// Call this before registration
func (e *Emitter) Copy() *Emitter {
	if e.Age != 0 {
		return nil
	}
	em := *e
	return &em
}

// Register the emitter
func (e *Emitter) Register() {
	emitters = append(emitters, e)
}

// Delete the emmiter
func (e *Emitter) Delete() {
	e.isDeleted = true
}

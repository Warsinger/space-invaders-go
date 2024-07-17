package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

type PlayerData struct {
	Score int
}

var Player = donburi.NewComponentType[PlayerData]()

func NewPlayer(w donburi.World) error {
	entity := w.Create(Player, Position, Velocity, Render)
	entry := w.Entry(entity)
	Position.SetValue(entry, PositionData{X: 350, Y: 460})
	Velocity.SetValue(entry, VelocityData{X: 5, Y: 0})
	Render.SetValue(entry, RenderData{&SpriteData{image: GetImage("ship")}})
	return nil
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

func (p *PlayerData) Update(w donburi.World, entry *donburi.Entry) error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		p.Move(Right, entry)
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		p.Move(Left, entry)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.Shoot(w, entry)
	}
	return nil
}

func (p *PlayerData) Move(dir Direction, entry *donburi.Entry) {
	pos := Position.Get(entry)
	v := Velocity.Get(entry)
	// TODO check for bounds
	delta := v.X
	if dir == Left {
		delta = -delta
	}
	pos.X += delta
}

func (p *PlayerData) Shoot(w donburi.World, entry *donburi.Entry) {
	pos := Position.Get(entry)
	entity := w.Create(Bullet, Position, Velocity, Render)
	bEntry := w.Entry(entity)
	Position.SetValue(bEntry, PositionData{X: pos.X + 24, Y: pos.Y - 10})
	Velocity.SetValue(bEntry, VelocityData{X: 0, Y: -3})
	Render.SetValue(bEntry, RenderData{&BulletRenderData{Color: color.RGBA{255, 215, 0, 255}}})
	Bullet.SetValue(bEntry, BulletData{Length: 10, Width: 3})
}

func (p *PlayerData) AddScore(score int) {
	p.Score += score
}

func (p *PlayerData) GetScore() int {
	return p.Score
}

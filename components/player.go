package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
)

type PlayerData struct {
	score int
	dead  bool
}

var Player = donburi.NewComponentType[PlayerData]()

func NewPlayer(w donburi.World) error {
	entity := w.Create(Player, Position, Velocity, Render)
	entry := w.Entry(entity)
	Position.SetValue(entry, PositionData{x: 350, y: 460})
	Velocity.SetValue(entry, VelocityData{x: 5, y: 0})
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
	if p.dead {
		return nil
	}
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
	pt := Render.Get(entry).GetRect(entry).Size()
	width := pt.X

	delta := v.x
	if dir == Left {
		delta = -delta
	}
	newX := pos.x + delta

	// check for bounds
	be := Board.MustFirst(entry.World)
	board := Board.Get(be)
	if newX > 0 && newX < board.Width-width {
		pos.x += delta
	}
}

func (p *PlayerData) Shoot(w donburi.World, entry *donburi.Entry) {
	if p.dead {
		return
	}
	pos := Position.Get(entry)
	entity := w.Create(Bullet, Position, Velocity, Render)
	bEntry := w.Entry(entity)
	Position.SetValue(bEntry, PositionData{x: pos.x + 24, y: pos.y - 10})
	Velocity.SetValue(bEntry, VelocityData{x: 0, y: -3})
	Render.SetValue(bEntry, RenderData{&BulletRenderData{color: color.RGBA{255, 215, 0, 255}}})
	Bullet.SetValue(bEntry, BulletData{length: 10, width: 3})
}

func (p *PlayerData) AddScore(score int) {
	p.score += score
}

func (p *PlayerData) GetScore() int {
	return p.score
}

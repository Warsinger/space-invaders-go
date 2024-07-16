package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type PlayerData struct {
}

var Player = donburi.NewComponentType[PlayerData]()

func NewPlayer(w donburi.World) {
	entity := w.Create(Player, Position, Velocity, Sprite)
	entry := w.Entry(entity)
	Position.SetValue(entry, PositionData{X: 350, Y: 460})
	Velocity.SetValue(entry, VelocityData{X: 5, Y: 0})
	Sprite.SetValue(entry, SpriteData{image: GetImage("ship")})
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

func (p *PlayerData) Update(pos *PositionData, v *VelocityData) error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		p.Move(Right, pos, v)
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		p.Move(Left, pos, v)
	}
	return nil
}

func (p *PlayerData) Move(dir Direction, pos *PositionData, v *VelocityData) {
	// TODO check for bounds
	var delta float64 = v.X
	if dir == Left {
		delta = -delta
	}
	pos.X += delta
}

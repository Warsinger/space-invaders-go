package game

import (
	"github.com/yohamta/donburi"
)

type AlienData struct {
}

var Alien = donburi.NewComponentType[AlienData]()

func NewAlien(w donburi.World) {
	entity := w.Create(Alien, Position, Velocity, Sprite)
	entry := w.Entry(entity)
	Position.SetValue(entry, PositionData{X: 100, Y: 100})
	Velocity.SetValue(entry, VelocityData{X: 3, Y: 50})
	Sprite.SetValue(entry, SpriteData{image: GetImage("alien")})
}

func (a *AlienData) Update(pos *PositionData, v *VelocityData) error {
	pos.X += v.X
	if pos.X > 750 {
		pos.X = 50
		pos.Y += v.Y
		if pos.Y > 550 {
			pos.Y = 50
		}
	}
	return nil
}

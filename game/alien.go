package game

import (
	"github.com/yohamta/donburi"
)

type AlienData struct {
}

var Alien = donburi.NewComponentType[AlienData]()

const xOffset = 60
const xBorder = 25
const yOffset = 60
const yBorder = 25

func NewAliens(w donburi.World, rows, columns int) {
	for r := 0; r < rows; r++ {
		for c := 0; c < columns; c++ {
			NewAlien(w, float64(xBorder+c*xOffset), float64(yBorder+r*yOffset))
		}
	}
}

func NewAlien(w donburi.World, x, y float64) {
	entity := w.Create(Alien, Position, Velocity, Sprite)
	entry := w.Entry(entity)
	Position.SetValue(entry, PositionData{X: x, Y: y})
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

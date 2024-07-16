package game

import (
	"github.com/yohamta/donburi"
)

type AlienData struct {
	XStart float64
	XRange float64
	YStart float64
	YRange float64
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
	Velocity.SetValue(entry, VelocityData{X: 1, Y: 10})
	Sprite.SetValue(entry, SpriteData{image: GetImage("alien")})
	Alien.SetValue(entry, AlienData{XStart: x, XRange: 75, YStart: y, YRange: 500})
}

func (a *AlienData) Update(pos *PositionData, v *VelocityData) error {
	pos.X += v.X
	if pos.X > a.XStart+a.XRange || pos.X < a.XStart {
		if pos.X < a.XStart {
			pos.X = a.XStart
		} else {
			pos.X = a.XStart + a.XRange
		}
		v.X = -v.X
		pos.Y += v.Y
		if pos.Y > a.YRange {
			pos.Y = a.YStart // TODO this should kill the player
		}
	}
	return nil
}

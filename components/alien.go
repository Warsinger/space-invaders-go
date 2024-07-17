package components

import (
	"image"

	"github.com/yohamta/donburi"
)

type AlienData struct {
	XStart     int
	XRange     int
	YStart     int
	YRange     int
	ScoreValue int
}

var Alien = donburi.NewComponentType[AlienData]()

const xOffset = 60
const xBorder = 25
const yOffset = 60
const yBorder = 25

func NewAliens(w donburi.World, rows, columns int) error {
	for r := 0; r < rows; r++ {
		for c := 0; c < columns; c++ {
			err := NewAlien(w, xBorder+c*xOffset, yBorder+r*yOffset)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewAlien(w donburi.World, x, y int) error {
	entity := w.Create(Alien, Position, Velocity, Render)
	entry := w.Entry(entity)
	Position.SetValue(entry, PositionData{X: x, Y: y})
	Velocity.SetValue(entry, VelocityData{X: 1, Y: 10})
	Render.SetValue(entry, RenderData{&SpriteData{image: GetImage("alien")}})
	Alien.SetValue(entry, AlienData{XStart: x, XRange: 75, YStart: y, YRange: 500, ScoreValue: 10})
	return nil
}

func (a *AlienData) Update(entry *donburi.Entry) error {
	pos := Position.Get(entry)
	v := Velocity.Get(entry)
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

func (a *AlienData) GetRect(entry *donburi.Entry) image.Rectangle {
	sprite := Render.Get(entry)
	return sprite.Renderer.GetRect(entry)
}

func (a *AlienData) GetScoreValue() int {
	return a.ScoreValue
}

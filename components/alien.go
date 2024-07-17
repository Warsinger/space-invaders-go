package components

import (
	"image"
	"log"

	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type AlienData struct {
	xStart     int
	xRange     int
	yStart     int
	yRange     int
	scoreValue int
}

var Alien = donburi.NewComponentType[AlienData]()

const xOffset = 60
const xBorder = 10
const yOffset = 60
const yBorder = 40

func NewAliens(w donburi.World, rows, columns int) error {
	query := donburi.NewQuery(filter.Contains(Board))
	be, found := query.First(w)
	board := Board.Get(be)
	if !found {
		log.Fatal("No Board found")
	}
	for r := 0; r < rows; r++ {
		for c := 0; c < columns; c++ {
			err := NewAlien(w, board, xBorder+c*xOffset, yBorder+r*yOffset)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewAlien(w donburi.World, b *BoardInfo, x, y int) error {
	entity := w.Create(Alien, Position, Velocity, Render)
	entry := w.Entry(entity)
	Position.SetValue(entry, PositionData{x: x, y: y})
	Velocity.SetValue(entry, VelocityData{x: 1, y: 10})
	Render.SetValue(entry, RenderData{&SpriteData{image: GetImage("alien")}})
	Alien.SetValue(entry, AlienData{xStart: x, xRange: 75, yStart: y, yRange: b.Height, scoreValue: 10})
	return nil
}

func (a *AlienData) Update(entry *donburi.Entry) error {
	pos := Position.Get(entry)
	v := Velocity.Get(entry)

	pos.x += v.x
	if pos.x > a.xStart+a.xRange || pos.x < a.xStart {
		if pos.x < a.xStart {
			pos.x = a.xStart
		} else {
			pos.x = a.xStart + a.xRange
		}
		v.x = -v.x
		pos.y += v.y
		if pos.y > a.yRange {
			pos.y = a.yStart // TODO this should kill the player
		}
	}
	return nil
}

func (a *AlienData) GetRect(entry *donburi.Entry) image.Rectangle {
	sprite := Render.Get(entry)
	return sprite.renderer.GetRect(entry)
}

func (a *AlienData) GetScoreValue() int {
	return a.scoreValue
}

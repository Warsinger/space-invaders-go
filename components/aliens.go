package components

import (
	"fmt"
	"image"
	"log"

	assets "space-invaders/assets"

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

const (
	xOffset       = 60
	xBorder       = 25
	yOffset       = 60
	yBorderTop    = 60
	yBorderBottom = 45
)

func NewAliens(w donburi.World, level, rows, columns int) error {
	query := donburi.NewQuery(filter.Contains(Board))
	be, found := query.First(w)
	board := Board.Get(be)
	if !found {
		log.Fatal("No Board found")
	}
	// past a certain level don't just increase the speed but increase the rows by 1
	if level > 10 {
		rows++
	}
	for r := 0; r < rows; r++ {
		for c := 0; c < columns; c++ {
			choose := (r+c)%2 + 1
			err := NewAlien(w, board, level, xBorder+c*xOffset, yBorderTop+r*yOffset, choose)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func NewAlien(w donburi.World, b *BoardInfo, level, x, y, choose int) error {
	entity := w.Create(Alien, Position, Velocity, Render)
	entry := w.Entry(entity)
	Position.SetValue(entry, PositionData{x: x, y: y})
	Velocity.SetValue(entry, VelocityData{x: level/3 + 1, y: 10})
	name := fmt.Sprintf("alien%v", choose)
	Render.SetValue(entry, RenderData{&SpriteData{image: assets.GetImage(name)}})
	Alien.SetValue(entry, AlienData{xStart: x, xRange: 75, yStart: y, yRange: b.Height, scoreValue: 10 + level*5})
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

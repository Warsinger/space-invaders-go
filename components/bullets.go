package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
)

type BulletData struct {
}

var Bullet = donburi.NewComponentType[BulletData]()

type BulletRenderData struct {
	Length, Width float64
	Color         color.Color
}

// var BulletRender = donburi.NewComponentType[BulletRenderData]()

func (brd *BulletRenderData) Draw(screen *ebiten.Image, pos *PositionData) {
	vector.StrokeLine(screen, float32(pos.X), float32(pos.Y), float32(pos.X), float32(pos.Y+brd.Length), float32(brd.Width), brd.Color, true)
}

func (brd *BulletData) Update(pos *PositionData, v *VelocityData) error {
	newY := pos.Y + v.Y
	if newY < 0 {
		pos.Y = 0 // TODO remove from the world
	} else {
		pos.Y = newY
	}
	return nil
}

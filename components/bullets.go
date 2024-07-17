package components

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
)

type BulletData struct {
	Length, Width int
}

var Bullet = donburi.NewComponentType[BulletData]()

type BulletRenderData struct {
	Color color.Color
}

func (brd *BulletRenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	b := Bullet.Get(entry)
	pos := Position.Get(entry)
	vector.StrokeLine(screen, float32(pos.X), float32(pos.Y), float32(pos.X), float32(pos.Y+b.Length), float32(b.Width), brd.Color, true)
}

func (bd *BulletData) Update(entry *donburi.Entry) error {
	pos := Position.Get(entry)
	v := Velocity.Get(entry)
	newY := pos.Y + v.Y
	if newY < 0 {
		entry.World.Remove(entry.Entity())
	} else {
		pos.Y = newY
	}
	return nil
}

func (brd *BulletRenderData) GetRect(entry *donburi.Entry) image.Rectangle {
	pos := Position.Get(entry)
	b := Bullet.Get(entry)
	return image.Rect(pos.X, pos.Y, pos.X+b.Width, pos.Y+b.Length)
}

package components

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

	be := Board.MustFirst(entry.World)
	board := Board.Get(be)

	Position.SetValue(entry, PositionData{x: board.Width / 2, y: board.Height - yBorder})
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

func (p *PlayerData) IsDead() bool {
	return p.dead
}

func (p *PlayerData) Kill() {
	p.dead = true
}

func (p *PlayerData) GetRect(entry *donburi.Entry) image.Rectangle {
	sprite := Render.Get(entry)
	return sprite.renderer.GetRect(entry)
}

func (p *PlayerData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	if p.dead {
		rect := p.GetRect(entry)
		vector.StrokeLine(screen, float32(rect.Min.X), float32(rect.Min.Y), float32(rect.Max.X), float32(rect.Max.Y), 3, color.RGBA{255, 0, 0, 255}, true)
		vector.StrokeLine(screen, float32(rect.Max.X), float32(rect.Min.Y), float32(rect.Min.X), float32(rect.Max.Y), 3, color.RGBA{255, 0, 0, 255}, true)
	}
	// draw score
	str := fmt.Sprintf("SCORE %04d", p.score)
	op := &text.DrawOptions{}
	// op.LineSpacing = scoreFace.Size * 1.5
	x, _ := text.Measure(str, ScoreFace, op.LineSpacing)
	op.GeoM.Translate(400-x/2, 5)
	text.Draw(screen, str, ScoreFace, op)

}

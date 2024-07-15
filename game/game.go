package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type GameInfo struct {
	World donburi.World
}
type Game interface {
	GetWorld() donburi.World

	Update() error
	Draw(screen *ebiten.Image)
	Layout(width, height int) (int, int)
}

func NewGame() Game {
	return &GameInfo{
		donburi.NewWorld(),
	}
}

func (g *GameInfo) GetWorld() donburi.World {
	return g.World
}

func (g *GameInfo) Update() error {
	return nil
}

func (g *GameInfo) Draw(screen *ebiten.Image) {
	screen.Clear()
}

func (g *GameInfo) Layout(width, height int) (int, int) {
	return width, height
}

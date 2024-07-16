package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// Component is any struct that holds some kind of data.
type PositionData struct {
	X, Y float64
}

type VelocityData struct {
	X, Y float64
}

// ComponentType represents kind of component which is used to create or query entities.
var Position = donburi.NewComponentType[PositionData]()
var Velocity = donburi.NewComponentType[VelocityData]()

type Renderer interface {
	Draw(screen *ebiten.Image, pos *PositionData)
}

type RenderData struct {
	Renderer Renderer
}

var Render = donburi.NewComponentType[RenderData]()

func (r *RenderData) Draw(screen *ebiten.Image, pos *PositionData) {
	r.Renderer.Draw(screen, pos)
}

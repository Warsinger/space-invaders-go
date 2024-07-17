package components

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// Component is any struct that holds some kind of data.
type PositionData struct {
	X, Y int
}

type VelocityData struct {
	X, Y int
}

// ComponentType represents kind of component which is used to create or query entities.
var Position = donburi.NewComponentType[PositionData]()
var Velocity = donburi.NewComponentType[VelocityData]()

type Renderer interface {
	Draw(screen *ebiten.Image, entry *donburi.Entry)
	GetRect(entry *donburi.Entry) image.Rectangle
}

type RenderData struct {
	Renderer Renderer
}

var Render = donburi.NewComponentType[RenderData]()

func (r *RenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	r.Renderer.Draw(screen, entry)
}

func (r *RenderData) GetRect(entry *donburi.Entry) image.Rectangle {
	return r.Renderer.GetRect(entry)
}

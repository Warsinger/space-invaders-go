package components

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/yohamta/donburi"
)

type BarrierData struct {
	lines  int
	length int
}

type BarrierRenderData struct {
	color     color.Color
	lineWidth int
}

var Barrier = donburi.NewComponentType[BarrierData]()

const maxBarriers = 6

func NewBarriers(w donburi.World, level int) error {
	count := min(level, maxBarriers)
	fmt.Printf("barrier count %v\n", count)
	be := Board.MustFirst(w)
	board := Board.Get(be)

	const lineLength = 50
	const linesPerBarrier = 4
	for i := 0; i < count; i++ {
		spacing := board.Width / count
		blank := (spacing - lineLength) / 2
		x := spacing*i + blank
		y := board.Height - yBorderBottom*2
		NewBarrier(w, linesPerBarrier, x, y)
	}

	return nil
}

func NewBarrier(w donburi.World, lines, x, y int) error {
	entity := w.Create(Barrier, Position, Velocity, Render)
	entry := w.Entry(entity)

	Position.SetValue(entry, PositionData{x: x, y: y})
	Render.SetValue(entry, RenderData{&BarrierRenderData{color: color.RGBA{100, 215, 100, 255}, lineWidth: 4}})
	Barrier.SetValue(entry, BarrierData{lines: lines, length: 50})
	return nil
}

func (b *BarrierData) ProcessHit(entry *donburi.Entry) {
	if b.lines > 0 {
		b.lines--
		if b.lines <= 0 {
			entry.Remove()
		}
	}
}

func (brd *BarrierRenderData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	b := Barrier.Get(entry)
	pos := Position.Get(entry)
	startX, startY := pos.x, pos.y
	for i := 0; i < b.lines; i++ {
		vector.StrokeLine(screen, float32(startX), float32(startY), float32(startX+b.length), float32(startY), float32(brd.lineWidth), brd.color, true)
		startY += brd.lineWidth
	}
}

func (brd *BarrierRenderData) GetRect(entry *donburi.Entry) image.Rectangle {
	pos := Position.Get(entry)
	b := Barrier.Get(entry)
	height := brd.lineWidth * b.lines
	return image.Rect(pos.x, pos.y, pos.x+b.length, pos.y+height)
}

func (b *BarrierData) GetRect(entry *donburi.Entry) image.Rectangle {
	r := Render.Get(entry)
	return r.renderer.GetRect(entry)
}

package game

import "github.com/yohamta/donburi"

type BoardInfo struct {
	Width, Height float64
}

var Board = donburi.NewComponentType[BoardInfo]()

func NewBoard(w donburi.World) (BoardInfo, error) {
	entity := w.Create(Board)
	entry := w.Entry(entity)
	b := BoardInfo{Width: 800, Height: 600}
	Board.SetValue(entry, b)
	return b, nil
}

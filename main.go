package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	g := &GameInfo{}

	g.Init()

	mX, mY := ebiten.Monitor().Size()
	xScale := float64(g.Board.GetWidth()) / float64(mX)
	yScale := float64(g.Board.GetHeight()) / float64(mY)
	scale := math.Max(xScale, yScale) * 1.1

	ebiten.SetWindowSize(int(float64(g.Board.GetWidth())/scale), int(float64(g.Board.GetHeight())/scale))
	ebiten.SetWindowTitle("Space Invaders")
	// ebiten.SetWindowDecorated(false)
	// ebiten.SetTPS(15)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

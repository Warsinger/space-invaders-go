package main

import (
	"log"
	"space-invaders/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	g := game.NewGame()
	g.Init()

	// mX, mY := ebiten.Monitor().Size()

	// ebiten.SetWindowSize(mX-50, mY-50)
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Space Invaders")
	// ebiten.SetWindowDecorated(false)
	// ebiten.SetTPS(15)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

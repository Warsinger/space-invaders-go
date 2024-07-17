package main

import (
	"log"
	"space-invaders/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	g, err := game.NewGame()
	if err != nil {
		log.Fatal(err)
	}
	err = g.Init()
	if err != nil {
		log.Fatal(err)
	}
	// mX, mY := ebiten.Monitor().Size()

	// ebiten.SetWindowSize(mX-50, mY-50)
	// ebiten.SetWindowDecorated(false)
	// ebiten.SetTPS(15)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

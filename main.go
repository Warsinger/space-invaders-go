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

	// ebiten.SetTPS(60)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

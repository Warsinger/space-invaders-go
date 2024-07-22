package main

import (
	"flag"
	"log"
	"space-invaders/game"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	level := flag.Int("level", 1, "Starting level")
	flag.Parse()

	g, err := game.NewGame(*level)
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

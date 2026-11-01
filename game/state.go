package game

import (
	"fmt"

	"space-invaders/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

// menuInput separates state transitions from keyboard/pointer sampling and simulation.
type menuInput struct{ start, reset, title, pause, quit bool }

func (g *GameInfo) saveHighScore() {
	if err := g.SaveScores(); err != nil {
		fmt.Printf("WARN save high score: %v\n", err)
	}
}

func (g *GameInfo) startRun() error {
	g.saveHighScore()
	if err := g.Clear(); err != nil {
		return err
	}
	if err := g.Init(); err != nil {
		return err
	}
	g.onTitle = false
	return nil
}

func (g *GameInfo) returnToTitle() error {
	g.saveHighScore()
	assets.StopSounds()
	if err := g.Clear(); err != nil {
		return err
	}
	g.onTitle = true
	return nil
}

// A transition consumes this tick so Space used to start cannot also fire a shot.
func (g *GameInfo) handleMenuInput(input menuInput) (bool, error) {
	if input.quit {
		if g.onTitle {
			g.saveHighScore()
		} else {
			g.EndGame()
		}
		return true, ebiten.Termination
	}
	if g.onTitle {
		if input.start {
			return true, g.startRun()
		}
		return true, nil
	}
	if input.title || input.reset {
		return true, g.returnToTitle()
	}
	if g.gameOver {
		return true, nil
	}
	if input.pause {
		g.paused = !g.paused
	}
	return g.paused, nil
}

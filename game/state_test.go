package game

import (
	"errors"
	"image"
	"reflect"
	"testing"

	comp "space-invaders/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

func titleFixture(t *testing.T) *GameInfo {
	t.Helper()
	t.Chdir(t.TempDir())
	world := donburi.NewWorld()
	if _, err := comp.NewBoard(world); err != nil {
		t.Fatal(err)
	}
	return &GameInfo{world: world, onTitle: true, startingLevel: 3, level: 3, highScore: 42}
}
func countComponent(g *GameInfo, component donburi.IComponentType) int {
	return donburi.NewQuery(filter.Contains(component)).Count(g.world)
}

func TestTitleWaitsWithoutSimulation(t *testing.T) {
	g := titleFixture(t)
	board := comp.Board.Get(comp.Board.MustFirst(g.world))
	if err := comp.NewAlien(g.world, board, 1, 25, 60, 1); err != nil {
		t.Fatal(err)
	}
	alien := comp.Alien.MustFirst(g.world)
	before := *comp.Position.Get(alien)
	for i := 0; i < 600; i++ {
		if err := g.updateWithInput(menuInput{pause: i%2 == 0, reset: i%3 == 0}); err != nil {
			t.Fatal(err)
		}
	}
	if !g.onTitle || g.gameOver || g.paused || g.score != 0 || countComponent(g, comp.Bullet) != 0 || !reflect.DeepEqual(before, *comp.Position.Get(alien)) {
		t.Fatal("title input/ticks advanced simulation or started a run")
	}
}

func TestStartPauseResetAndReturnToTitle(t *testing.T) {
	g := titleFixture(t)
	if err := g.updateWithInput(menuInput{start: true}); err != nil {
		t.Fatal(err)
	}
	if g.onTitle || g.gameOver || g.paused || countComponent(g, comp.Player) != 1 || countComponent(g, comp.Alien) != 32 || countComponent(g, comp.Bullet) != 0 {
		t.Fatal("start did not create exactly one fresh run or also fired a shot")
	}
	stopped, err := g.handleMenuInput(menuInput{pause: true})
	if err != nil || !stopped || !g.paused {
		t.Fatal("pause did not stop updates")
	}
	stopped, err = g.handleMenuInput(menuInput{pause: true})
	if err != nil || stopped || g.paused {
		t.Fatal("resume did not restore play")
	}
	g.AddScore(75)
	g.level = 8
	g.paused = true
	if _, err := g.handleMenuInput(menuInput{reset: true}); err != nil {
		t.Fatal(err)
	}
	if !g.onTitle || g.paused || g.score != 0 || g.level != 3 || g.highScore != 75 || countComponent(g, comp.Player) != 0 || countComponent(g, comp.Alien) != 0 {
		t.Fatal("R did not return to title and clear the run while preserving high score")
	}
	if _, err := g.handleMenuInput(menuInput{start: true}); err != nil {
		t.Fatal(err)
	}
	if g.onTitle || countComponent(g, comp.Player) != 1 || countComponent(g, comp.Alien) != 32 {
		t.Fatal("new start accumulated or omitted run entities")
	}
	g.gameOver = true
	if _, err := g.handleMenuInput(menuInput{start: true, pause: true}); err != nil {
		t.Fatal(err)
	}
	if !g.gameOver {
		t.Fatal("Space/P escaped game over without returning to title")
	}
	if _, err := g.handleMenuInput(menuInput{title: true}); err != nil {
		t.Fatal(err)
	}
	if !g.onTitle || g.gameOver || g.paused || countComponent(g, comp.Player) != 0 || countComponent(g, comp.Alien) != 0 || g.highScore != 75 || LoadScores() != 75 {
		t.Fatal("return to title failed to stop/clear run or retain saved high score")
	}
	if _, err := g.handleMenuInput(menuInput{start: true}); err != nil {
		t.Fatal(err)
	}
	if g.onTitle || g.score != 0 || g.level != 3 || g.highScore != 75 {
		t.Fatal("title could not start another fresh run")
	}
}

func TestNativeQuitFromTitleSavesWithoutStarting(t *testing.T) {
	g := titleFixture(t)
	stopped, err := g.handleMenuInput(menuInput{quit: true, start: true})
	if !stopped || !errors.Is(err, ebiten.Termination) || !g.onTitle || countComponent(g, comp.Player) != 0 || LoadScores() != 42 {
		t.Fatal("title quit did not save/terminate without starting")
	}
}

func TestStartTargetFitsLogicalBoardAndHasExactEdges(t *testing.T) {
	bounds := startButtonBounds(600, 800)
	if !bounds.In(image.Rect(0, 0, 600, 800)) || bounds.Dy() < 88 {
		t.Fatal("start target clipped or too small")
	}
	for _, point := range []image.Point{bounds.Min, {300, 700}, {bounds.Max.X - 1, bounds.Max.Y - 1}} {
		if !point.In(bounds) {
			t.Fatal("button rejected inside point")
		}
	}
	for _, point := range []image.Point{{bounds.Min.X - 1, 700}, {bounds.Max.X, 700}, {300, bounds.Max.Y}} {
		if point.In(bounds) {
			t.Fatal("button accepted outside point")
		}
	}
}

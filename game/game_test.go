package game

import (
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
	comp "space-invaders/components"
)

type fixedRenderer struct{ rect image.Rectangle }

func (fixedRenderer) Draw(*ebiten.Image, *donburi.Entry)       {}
func (r fixedRenderer) GetRect(*donburi.Entry) image.Rectangle { return r.rect }

func TestScoresAndClear(t *testing.T) {
	t.Chdir(t.TempDir())
	g := &GameInfo{world: donburi.NewWorld(), startingLevel: 3, level: 5}
	g.AddScore(20)
	g.AddScore(15)
	if g.GetScore() != 35 || LoadScores() != 35 {
		t.Fatal("high score must save as score improves")
	}
	g.gameOver = true
	g.paused = true
	if err := g.Clear(); err != nil {
		t.Fatal(err)
	}
	if g.score != 0 || g.highScore != 35 || g.level != 3 || g.gameOver || g.paused {
		t.Fatal("clear lost high score or failed to reset run")
	}
	g.AddScore(5)
	if LoadScores() != 35 {
		t.Fatal("lower new run must not replace high score")
	}
}

func TestBulletConsumedByOneTarget(t *testing.T) {
	w := donburi.NewWorld()
	g := &GameInfo{world: w}
	for i := 0; i < 2; i++ {
		e := w.Entry(w.Create(comp.Alien, comp.Render))
		comp.Render.SetValue(e, comp.RenderData{})
		*comp.Render.Get(e) = renderAt(image.Rect(0, 0, 20, 20))
	}
	b := w.Entry(w.Create(comp.Bullet, comp.Render))
	*comp.Render.Get(b) = renderAt(image.Rect(0, 0, 4, 16))
	if err := g.DetectCollisions(); err != nil {
		t.Fatal(err)
	}
	if donburi.NewQuery(filter.Contains(comp.Alien)).Count(w) != 1 || donburi.NewQuery(filter.Contains(comp.Bullet)).Count(w) != 0 {
		t.Fatal("one bullet must remove exactly one alien")
	}
}

func TestAlienBulletHitsPlayerWithoutOtherTargets(t *testing.T) {
	t.Chdir(t.TempDir())
	w := donburi.NewWorld()
	g := &GameInfo{world: w}
	p := w.Entry(w.Create(comp.Player, comp.Render))
	*comp.Render.Get(p) = renderAt(image.Rect(0, 0, 20, 20))
	shooter := w.Entry(w.Create(comp.Alien, comp.Position))
	comp.Alien.Get(shooter).Shoot(shooter)
	shooter.Remove()
	b := comp.Bullet.MustFirst(w)
	*comp.Render.Get(b) = renderAt(image.Rect(0, 0, 4, 16))
	if err := g.DetectCollisions(); err != nil {
		t.Fatal(err)
	}
	if !g.gameOver || !comp.Player.Get(p).IsDead() || donburi.NewQuery(filter.Contains(comp.Bullet)).Count(w) != 0 {
		t.Fatal("alien bullet must kill player even with no aliens/barriers left")
	}
}

func renderAt(rect image.Rectangle) comp.RenderData { return comp.NewRenderData(fixedRenderer{rect}) }

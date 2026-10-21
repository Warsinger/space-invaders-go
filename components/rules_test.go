package components

import (
	"image"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type boundsRenderer struct{}

func (boundsRenderer) Draw(*ebiten.Image, *donburi.Entry) {}
func (boundsRenderer) GetRect(e *donburi.Entry) image.Rectangle {
	p := Position.Get(e)
	return image.Rect(p.x, p.y, p.x+48, p.y+48)
}

func TestPlayerBounds(t *testing.T) {
	w := donburi.NewWorld()
	NewBoard(w)
	e := w.Entry(w.Create(Player, Position, Velocity, Render))
	Velocity.SetValue(e, VelocityData{x: 5})
	Render.SetValue(e, RenderData{boundsRenderer{}})
	p := Player.Get(e)
	for _, tc := range []struct {
		x    int
		dir  Direction
		want int
	}{{1, Left, 1}, {10, Left, 5}, {550, Right, 550}, {100, Right, 105}} {
		Position.SetValue(e, PositionData{x: tc.x})
		p.Move(tc.dir, e)
		if got := Position.Get(e).x; got != tc.want {
			t.Errorf("move from %d=%d want %d", tc.x, got, tc.want)
		}
	}
	p.Kill()
	if !p.IsDead() {
		t.Fatal("player must be dead")
	}
}

func TestBulletExpiry(t *testing.T) {
	for _, tc := range []struct {
		y, v   int
		remain bool
	}{{1, -3, false}, {799, 3, false}, {3, -3, true}, {400, 3, true}} {
		w := donburi.NewWorld()
		NewBoard(w)
		e := w.Entry(w.Create(Bullet, Position, Velocity))
		Position.SetValue(e, PositionData{y: tc.y})
		Velocity.SetValue(e, VelocityData{y: tc.v})
		if err := Bullet.Get(e).Update(e); err != nil {
			t.Fatal(err)
		}
		if got := donburi.NewQuery(filter.Contains(Bullet)).Count(w) > 0; got != tc.remain {
			t.Errorf("bullet(%d,%d) remains=%v", tc.y, tc.v, got)
		}
	}
}

func TestBarriers(t *testing.T) {
	for _, tc := range []struct{ level, count int }{{1, 2}, {2, 2}, {3, 3}, {6, 6}, {7, 2}, {12, 6}} {
		w := donburi.NewWorld()
		NewBoard(w)
		if err := NewBarriers(w, tc.level); err != nil {
			t.Fatal(err)
		}
		if donburi.NewQuery(filter.Contains(Barrier)).Count(w) != tc.count {
			t.Fatalf("level %d barrier count", tc.level)
		}
		e := Barrier.MustFirst(w)
		b := Barrier.Get(e)
		for i := 0; i < 3; i++ {
			b.ProcessHit(e)
		}
		if b.lines != 1 {
			t.Fatal("three hits must leave one line")
		}
		b.ProcessHit(e)
		if donburi.NewQuery(filter.Contains(Barrier)).Count(w) != tc.count-1 {
			t.Fatal("fourth hit must remove barrier")
		}
	}
}

func TestAlienReversesAndDescends(t *testing.T) {
	w := donburi.NewWorld()
	e := w.Entry(w.Create(Alien, Position, Velocity))
	a := Alien.Get(e)
	a.xStart = 25
	a.xRange = 75
	Position.SetValue(e, PositionData{x: 100, y: 60})
	Velocity.SetValue(e, VelocityData{x: 2, y: 10})
	if err := a.Update(e); err != nil {
		t.Fatal(err)
	}
	if p, v := Position.Get(e), Velocity.Get(e); p.x != 100 || p.y != 70 || v.x != -2 {
		t.Fatal("alien did not reverse at right limit")
	}
}

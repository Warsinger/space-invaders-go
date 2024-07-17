package game

import (
	"fmt"
	comp "space-invaders/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	ecslib "github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type GameInfo struct {
	world    donburi.World
	ecs      *ecslib.ECS
	gameOver bool
}

func NewGame() (*GameInfo, error) {
	world := donburi.NewWorld()
	ecs := ecslib.NewECS(world)
	board, err := comp.NewBoard(world)
	if err != nil {
		return nil, err
	}
	err = comp.LoadAssets()
	if err != nil {
		return nil, err
	}
	ebiten.SetWindowSize(int(board.Width), int(board.Height))
	ebiten.SetWindowTitle("Space Invaders")

	return &GameInfo{
		world,
		ecs,
		false,
	}, nil
}

func (g *GameInfo) Init() error {
	err := comp.NewPlayer(g.world)
	if err != nil {
		return err
	}
	err = comp.NewAliens(g.world, 4, 10)
	if err != nil {
		return err
	}
	return nil
}

func (g *GameInfo) Clear() error {
	g.gameOver = false
	query := donburi.NewQuery(filter.Or(
		filter.Contains(comp.Bullet),
		filter.Contains(comp.Player),
		filter.Contains(comp.Alien),
	))
	query.Each(g.world, func(e *donburi.Entry) {
		g.world.Remove(e.Entity())
	})
	return nil
}

func (g *GameInfo) GetWorld() donburi.World {
	return g.world
}
func (g *GameInfo) GetECS() *ecslib.ECS {
	return g.ecs
}

func (g *GameInfo) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Clear()
		g.Init()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	if g.gameOver {
		return nil
	}

	// query for all entities that have position and velocity and ???
	// and have them do their updates
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(comp.Position, comp.Velocity, comp.Render),
			filter.Or(
				filter.Contains(comp.Player),
				filter.Contains(comp.Alien),
				filter.Contains(comp.Bullet),
			),
		),
	)
	var err error = nil
	// update all entities
	query.Each(g.world, func(entry *donburi.Entry) {
		if entry.HasComponent(comp.Player) {
			player := comp.Player.Get(entry)
			err = player.Update(g.world, entry)
			if err != nil {
				return
			}
		}

		if entry.HasComponent(comp.Alien) {
			alien := comp.Alien.Get(entry)
			err = alien.Update(entry)
			if err != nil {
				return
			}

		}

		if entry.HasComponent(comp.Bullet) {
			b := comp.Bullet.Get(entry)
			err = b.Update(entry)
			if err != nil {
				return
			}
		}
	})
	// after updating all positions check for collisions
	// get all the bullets, for each bullet loop through all the aliens (or other objects) and see if there are collisions
	// if there is a collition, remove both objects (or subtract from their health)
	// accumultate points for killing aliens
	err = g.DetectCollisions()
	if err != nil {
		return err
	}
	return err
}

func (g *GameInfo) DetectCollisions() error {
	var err error = nil
	query := donburi.NewQuery(filter.Contains(comp.Bullet))
	query.Each(g.world, func(be *donburi.Entry) {
		brd := comp.Render.Get(be)
		bRect := brd.GetRect(be)

		query := donburi.NewQuery(filter.Contains(comp.Alien))
		query.Each(g.world, func(ae *donburi.Entry) {
			alien := comp.Alien.Get(ae)
			aRect := alien.GetRect(ae)
			if bRect.Overlaps(aRect) {
				// increment score
				pe := comp.Player.MustFirst(g.world)
				player := comp.Player.Get(pe)
				player.AddScore(alien.GetScoreValue())
				fmt.Printf("Player score: %v\n", player.GetScore())

				// remove bullet and alien
				g.world.Remove(ae.Entity())
				g.world.Remove(be.Entity())

			}

		})
	})

	// check for aliens reaching bottem
	pe := comp.Player.MustFirst(g.world)
	player := comp.Player.Get(pe)
	pRect := player.GetRect(pe)
	query = donburi.NewQuery(filter.Contains(comp.Alien))
	query.Each(g.world, func(ae *donburi.Entry) {
		alien := comp.Alien.Get(ae)
		aRect := alien.GetRect(ae)

		if aRect.Max.Y >= pRect.Min.Y {
			player.Kill()
			g.gameOver = true
		}
	})

	return err
}

func (g *GameInfo) Draw(screen *ebiten.Image) {
	screen.Clear()

	// query for all entities
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(comp.Position, comp.Render),
		),
	)

	// draw all entities
	query.Each(g.world, func(entry *donburi.Entry) {
		r := comp.Render.Get(entry)
		r.Draw(screen, entry)
		if entry.HasComponent(comp.Player) {
			player := comp.Player.Get(entry)
			if player.IsDead() {
				player.DrawDead(screen, entry)
			}
		}
	})
}

func (g *GameInfo) Layout(width, height int) (int, int) {
	return width, height
}

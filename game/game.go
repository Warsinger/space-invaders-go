package game

import (
	comp "space-invaders/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yohamta/donburi"
	ecslib "github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type GameInfo struct {
	world donburi.World
	ecs   *ecslib.ECS
}
type Game interface {
	GetWorld() donburi.World
	GetECS() *ecslib.ECS

	Update() error
	Draw(screen *ebiten.Image)
	Layout(width, height int) (int, int)
	Init() error
}

func NewGame() Game {
	world := donburi.NewWorld()
	ecs := ecslib.NewECS(world)
	return &GameInfo{
		world,
		ecs,
	}
}

func (g *GameInfo) GetWorld() donburi.World {
	return g.world
}
func (g *GameInfo) GetECS() *ecslib.ECS {
	return g.ecs
}

func (g *GameInfo) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Init()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
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
		position := comp.Position.Get(entry)
		velocity := comp.Velocity.Get(entry)

		if entry.HasComponent(comp.Player) {
			player := comp.Player.Get(entry)
			err = player.Update(g.world, position, velocity)
			if err != nil {
				return
			}
		}

		if entry.HasComponent(comp.Alien) {
			alien := comp.Alien.Get(entry)
			err = alien.Update(position, velocity)
			if err != nil {
				return
			}

		}

		if entry.HasComponent(comp.Bullet) {
			b := comp.Bullet.Get(entry)
			err = b.Update(position, velocity)
			if err != nil {
				return
			}
		}
	})

	return err
}

func (g *GameInfo) Init() error {
	comp.LoadAssets()
	comp.NewPlayer(g.world)
	comp.NewAliens(g.world, 4, 10)
	return nil
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
		position := comp.Position.Get(entry)
		r := comp.Render.Get(entry)
		r.Draw(screen, position)

	})
}

func (g *GameInfo) Layout(width, height int) (int, int) {
	return width, height
}

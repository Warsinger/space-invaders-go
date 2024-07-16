package game

import (
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
	// query for all entities
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(Position, Velocity),
			filter.Or(
				filter.Contains(Player),
				filter.Contains(Alien),
			),
		),
	)
	var err error = nil
	// update all entities
	query.Each(g.world, func(entry *donburi.Entry) {
		position := Position.Get(entry)
		velocity := Velocity.Get(entry)

		if entry.HasComponent(Player) {
			player := Player.Get(entry)
			err = player.Update(position, velocity)
			if err != nil {
				return
			}
		}
		if entry.HasComponent(Alien) {
			alien := Alien.Get(entry)
			err = alien.Update(position, velocity)
			if err != nil {
				return
			}
		}
	})

	return err
}

func (g *GameInfo) Init() error {
	LoadAssets()
	NewPlayer(g.world)
	NewAlien(g.world)
	return nil
}

func (g *GameInfo) Draw(screen *ebiten.Image) {
	screen.Clear()

	// query for all entities
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(Position, Sprite),
		),
	)

	// draw all entities
	query.Each(g.world, func(entry *donburi.Entry) {
		position := Position.Get(entry)
		sprite := Sprite.Get(entry)

		sprite.Draw(screen, position)
	})

}

func (g *GameInfo) Layout(width, height int) (int, int) {
	return width, height
}

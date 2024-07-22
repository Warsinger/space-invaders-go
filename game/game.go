package game

import (
	"fmt"
	"os"
	assets "space-invaders/assets"
	comp "space-invaders/components"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
	ecslib "github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type GameInfo struct {
	world     donburi.World
	ecs       *ecslib.ECS
	gameOver  bool
	paused    bool
	level     int
	score     int
	highScore int
}

func NewGame() (*GameInfo, error) {
	world := donburi.NewWorld()
	ecs := ecslib.NewECS(world)
	board, err := comp.NewBoard(world)
	if err != nil {
		return nil, err
	}
	err = assets.LoadAssets()
	if err != nil {
		return nil, err
	}

	highScore := LoadScores()

	ebiten.SetWindowSize(int(board.Width), int(board.Height))
	ebiten.SetWindowTitle("Space Invaders")

	return &GameInfo{
		world:     world,
		ecs:       ecs,
		level:     1,
		highScore: highScore,
	}, nil
}

const highScoreFile = "score/highscore.txt"

func LoadScores() int {
	var highScore int = 0
	bytes, err := os.ReadFile(highScoreFile)
	if err == nil {
		highScore, err = strconv.Atoi(string(bytes))
		if err != nil {
			fmt.Printf("WARN high score formatting err %v\n", err)
		}
	}

	return highScore
}

func (g *GameInfo) SaveScores() error {
	str := strconv.Itoa(g.highScore)

	err := os.WriteFile(highScoreFile, []byte(str), 0644)
	return err
}

func (g *GameInfo) Init() error {
	err := comp.NewPlayer(g.world)
	if err != nil {
		return err
	}
	err = comp.NewAliens(g.world, g.level, 4, 12)
	if err != nil {
		return err
	}

	err = comp.NewBarriers(g.world, g.level)
	if err != nil {
		return err
	}
	return nil
}

func (g *GameInfo) Clear() error {
	g.gameOver = false
	g.paused = false
	g.score = 0
	g.level = 1
	query := donburi.NewQuery(filter.Or(
		filter.Contains(comp.Bullet),
		filter.Contains(comp.Barrier),
		filter.Contains(comp.Player),
		filter.Contains(comp.Alien),
	))
	query.Each(g.world, func(e *donburi.Entry) {
		e.Remove()
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
		g.EndGame()
		return ebiten.Termination
	}

	if g.gameOver {
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.paused = !g.paused
	}
	if g.paused {
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

	// check for all aliens destroyed or aliens reaching bottem
	pe := comp.Player.MustFirst(g.world)
	player := comp.Player.Get(pe)
	pRect := player.GetRect(pe)
	query = donburi.NewQuery(filter.Contains(comp.Alien))
	if query.Count(g.world) == 0 {
		err := g.NewLevel()
		if err != nil {
			return err
		}
	} else {
		query.Each(g.world, func(ae *donburi.Entry) {
			alien := comp.Alien.Get(ae)
			aRect := alien.GetRect(ae)

			if aRect.Max.Y >= pRect.Min.Y {
				player.Kill()
				g.EndGame()
			}
		})
	}

	return err
}
func (g *GameInfo) EndGame() {
	assets.PlaySound("killed")
	g.gameOver = true
	g.SaveScores()
}

func (g *GameInfo) NewLevel() error {
	g.level++
	g.CleanBoard()
	err := comp.NewAliens(g.world, g.level, 4, 12)
	if err != nil {
		return err
	}

	err = comp.NewBarriers(g.world, g.level)
	if err != nil {
		return err
	}

	return nil
}

func (g *GameInfo) CleanBoard() {
	query := donburi.NewQuery(filter.Or(
		filter.Contains(comp.Bullet),
		filter.Contains(comp.Barrier),
		filter.Contains(comp.Alien),
	))

	query.Each(g.world, func(be *donburi.Entry) {
		be.Remove()
	})
}

func (g *GameInfo) DetectCollisions() error {
	var err error = nil
	query := donburi.NewQuery(filter.Contains(comp.Bullet))
	query.Each(g.world, func(bulletEntry *donburi.Entry) {
		bulletRender := comp.Render.Get(bulletEntry)
		bulletRect := bulletRender.GetRect(bulletEntry)

		query := donburi.NewQuery(filter.Or(
			filter.Contains(comp.Alien),
			filter.Contains(comp.Barrier),
		))
		query.Each(g.world, func(e *donburi.Entry) {
			if e.HasComponent(comp.Alien) {
				alien := comp.Alien.Get(e)
				alienRect := alien.GetRect(e)
				if bulletRect.Overlaps(alienRect) {
					g.AddScore(alien.GetScoreValue())

					// remove bullet and alien
					e.Remove()
					bulletEntry.Remove()
					assets.PlaySound("explosion")
				}
			} else if e.HasComponent(comp.Barrier) {
				barrier := comp.Barrier.Get(e)
				barrierRect := barrier.GetRect(e)
				if bulletRect.Overlaps(barrierRect) {
					barrier.ProcessHit(e)
					bulletEntry.Remove()
				}
			}

		})
	})

	return err
}

func (g *GameInfo) Draw(screen *ebiten.Image) {
	screen.Clear()

	img := assets.GetImage("background")
	opts := &ebiten.DrawImageOptions{}
	screen.DrawImage(img, opts)

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
			player.Draw(screen, entry)

		}
	})

	// draw level
	str := fmt.Sprintf("LEVEL %02d", g.level)
	op := &text.DrawOptions{}
	op.GeoM.Translate(5, 5)
	text.Draw(screen, str, assets.ScoreFace, op)

	// draw score
	str = fmt.Sprintf("SCORE %05d", g.score)
	op = &text.DrawOptions{}
	x, y := text.Measure(str, assets.ScoreFace, op.LineSpacing)
	op.GeoM.Translate(400-x/2, 5+y)
	text.Draw(screen, str, assets.ScoreFace, op)

	// draw high score
	str = fmt.Sprintf("HIGH %05d", g.highScore)
	op = &text.DrawOptions{}
	x, _ = text.Measure(str, assets.ScoreFace, op.LineSpacing)
	op.GeoM.Translate(795-x, 5)
	text.Draw(screen, str, assets.ScoreFace, op)

	if g.gameOver {
		// draw game over
		str := "GAME OVER"
		op := &text.DrawOptions{}
		x, y := text.Measure(str, assets.ScoreFace, op.LineSpacing)
		op.GeoM.Translate(400-x/2, 300-y/2)
		text.Draw(screen, str, assets.ScoreFace, op)
	} else if g.paused {
		// draw paused
		str := "PAUSED"
		op := &text.DrawOptions{}
		x, y := text.Measure(str, assets.ScoreFace, op.LineSpacing)
		op.GeoM.Translate(400-x/2, 300-y/2)
		text.Draw(screen, str, assets.ScoreFace, op)
	}
}

func (g *GameInfo) Layout(width, height int) (int, int) {
	return width, height
}

func (g *GameInfo) AddScore(score int) {
	g.score += score
	if g.score > g.highScore {
		g.highScore = g.score
	}
}

func (g *GameInfo) GetScore() int {
	return g.score
}

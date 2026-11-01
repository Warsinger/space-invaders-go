package game

import (
	"fmt"
	"image"
	"image/color"

	"space-invaders/assets"
	comp "space-invaders/components"
	"space-invaders/platform"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Keep the title inside the game, like Tower Defense, with a shared mouse/touch target.
func startButtonBounds(width, height int) image.Rectangle {
	return image.Rect(width/2-160, height-150, width/2+160, height-62)
}

func startButtonPressed(bounds image.Rectangle) bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if image.Pt(x, y).In(bounds) {
			return true
		}
	}
	for _, id := range inpututil.AppendJustPressedTouchIDs(nil) {
		x, y := ebiten.TouchPosition(id)
		if image.Pt(x, y).In(bounds) {
			return true
		}
	}
	return false
}

func centeredText(screen *ebiten.Image, message string, face *text.GoTextFace, width, y float64) {
	options := &text.DrawOptions{}
	options.LineSpacing = face.Size * 1.5
	measured, _ := text.Measure(message, face, options.LineSpacing)
	options.GeoM.Translate((width-measured)/2, y)
	text.Draw(screen, message, face, options)
}

func (g *GameInfo) drawTitle(screen *ebiten.Image) {
	board := comp.Board.Get(comp.Board.MustFirst(g.world))
	width, height := board.Width, board.Height
	screen.Clear()
	screen.DrawImage(assets.GetImage("backgroundV"), &ebiten.DrawImageOptions{})
	titleFace := *assets.ScoreFace
	titleFace.Size = 42
	infoFace := *assets.ScoreFace
	infoFace.Size = 20
	scoreText := fmt.Sprintf("HIGH %05d", g.highScore)
	scoreWidth, _ := text.Measure(scoreText, assets.ScoreFace, 0)
	scoreOptions := &text.DrawOptions{}
	scoreOptions.GeoM.Translate(float64(width)-scoreWidth-20, 20)
	text.Draw(screen, scoreText, assets.ScoreFace, scoreOptions)
	centeredText(screen, "SPACE INVADERS", &titleFace, float64(width), 105)
	centeredText(screen, fmt.Sprintf("Starting level %02d", g.startingLevel), &infoFace, float64(width), 170)
	ship := assets.GetImage("ship")
	shipOptions := &ebiten.DrawImageOptions{}
	shipOptions.GeoM.Scale(2, 2)
	shipOptions.GeoM.Translate(float64(width)/2-float64(ship.Bounds().Dx()), 230)
	screen.DrawImage(ship, shipOptions)
	centeredText(screen, "Defend Earth, one wave at a time.", &infoFace, float64(width), 350)
	for i, name := range []string{"alien1", "alien2", "alien1", "alien2"} {
		img := assets.GetImage(name)
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(float64(width)/2-120+float64(i*64), 400)
		screen.DrawImage(img, options)
	}
	centeredText(screen, "Left / Right  move    Space  fire\nP  pause    R / Esc  title", &infoFace, float64(width), 495)
	if !platform.Browser {
		centeredText(screen, "Q  save and quit", &infoFace, float64(width), 565)
	}
	bounds := startButtonBounds(width, height)
	buttonColor := color.RGBA{5, 50, 200, 255} // Tower Defense's Start Game button palette.
	x, y := ebiten.CursorPosition()
	if image.Pt(x, y).In(bounds) {
		buttonColor = color.RGBA{0, 0, 150, 255}
	}
	vector.DrawFilledRect(screen, float32(bounds.Min.X), float32(bounds.Min.Y), float32(bounds.Dx()), float32(bounds.Dy()), buttonColor, true)
	vector.StrokeRect(screen, float32(bounds.Min.X), float32(bounds.Min.Y), float32(bounds.Dx()), float32(bounds.Dy()), 2, color.RGBA{150, 180, 255, 255}, true)
	centeredText(screen, "Start Game", assets.ScoreFace, float64(width), float64(bounds.Min.Y)+25)
	centeredText(screen, "Click / tap Start Game, or press Space / Enter", &infoFace, float64(width), float64(height-200))
}

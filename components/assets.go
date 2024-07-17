package components

import (
	"bytes"
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
)

var images map[string]*ebiten.Image = make(map[string]*ebiten.Image)

type SpriteData struct {
	image *ebiten.Image
}

var ScoreFace *text.GoTextFace

func LoadAssets() error {
	err := loadImageAsset("alien")
	if err != nil {
		return err
	}
	err = loadImageAsset("ship")
	if err != nil {
		return err
	}

	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.ArcadeN_ttf))
	if err != nil {
		log.Fatal(err)
	}

	ScoreFace = &text.GoTextFace{
		Source: s,
		Size:   24,
	}

	return nil
}

func loadImageAsset(name string) error {
	filepath := fmt.Sprintf("assets/%s.png", name)
	img, _, err := ebitenutil.NewImageFromFile(filepath)
	if err != nil {
		log.Fatalf("failed to load image %v: %v", name, err)
		return err
	}

	images[name] = img
	return nil
}

func GetImage(name string) *ebiten.Image {
	return images[name]
}

func (s *SpriteData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	pos := Position.Get(entry)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(pos.x), float64(pos.y))
	// opts.GeoM.Scale(1, 1)
	screen.DrawImage(s.image, opts)
}

func (s *SpriteData) GetRect(entry *donburi.Entry) image.Rectangle {
	pos := Position.Get(entry)
	rect := s.image.Bounds()
	return rect.Add(image.Pt(pos.x, pos.y))
}

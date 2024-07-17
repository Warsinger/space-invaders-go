package components

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yohamta/donburi"
)

var images map[string]*ebiten.Image = make(map[string]*ebiten.Image)

type SpriteData struct {
	image *ebiten.Image
}

func LoadAssets() error {
	err := loadImageAsset("alien")
	if err != nil {
		return err
	}
	err = loadImageAsset("ship")
	if err != nil {
		return err
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
	opts.GeoM.Translate(float64(pos.X), float64(pos.Y))
	// opts.GeoM.Scale(1, 1)
	screen.DrawImage(s.image, opts)
}

func (s *SpriteData) GetRect(entry *donburi.Entry) image.Rectangle {
	pos := Position.Get(entry)
	rect := s.image.Bounds()
	return rect.Add(image.Pt(pos.X, pos.Y))
}

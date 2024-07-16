package components

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

func (s *SpriteData) Draw(screen *ebiten.Image, pos *PositionData) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(pos.X, pos.Y)
	// opts.GeoM.Scale(1, 1)
	screen.DrawImage(s.image, opts)
}

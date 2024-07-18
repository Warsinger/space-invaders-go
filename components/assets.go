package components

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/donburi"
)

var images map[string]*ebiten.Image = make(map[string]*ebiten.Image)
var sounds map[string]*audio.Player = make(map[string]*audio.Player)

type SpriteData struct {
	image *ebiten.Image
}

type AudioData struct {
	// sound * audio
}

var ScoreFace *text.GoTextFace
var audioContext *audio.Context

func LoadAssets() error {
	err := loadImages()
	if err != nil {
		return err
	}

	err = loadFonts()
	if err != nil {
		return err
	}

	err = loadAudio()
	if err != nil {
		return err
	}

	return nil
}
func loadFonts() error {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.ArcadeN_ttf))
	if err != nil {
		return err
	}

	ScoreFace = &text.GoTextFace{
		Source: s,
		Size:   24,
	}
	return nil
}

func loadImages() error {
	err := loadImageAsset("alien1")
	if err != nil {
		return err
	}
	err = loadImageAsset("alien2")
	if err != nil {
		return err
	}
	err = loadImageAsset("ship")
	if err != nil {
		return err
	}
	err = loadImageAsset("background")
	if err != nil {
		return err
	}
	return nil
}

func loadAudio() error {
	audioContext = audio.NewContext(44100)
	err := loadAudioAsset("explosion")
	if err != nil {
		return err
	}
	err = loadAudioAsset("killed")
	if err != nil {
		return err
	}
	err = loadAudioAsset("shoot")
	if err != nil {
		return err
	}

	return nil
}

func loadAudioAsset(name string) error {
	filepath := fmt.Sprintf("assets/sounds/%s.wav", name)
	audioFile, err := os.Open(filepath)
	if err != nil {
		return err
	}

	d, err := wav.DecodeWithoutResampling(audioFile)
	if err != nil {
		return err
	}
	player, err := audioContext.NewPlayer(d)
	if err != nil {
		return err
	}
	sounds[name] = player

	return nil
}

func loadImageAsset(name string) error {
	filepath := fmt.Sprintf("assets/images/%s.png", name)
	img, _, err := ebitenutil.NewImageFromFile(filepath)
	if err != nil {
		log.Fatalf("failed to load image %v: %v", name, err)
		return err
	}

	images[name] = img
	return nil
}

func PlaySound(name string) {
	sound := sounds[name]
	sound.Rewind()
	sound.Play()
}

func GetImage(name string) *ebiten.Image {
	return images[name]
}

func (s *SpriteData) Draw(screen *ebiten.Image, entry *donburi.Entry) {
	pos := Position.Get(entry)
	v := Velocity.Get(entry)
	opts := &ebiten.DrawImageOptions{}

	if v.x < 0 {
		// this flips the image when it is going to the left
		opts.GeoM.Translate(-float64(pos.x+s.image.Bounds().Dx()), float64(pos.y))
		opts.GeoM.Scale(-1, 1)
	} else {
		opts.GeoM.Translate(float64(pos.x), float64(pos.y))
	}
	screen.DrawImage(s.image, opts)
}

func (s *SpriteData) GetRect(entry *donburi.Entry) image.Rectangle {
	pos := Position.Get(entry)
	rect := s.image.Bounds()
	return rect.Add(image.Pt(pos.x, pos.y))
}

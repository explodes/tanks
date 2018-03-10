package ebu

import "github.com/hajimehoshi/ebiten"

type Scene interface {
	Update(dt float64, image *ebiten.Image) error
}

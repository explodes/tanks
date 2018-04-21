package core

import "github.com/hajimehoshi/ebiten"

type Scene interface {
	Update(dt float64) error
	Draw(image *ebiten.Image)
}

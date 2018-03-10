package eb1

import (
	"fmt"

	"github.com/explodes/scratch/ebgames/ebu"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type titleScene struct {
	g *Game
}

func NewTitleScene(game *Game) ebu.Scene {
	return &titleScene{
		g: game,
	}
}

func (s *titleScene) Update(dt float64, image *ebiten.Image) error {
	ebitenutil.DebugPrint(image, fmt.Sprintf("%.05f", dt))
	return nil
}

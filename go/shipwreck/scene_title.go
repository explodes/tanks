package tanks

import (
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/colornames"

	_ "image/png"

	"github.com/explodes/tanks/go/core"
	"github.com/explodes/tempura"
)

var _ core.Scene = (*titleScene)(nil)

type titleScene struct {
	g    *Game
	time float64

	title tempura.Text
}

func NewTitleScene(game *Game) (core.Scene, error) {
	loader := game.context.Loader()

	titleFace, err := loader.Face("fonts/DampfPlatz.ttf", 120)
	if err != nil {
		return nil, err
	}

	s := &titleScene{
		g:     game,
		title: tempura.NewText(titleFace, colornames.White, Title),
	}
	return s, nil
}

func (s *titleScene) Update(dt float64) error {
	s.time += dt

	if Begin() {
		return s.g.SetNewScene(NewGameScene)
	}

	return nil
}

func (s *titleScene) Draw(image *ebiten.Image) {
	s.title.Draw(image, core.ScreenWidth/2-s.title.W/2, core.ScreenHeight/2+s.title.H/2, tempura.AlignLeft)
}

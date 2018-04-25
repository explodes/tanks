package tanks

import (
	"github.com/explodes/tanks/go/core"
	"github.com/explodes/tempura"
)

const (
	Title = "Shipwreck"
)

var _ core.Game = (*Game)(nil)

type Game struct {
	core.GameSceneLoop
	context core.Context

	redScore  int
	blueScore int
}

func NewGame(context core.Context) (core.Game, error) {
	if core.Debug {
		defer tempura.LogStart("%s init", Title).End()
	}

	game := &Game{
		context: context,
	}

	if err := game.SetNewScene(NewTitleScene); err != nil {
		return nil, err
	}

	return game, nil
}

func (g *Game) SetNewScene(factory func(*Game) (scene core.Scene, err error)) error {
	if core.Debug {
		defer tempura.LogStart("New %s scene", Title).End()
	}
	scene, err := factory(g)
	if err != nil {
		return err
	}
	return g.SetScene(scene)
}

func (g *Game) OnMuted(muted bool) {
}

func (g *Game) Close() error {
	return nil
}

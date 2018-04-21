package tanks

import (
	"github.com/explodes/tanks/go/core"
	"github.com/explodes/tempura"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
)

const (
	Title = "Tanks"

	bgmVolume = 0.5
)

var _ core.Game = (*Game)(nil)

type Game struct {
	context core.Context
	scene   core.Scene

	redScore  int
	blueScore int

	bgm *audio.Player
}

func NewGame(context core.Context) (core.Game, error) {
	if core.Debug {
		defer tempura.LogStart("Tanks game init").End()
	}
	bgm, err := context.Loader().AudioLoop(context.AudioContext(), "mp3", "music/octane.mp3")
	if err != nil {
		return nil, err
	}

	game := &Game{
		context: context,
		bgm:     bgm,
	}

	if err := game.SetNewScene(NewTitleScene); err != nil {
		return nil, err
	}

	bgm.SetVolume(bgmVolume)
	bgm.Play()

	return game, nil
}

func (g *Game) SetNewScene(factory func(*Game) (scene core.Scene, err error)) error {
	if core.Debug {
		defer tempura.LogStart("Set new tank scene").End()
	}
	scene, err := factory(g)
	if err != nil {
		return err
	}
	return g.SetScene(scene)
}

func (g *Game) SetScene(scene core.Scene) error {
	core.DebugLog("new tank scene: %T", scene)
	g.scene = scene
	return nil
}

func (g *Game) Update(dt float64) error {
	if g.scene != nil {
		if err := g.scene.Update(dt); err != nil {
			return err
		}
	}
	return nil
}

func (g *Game) Draw(image *ebiten.Image) {
	if g.scene != nil {
		g.scene.Draw(image)
	}
}

func (g *Game) OnMuted(muted bool) {
	if muted {
		g.bgm.SetVolume(0)
	} else {
		g.bgm.SetVolume(bgmVolume)
	}
}

func (g *Game) Close() error {
	return g.bgm.Close()
}

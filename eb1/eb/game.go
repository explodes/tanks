package eb1

import (
	"github.com/explodes/scratch/ebgames/ebu"
	"github.com/hajimehoshi/ebiten"
)

const (
	ScreenWidth  = 320
	ScreenHeight = 240
)

type Game struct {
	time      float64
	stopwatch ebu.Stopwatch
	scene     ebu.Scene
}

func NewGame() (*Game, error) {
	game := &Game{
		stopwatch: ebu.NewStopwatch(),
	}

	game.SetScene(NewTitleScene(game))

	return game, nil
}

func (g *Game) SetScene(scene ebu.Scene) {
	g.scene = scene
}

func (g *Game) Update(image *ebiten.Image) error {
	if ebiten.IsRunningSlowly() {
		return nil
	}
	if g.scene != nil {
		return g.scene.Update(g.stopwatch.TimeDelta(), image)
	}
	return nil
}

func (g *Game) Pause() {
	g.stopwatch.Pause()
}

func (g *Game) Resume() {
	g.stopwatch.Resume()
}

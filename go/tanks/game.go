package tanks

import (
	"math/rand"
	"time"

	"github.com/explodes/tanks/go/tanks/res"
	"github.com/explodes/tanks/go/tanksutil"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
)

const (
	Title        = "Tanks"
	ScreenWidth  = 768
	ScreenHeight = 432

	audioSampleRate = 44100
)

var (
	ScreenBounds = tanksutil.R(0, 0, ScreenWidth, ScreenHeight)
)

type Game struct {
	time         float64
	loader       *tanksutil.Loader
	stopwatch    tanksutil.Stopwatch
	scene        Scene
	input        Input
	audioContext *audio.Context

	redScore  int
	blueScore int
}

type Scene interface {
	Update(dt float64) error
	Draw(image *ebiten.Image)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewGame() (*Game, error) {
	loader := tanksutil.NewLoader(res.Asset)
	audioContext, err := audio.NewContext(audioSampleRate)
	if err != nil {
		return nil, err
	}
	bgm, err := loader.AudioLoop(audioContext, "mp3", "music/octane.mp3")
	if err != nil {
		return nil, err
	}

	game := &Game{
		loader:       loader,
		stopwatch:    tanksutil.NewStopwatch(),
		input:        NewInput(),
		audioContext: audioContext,
	}

	scene, err := NewTitleScene(game)
	if err != nil {
		return nil, err
	}
	game.SetScene(scene)

	bgm.SetVolume(0.5)
	bgm.Play()

	return game, nil
}

func (g *Game) SetNewScene(scene Scene, err error) error {
	if err != nil {
		return err
	}
	return g.SetScene(scene)
}

func (g *Game) SetScene(scene Scene) error {
	DebugLog("new scene: %T", scene)
	g.scene = scene
	return nil
}

func (g *Game) Update(image *ebiten.Image) error {
	dt := g.stopwatch.TimeDelta()
	g.time += dt

	if g.scene != nil {
		if err := g.scene.Update(dt); err != nil {
			return err
		}
		if !ebiten.IsRunningSlowly() {
			g.scene.Draw(image)
		}
	}

	return nil
}

func (g *Game) Pause() {
	g.stopwatch.Pause()
}

func (g *Game) Resume() {
	g.stopwatch.Resume()
}

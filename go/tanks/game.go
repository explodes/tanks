package tanks

import (
	"math/rand"
	"time"

	"errors"

	"github.com/explodes/tanks/go/tanks/res"
	"github.com/explodes/tempura"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
)

const (
	Title        = "Tanks"
	ScreenWidth  = 768
	ScreenHeight = 432

	audioSampleRate = 44100

	bgmVolume = 0.5
)

var (
	ScreenBounds = tempura.R(0, 0, ScreenWidth, ScreenHeight)

	regularTermination = errors.New("goodbye!")
)

type Game struct {
	time         float64
	loader       tempura.Loader
	stopwatch    tempura.Stopwatch
	scene        Scene
	input        Input
	audioContext *audio.Context

	fullscreen bool
	muted      bool

	redScore  int
	blueScore int

	bgm *audio.Player
}

type Scene interface {
	Update(dt float64) error
	Draw(image *ebiten.Image)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewGame() (*Game, error) {
	if debug {
		defer tempura.LogStart("Game init").End()
	}
	loader := tempura.NewCachedLoader(tempura.NewLoaderDebug(res.Asset, debug))
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
		stopwatch:    tempura.NewStopwatch(),
		input:        NewInput(),
		audioContext: audioContext,
		bgm:          bgm,
	}

	if err := game.SetNewScene(NewTitleScene); err != nil {
		return nil, err
	}

	bgm.SetVolume(bgmVolume)
	bgm.Play()

	return game, nil
}

func (g *Game) SetNewScene(factory func(*Game) (scene Scene, err error)) error {
	if debug {
		defer tempura.LogStart("Set New Scene").End()
	}
	scene, err := factory(g)
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

	if g.input.Exit() {
		return regularTermination
	}

	if g.input.ToggleFullscreen() {
		g.fullscreen = !g.fullscreen
		ebiten.SetFullscreen(g.fullscreen)
	}

	if g.input.ToggleMute() {
		g.muted = !g.muted
		if g.muted {
			g.bgm.SetVolume(0)
		} else {
			g.bgm.SetVolume(bgmVolume)
		}
	}

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

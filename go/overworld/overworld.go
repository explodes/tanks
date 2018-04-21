package overworld

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/explodes/tanks/go/core"
	"github.com/explodes/tanks/go/games"
	"github.com/explodes/tanks/go/res"
	"github.com/explodes/tempura"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/audio"
	"github.com/pkg/errors"
)

var _ core.Context = (*Overworld)(nil)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Overworld struct {
	loader       tempura.Loader
	audioContext *audio.Context
	muted        bool
	stopwatch    tempura.Stopwatch
	fullscreen   bool

	game core.Game
}

func NewOverworld() (*Overworld, error) {
	loader := tempura.NewCachedLoader(tempura.NewLoaderDebug(res.Asset, core.Debug))
	audioContext, err := audio.NewContext(core.AudioSampleRate)
	if err != nil {
		return nil, err
	}
	overworld := &Overworld{
		loader:       loader,
		audioContext: audioContext,
		muted:        false,
		stopwatch:    tempura.NewStopwatch(),
		fullscreen:   false,
	}
	if err := overworld.startOverworldGame(); err != nil {
		return nil, err
	}
	return overworld, nil
}

func (o *Overworld) Loader() tempura.Loader {
	return o.loader
}

func (o *Overworld) AudioContext() *audio.Context {
	return o.audioContext
}

func (o *Overworld) Muted() bool {
	return o.muted
}

func (o *Overworld) startOverworldGame() error {
	return o.LoadGame("tanks")
}

func (o *Overworld) LoadGame(name string) error {
	if core.Debug {
		defer tempura.LogStart(fmt.Sprintf("load game %s", name)).End()
	}
	factory := games.GetGameFactory(name)
	if factory == nil {
		return errors.Errorf("%s game does not exist", name)
	}
	game, err := factory(o)
	if err != nil {
		return err
	}
	return o.setGame(game)
}
func (o *Overworld) setGame(game core.Game) error {
	core.DebugLog("new game: %T", game)
	o.game = game
	return nil
}

func (o *Overworld) Update(image *ebiten.Image) error {
	dt := o.stopwatch.TimeDelta()
	if Exit() {
		core.DebugLog("exit")
		return core.RegularTermination
	}
	if ToggleFullscreen() {
		o.fullscreen = !o.fullscreen
		core.DebugLog("fullscreen: %v", o.fullscreen)
		ebiten.SetFullscreen(o.fullscreen)
	}
	if ToggleMute() {
		o.muted = !o.muted
		core.DebugLog("mute: %v", o.muted)
		o.game.OnMuted(o.muted)
	}
	if o.game != nil {
		if err := o.game.Update(dt); err == core.GameTermination {
			o.closeGame()
			return o.startOverworldGame()
		} else if err != nil {
			o.closeGame()
			return err
		}
		if !ebiten.IsRunningSlowly() {
			o.game.Draw(image)
		}
	}
	return nil
}

func (o *Overworld) closeGame() {
	if o.game == nil {
		return
	}
	if err := o.game.Close(); err != nil {
		core.Log("error closing game: %v", err)
	}
	o.game = nil
}

func (o *Overworld) Pause() {
	core.DebugLog("paused")
	o.stopwatch.Pause()
}

func (o *Overworld) Resume() {
	core.DebugLog("resumed")
	o.stopwatch.Resume()
}

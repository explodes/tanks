package mobile

import (
	"github.com/explodes/tanks/go/tanks"
	"github.com/hajimehoshi/ebiten/mobile"
)

var (
	running bool
	game    *tanks.Game
)

const (
	ScreenWidth  = tanks.ScreenWidth
	ScreenHeight = tanks.ScreenHeight
)

// IsRunning returns a boolean value indicating whether the game is running.
func IsRunning() bool {
	return running
}

// Start starts the game.
func Start(scale float64) error {
	running = true
	var err error
	game, err = tanks.NewGame()
	if err != nil {
		return err
	}
	if err := mobile.Start(game.Update, ScreenWidth, ScreenHeight, scale, tanks.Title); err != nil {
		return err
	}
	return nil
}

// Update proceeds the game.
func Update() error {
	return mobile.Update()
}

func Pause() {
	if game != nil {
		game.Pause()
	}
}

func Resume() {
	if game != nil {
		game.Resume()
	}
}

// UpdateTouchesOnAndroid dispatches touch events on Android.
func UpdateTouchesOnAndroid(action int, id int, x, y int) {
	mobile.UpdateTouchesOnAndroid(action, id, x, y)
}

// UpdateTouchesOnIOS dispatches touch events on iOS.
func UpdateTouchesOnIOS(phase int, ptr int64, x, y int) {
	// Prepare this function if you also want to make your game run on iOS.
	mobile.UpdateTouchesOnIOS(phase, ptr, x, y)
}

// +build darwin freebsd linux windows js
// +build !android
// +build !ios

package tanks

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

func (i *inputImpl) ToggleFullscreen() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyF)
}

func (i *inputImpl) Begin() bool {
	return ebiten.IsKeyPressed(ebiten.KeySpace)
}

func (i *inputImpl) BlueRotate() bool {
	return ebiten.IsKeyPressed(ebiten.KeyA)
}

func (i *inputImpl) RedRotate() bool {
	return ebiten.IsKeyPressed(ebiten.KeyL)
}

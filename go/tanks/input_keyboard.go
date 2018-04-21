// +build darwin freebsd linux windows js
// +build !android
// +build !ios

package tanks

import (
	"github.com/hajimehoshi/ebiten"
)

func Begin() bool {
	return ebiten.IsKeyPressed(ebiten.KeySpace)
}

func BlueRotate() bool {
	return ebiten.IsKeyPressed(ebiten.KeyA)
}

func RedRotate() bool {
	return ebiten.IsKeyPressed(ebiten.KeyL)
}

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

func BlueSwitchDirection() bool {
	return ebiten.IsKeyPressed(ebiten.KeyA)
}

func RedSwitchDirection() bool {
	return ebiten.IsKeyPressed(ebiten.KeyL)
}
